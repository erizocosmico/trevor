package trevor

import (
	"testing"
)

func Test_sortAnalysisResults(t *testing.T) {
	input, expected := getTestCases()

	sortAnalysisResults(input)

	for i, _ := range input {
		if input[i].name != expected[i].name {
			t.Errorf("expected %s to be at position %d, %s found", expected[i].name, i, input[i].name)
		}
	}
}

func Test_getBestResult(t *testing.T) {
	input, expected := getTestCases()

	result := getBestResult(input)
	if result.name != expected[0].name {
		t.Errorf("expected %s to be best result, %s found", expected[0].name, result.name)
	}
}

func getTestCases() ([]analysisResult, []analysisResult) {
	phrases := newAnalysisResult(1.0, false, 1, "phrases", nil)
	movies := newAnalysisResult(1.5, false, 1, "movies", nil)
	gifs := newAnalysisResult(1.5, true, 2, "gifs", nil)
	maps := newAnalysisResult(0.5, false, 1, "maps", nil)
	pictures := newAnalysisResult(1.5, true, 3, "pictures", nil)
	jokes := newAnalysisResult(1.5, false, 2, "jokes", nil)

	input := []analysisResult{
		phrases, movies, jokes, maps, pictures, gifs,
	}

	expected := []analysisResult{
		pictures, gifs, jokes, movies, phrases, maps,
	}

	return input, expected
}
