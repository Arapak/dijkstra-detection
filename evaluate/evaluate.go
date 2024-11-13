package evaluate

import (
	"ai-dijkstra/config"
	"ai-dijkstra/prepare_data"
	"ai-dijkstra/tokenize_code"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Evaluate(tokensNumeric [][]int, modelFile string, cpuOnly bool) (predictions []float64, err error) {
	device := "GPU"
	if cpuOnly {
		device = "CPU"
	}
	command := fmt.Sprintf("python3 evaluate/evaluate.py %v %v", modelFile, device)
	cmds := tokenize_code.SplitCmd(command)
	cmd := exec.Command(cmds[0], cmds[1:]...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	var input bytes.Buffer
	fmt.Fprintf(&input, strconv.Itoa(len(tokensNumeric))+"\n")
	for _, tokens := range tokensNumeric {
		for _, token := range tokens {
			fmt.Fprintf(&input, strconv.Itoa(token)+"\n")
		}
	}
	cmd.Stdin = bytes.NewReader(input.Bytes())
	err = cmd.Run()
	if err != nil {
		return
	}

	output := strings.TrimSpace(stdout.String())
	predictionsString := strings.Split(output, "\n")
	predictions = make([]float64, len(tokensNumeric))
	for i := range predictions {
		predictions[i], err = strconv.ParseFloat(predictionsString[i+len(predictionsString)-len(predictions)], 32)
		if err != nil {
			return
		}
	}
	return
}

func LoadDicitonary(dictionaryFile string) (dict map[string]int, err error) {
	dictionaryBytes, err := os.ReadFile(dictionaryFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(dictionaryBytes, &dict)
	return
}

func EvaluateFile(file, modelFile, dictionaryFile string, cpuOnly bool) (prediction float64, err error) {
	tokens, err := tokenize_code.TokenizeCode(file)
	if err != nil {
		return
	}
	dict, err := LoadDicitonary(dictionaryFile)
	if err != nil {
		return
	}
	tokensNumeric := prepare_data.ConvertToNumeric(dict, tokens)
	tokensNumericArray := [][]int{tokensNumeric}
	predictions, err := Evaluate(tokensNumericArray, modelFile, cpuOnly)
	if err != nil {
		return
	}
	return predictions[0], nil
}

func EvaluateDir(path, modelFile, dictionaryFile string, useCache bool, cpuOnly bool) (files []string, predictions []float64, err error) {
	files, tokens, err := prepare_data.TokenizeFiles(path, config.NumberOfWorkers, useCache)
	if err != nil {
		return
	}
	dict, err := LoadDicitonary(dictionaryFile)
	if err != nil {
		return
	}
	tokensNumeric := make([][]int, len(tokens))
	for i := range tokens {
		tokensNumeric[i] = prepare_data.ConvertToNumeric(dict, tokens[i])
	}
	predictions, err = Evaluate(tokensNumeric, modelFile, cpuOnly)
	return
}
