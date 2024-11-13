package prepare_data

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"ai-dijkstra/tokenize_code"

	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
)

func createDictionary(tokens [][]string) map[string]int {
	allTokens := make(map[string]bool)
	for _, t := range tokens {
		for _, token := range t {
			allTokens[token] = true
		}
	}
	allTokens[tokenize_code.UnknownToken] = true

	sortedTokens := make([]string, len(allTokens))
	i := 0
	for token := range allTokens {
		sortedTokens[i] = token
		i++
	}
	sort.Strings(sortedTokens)

	dict := make(map[string]int)
	for _, token := range sortedTokens {
		dict[token] = len(dict) + 1
	}
	return dict
}

const maxSeqLength = 1024

func ConvertToNumeric(dict map[string]int, tokens []string) (tokensNumeric []int) {
	tokensNumeric = make([]int, maxSeqLength)
	for i, token := range tokens {
		if i == maxSeqLength {
			break
		}
		if tokenId, ok := dict[token]; ok {
			tokensNumeric[i] = tokenId
		} else {
			tokensNumeric[i] = dict["unknown"]
		}
	}
	return
}

func printStatus(finished, errored, all int, startTime time.Time, last bool) {
	ansi.EraseInLine(2)
	message := fmt.Sprintf("finished: %v/%v errored: %v time: %v", finished, all, errored, time.Since(startTime).Round(time.Millisecond).String())
	if !last {
		message += fmt.Sprintf(" estimated_left: %v", time.Duration(int64(time.Since(startTime))*int64(all-finished)/int64(finished)).Round(time.Millisecond).String())
	}
	ansi.Printf("%v\n", message)
	if !last {
		ansi.CursorUp(1)
	}
}

func TokenizeFiles(dataPath string, numberOfWorkers int, useCache bool) (filenames []string, tokens [][]string, err error) {
	if useCache && cache == nil {
		err = readCache()
		if err != nil {
			return
		}
	}

	files, err := FindFiles(dataPath, ".cpp")
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(numberOfWorkers)
	mu := sync.Mutex{}
	index := 0
	finished := 0
	errored := 0
	start := time.Now()
	for i := 1; i <= numberOfWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for {
				mu.Lock()
				if index >= len(files) {
					mu.Unlock()
					break
				}
				code := files[index]
				index++

				mu.Unlock()

				var currentTokens []string
				var err error
				if useCache {
					currentTokens, err = tokenizeCodeCached(code)
				} else {
					currentTokens, err = tokenize_code.TokenizeCode(code)
				}

				mu.Lock()
				finished++
				if err != nil {
					errored++
				} else {
					filenames = append(filenames, code)
					tokens = append(tokens, currentTokens)
				}
				if finished%5 == 0 {
					printStatus(finished, errored, len(files), start, false)
				}
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	printStatus(finished, errored, len(files), start, true)
	err = writeCache()
	return

}

func ConvertTokensToNumeric(tokens [][]string) (dict map[string]int, tokensNumeric [][]int) {
	dict = createDictionary(tokens)

	for i := range tokens {
		if len(tokens[i]) == 0 {
			continue
		}
		tokensNumeric = append(tokensNumeric, ConvertToNumeric(dict, tokens[i]))
	}
	return
}

type DataPaths struct {
	Data   string
	Result string
}

func PrepareData(paths []DataPaths, dictionaryPath string, numberOfWorkers int) error {
	startTime := time.Now()

	tokens := make([][]string, len(paths))
	lengths := make([]int, len(paths)+1)
	for i := range tokens {
		color.Blue("Tokenize %v", paths[i].Data)
		_, currentTokens, err := TokenizeFiles(paths[i].Data, numberOfWorkers, true)
		if err != nil {
			return err
		}
		lengths[i+1] = lengths[i] + len(currentTokens)
		tokens = append(tokens, currentTokens...)
	}

	dict, tokensNumeric := ConvertTokensToNumeric(tokens)

	for i := range paths {
		err := writeTokensToFile(tokensNumeric[lengths[i]:lengths[i+1]], paths[i].Result)
		// err = writeJsonToFile(tokensNumeric[lengths[i]:lengths[i+1]], paths[i].Result)
		if err != nil {
			return err
		}
	}

	err := writeJsonToFile(dict, dictionaryPath)
	if err != nil {
		return err
	}
	color.Green("TIME: %v\n", time.Since(startTime).Round(time.Millisecond).String())
	return nil
}
