package tokenize_code

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode"

	"github.com/go-clang/clang-v11/clang"
)

func SplitCmd(s string) (res []string) {
	// https://github.com/vrischmann/shlex/blob/master/shlex.go
	var buf bytes.Buffer
	insideQuotes := false
	for _, r := range s {
		switch {
		case unicode.IsSpace(r) && !insideQuotes:
			if buf.Len() > 0 {
				res = append(res, buf.String())
				buf.Reset()
			}
		case r == '"' || r == '\'':
			if insideQuotes {
				res = append(res, buf.String())
				buf.Reset()
				insideQuotes = false
				continue
			}
			insideQuotes = true
		default:
			buf.WriteRune(r)
		}
	}
	if buf.Len() > 0 {
		res = append(res, buf.String())
	}
	return
}

var macroExpansionCommand = `clang-15 -DONLINE_JUDGE -E -nostdinc++ -nobuiltininc -nostdlibinc -P -dI %v`
var createASTCommand = `clang++ -std=c++17 -DONLINE_JUDGE -include-pch /usr/include/x86_64-linux-gnu/c++/10/bits/stdc++.h.pch -emit-ast %v -o %v`

func runCommand(command string, outFile *os.File) error {
	cmds := SplitCmd(command)
	cmd := exec.Command(cmds[0], cmds[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if outFile != nil {
		cmd.Stdout = outFile
	}
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err.Error(), stderr.String())
	}
	return err
}

func createAST(filename string) (string, string, error) {
	snippetCode, err := os.CreateTemp("snippets", "snippetCode*.cpp")
	if err != nil {
		log.Fatal(err)
	}
	runCommand(fmt.Sprintf(macroExpansionCommand, filename), snippetCode)
	snippetAst, err := os.CreateTemp("snippets", "snippetAst*.ast")
	if err != nil {
		log.Fatal(err)
	}
	err = runCommand(fmt.Sprintf(createASTCommand, snippetCode.Name(), snippetAst.Name()), nil)
	return snippetCode.Name(), snippetAst.Name(), err
}

func getNames(tu *clang.TranslationUnit, cursor *clang.Cursor) (names map[string]int) {
	names = make(map[string]int)
	location := tu.Spelling()
	extractNames := func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		file, _, _, _ := cursor.Location().SpellingLocation()
		if file.Name() != location {
			return clang.ChildVisit_Continue
		}

		if cursor.Kind().IsDeclaration() && cursor.Spelling() != "" {
			if _, ok := names[cursor.Spelling()]; !ok {
				names[cursor.Spelling()] = len(names)
			}
		}
		return clang.ChildVisit_Recurse
	}
	cursor.Visit(extractNames)
	return
}

const Identifier = "<identifier>"
const StringLiteral = "<string literal>"
const CharacterLiteral = "<charcter literal>"
const NumericLiteral = "<numeric literal>"
const UnknownToken = "<unknown>"

func contains(m map[string]int, s string) bool {
	_, ok := m[s]
	return ok
}

func TokenizeCode(filename string) (tokens []string, err error) {
	codeTmpFile, astTmpFile, err := createAST(filename)
	defer os.Remove(codeTmpFile)
	defer os.Remove(astTmpFile)
	if err != nil {
		return
	}

	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()
	tu := idx.TranslationUnit(astTmpFile)
	defer tu.Dispose()
	cursor := tu.TranslationUnitCursor()

	names := getNames(&tu, &cursor)

	if contains(names, "spfa") || contains(names, "SPFA") || contains(names, "Spfa") {
		err = errors.New("file contains spfa")
		return
	}

	clangTokens := tu.Tokenize(cursor.Extent())

	for _, token := range clangTokens {
		spelling := tu.TokenSpelling(token)
		if token.Kind() == clang.Token_Comment {
			continue
		} else if _, ok := names[spelling]; ok && token.Kind() == clang.Token_Identifier {
			tokens = append(tokens, Identifier)
		} else if token.Kind() == clang.Token_Literal && !(spelling == "0" || spelling == "1") {
			if strings.HasPrefix(spelling, `"`) && strings.HasSuffix(spelling, `"`) {
				tokens = append(tokens, StringLiteral)
			} else if strings.HasPrefix(spelling, `'`) && strings.HasSuffix(spelling, `'`) {
				tokens = append(tokens, CharacterLiteral)
			} else {
				tokens = append(tokens, NumericLiteral)
			}
		} else {
			tokens = append(tokens, spelling)
		}
	}
	return
}
