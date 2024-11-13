package train

import (
	"ai-dijkstra/tokenize_code"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

const trainCommand = "python3 train/train.py %v %v %v %v %v %v %v"

func Train(trainDataPath, testDataPath, validationDataPath, dictionaryPath string, numberOfEpochs int, exportPath string, onlyCpu bool) error {
	device := "GPU"
	if onlyCpu {
		device = "CPU"
	}
	command := fmt.Sprintf(trainCommand, trainDataPath, testDataPath, validationDataPath, dictionaryPath, numberOfEpochs, exportPath, device)
	cmds := tokenize_code.SplitCmd(command)
	cmd := exec.Command(cmds[0], cmds[1:]...)
	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v\n%v", err.Error(), stderr.String())
	}
	return err
}
