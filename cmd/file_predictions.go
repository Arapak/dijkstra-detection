package cmd

type FilePredictions struct {
	Filename   string
	Prediction float64
}
type FilePredictionsArray []FilePredictions

func (s FilePredictionsArray) Len() int {
	return len(s)
}

func (s FilePredictionsArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FilePredictionsArray) Less(i, j int) bool {
	return s[i].Prediction < s[j].Prediction
}
