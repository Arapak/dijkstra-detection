package cmd

import (
	"ai-dijkstra/config"
	"ai-dijkstra/evaluate"
	"ai-dijkstra/prepare_data"
	"ai-dijkstra/pull_codes"
	"ai-dijkstra/train"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/docopt/docopt-go"

	"github.com/fatih/color"
)

type ParsedArgs struct {
	PrepareData     bool `docopt:"prepare_data"`
	TestPrepareData bool `docopt:"test_prepare_data"`
	Train           bool `docopt:"train"`
	Run             bool `docopt:"run"`
	EvaluateFile    bool `docopt:"evaluate_file"`
	EvaluateDir     bool `docopt:"evaluate_dir"`
	PullCodes       bool `docopt:"pull_codes"`
	Cache           bool `docopt:"--cache"`
	Cpu             bool `docopt:"--cpu"`
	Epochs          string
	Path            string
}

var Args *ParsedArgs

func Eval(opts docopt.Opts) error {
	Args = &ParsedArgs{}
	err := opts.Bind(Args)
	if err != nil {
		return err
	}
	if Args.PrepareData {
		return PrepareData()
	} else if Args.TestPrepareData {
		mockPaths := []prepare_data.DataPaths{{Data: config.MockDataPath + "1", Result: config.ResultsMockDataPath + "data1"}, {Data: config.MockDataPath + "0", Result: config.ResultsMockDataPath + "data0"}}
		return prepare_data.PrepareData(mockPaths, config.MockDicionaryPath, config.NumberOfWorkers)
	} else if Args.Train {
		return Train()
	} else if Args.Run {
		err := PrepareData()
		if err != nil {
			return err
		}
		return Train()
	} else if Args.EvaluateFile {
		prediction, err := evaluate.EvaluateFile(Args.Path, config.ModelPath, config.DicionaryPath, Args.Cpu)
		fmt.Printf("%0.3f\n", prediction)
		return err
	} else if Args.EvaluateDir {
		files, predictions, err := evaluate.EvaluateDir(Args.Path, config.ModelPath, config.DicionaryPath, Args.Cache, Args.Cpu)
		array := make([]FilePredictions, len(files))
		for i := range files {
			array[i] = FilePredictions{files[i], predictions[i]}
		}
		sort.Sort(FilePredictionsArray(array))
		for i := range array {
			fmt.Printf("%v %0.3f\n", array[i].Filename, array[i].Prediction)
		}
		return err
	} else if Args.PullCodes {
		dijkstraProblems := []pull_codes.Problem{{ContestID: "786", ProblemID: "B"}}
		// dijkstraProblems := []pull_codes.Problem{{ContestID: "229", ProblemID: "B"}, {ContestID: "545", ProblemID: "E"}}
		// dijkstraProblems := []pull_codes.Problem{{ContestID: "757", ProblemID: "F"},{ContestID: "716", ProblemID: "D"}}
		// noDijkstraProblems := []pull_codes.Problem{{ContestID: "52", ProblemID: "C"}, {ContestID: "617", ProblemID: "E"}, {ContestID: "1657", ProblemID: "D"}, {ContestID: "1343", ProblemID: "E"}, {ContestID: "746", ProblemID: "D"}, {ContestID: "1732", ProblemID: "C1"}, {ContestID: "1228", ProblemID: "D"}, {ContestID: "371", ProblemID: "D"}}
		noDijkstraProblems := []pull_codes.Problem{{ContestID: "1493", ProblemID: "D"}, {ContestID: "1183", ProblemID: "E"}, {ContestID: "1712", ProblemID: "D"}, {ContestID: "1366", ProblemID: "D"}, {ContestID: "371", ProblemID: "D"}}
		// noDijkstraProblems := []pull_codes.Problem{{ContestID: "715", ProblemID: "C"}, {ContestID: "741", ProblemID: "D"}, {ContestID: "723", ProblemID: "F"}, {ContestID: "687", ProblemID: "D"}, {ContestID: "746", ProblemID: "G"}, {ContestID: "750", ProblemID: "F"}, {ContestID: "828", ProblemID: "F"}, {ContestID: "832", ProblemID: "D"}}
		// noDijkstraProblems := []pull_codes.Problem{{ContestID: "769", ProblemID: "C"}, {ContestID: "796", ProblemID: "D"}, {ContestID: "821", ProblemID: "D"}, {ContestID: "734", ProblemID: "E"}, {ContestID: "802", ProblemID: "K"}, {ContestID: "846", ProblemID: "E"}, {ContestID: "718", ProblemID: "D"}}
		// problems := []Problem{{"1763", "B"}, {"715", "B"}, {"843", "D"}, {"757", "F"}}
		// 1763B - no dijsktra, else dijkstra
		// problems := [2]Problem{{"1779", "C"}, {"1779", "D"}}
		// problems := [5]Problem{{"59", "E"}, {"449", "B"}, {"464", "E"}, {"567", "E"}, {"715", "B"}}
		dijkstraProblems = append(dijkstraProblems, noDijkstraProblems...)
		return pull_codes.PullCodes(dijkstraProblems)
	}
	color.Red("this function is not available")
	return nil
}

func PrepareData() error {
	testPaths := []prepare_data.DataPaths{{Data: config.TestDataPath + "1", Result: config.ResultsTestDataPath + "data1"}, {Data: config.TestDataPath + "0", Result: config.ResultsTestDataPath + "data0"}}
	validationPaths := []prepare_data.DataPaths{{Data: config.ValidationDataPath + "1", Result: config.ResultsValidationDataPath + "data1"}, {Data: config.ValidationDataPath + "0", Result: config.ResultsValidationDataPath + "data0"}}
	trainPaths := []prepare_data.DataPaths{{Data: config.TrainDataPath + "1", Result: config.ResultsTrainDataPath + "data1"}, {Data: config.TrainDataPath + "0", Result: config.ResultsTrainDataPath + "data0"}}
	return prepare_data.PrepareData(append(append(testPaths, validationPaths...), trainPaths...), config.DicionaryPath, config.NumberOfWorkers)
}

func Train() (err error) {
	epochs := config.DefaultNumOfEpochs
	if Args.Epochs != "" {
		epochs, err = strconv.Atoi(Args.Epochs)
		if err != nil || epochs <= 0 {
			return errors.New("arg epochs must be a positive integer")
		}
	}
	return train.Train(config.ResultsTrainDataPath, config.ResultsTestDataPath, config.ResultsValidationDataPath, config.DicionaryPath, epochs, config.ModelPath, Args.Cpu)
}
