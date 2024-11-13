package prepare_data

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func intToByteSlice(n int) []byte {
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, uint16(n))
	return result
}

func writeTokensToFile(tokens [][]int, filename string) error {
	if filepath.Ext(filename) == "" {
		filename += ".bin"
	}
	var tokensBytes bytes.Buffer
	for _, t := range tokens {
		for _, token := range t {
			tokensBytes.Write(intToByteSlice(token))
		}
	}
	err := os.WriteFile(filename, tokensBytes.Bytes(), 0644)
	return err
}

func writeJsonToFile(t interface{}, filename string) error {
	if filepath.Ext(filename) == "" {
		filename += ".json"
	}
	keyFile, err := calcSha(filename)
	if err != nil {
		return err
	}
	json, err := JSONMarshal(t)
	if err != nil {
		return err
	}
	keyBytes := hashBytes(json)
	if keyFile == keyBytes {
		return nil
	}
	return os.WriteFile(filename, json, 0644)
}

func FindFiles(root, ext string) (files []string, err error) {
	err = filepath.WalkDir(root, func(fileName string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			files = append(files, fileName)
		}
		return nil
	})
	return
}
