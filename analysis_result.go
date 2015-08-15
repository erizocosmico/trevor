package trevor

import "sort"

type analysisResult struct {
	score        float64
	isExactMatch bool
	precedence   int
	name         string
	metadata     interface{}
}

func newAnalysisResult(score float64, isExactMatch bool, precedence int, name string, metadata interface{}) analysisResult {
	return analysisResult{score: score, isExactMatch: isExactMatch, precedence: precedence, name: name, metadata: metadata}
}

type byMatch []analysisResult

func (b byMatch) Len() int {
	return len(b)
}

func (b byMatch) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byMatch) Less(i, j int) bool {
	return b[i].isExactMatch && !b[j].isExactMatch
}

type byScore []analysisResult

func (b byScore) Len() int {
	return len(b)
}

func (b byScore) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byScore) Less(i, j int) bool {
	return b[i].score > b[j].score
}

type byPrecedence []analysisResult

func (b byPrecedence) Len() int {
	return len(b)
}

func (b byPrecedence) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byPrecedence) Less(i, j int) bool {
	return b[i].precedence > b[j].precedence
}

func sortAnalysisResults(results []analysisResult) {
	sort.Sort(byPrecedence(results))
	sort.Sort(byScore(results))
	sort.Sort(byMatch(results))
}

func getResults(plugins []Plugin, req *Request) []analysisResult {
	results := make([]analysisResult, len(plugins))
	for i, plugin := range plugins {
		score, metadata := plugin.Analyze(req)
		results[i] = newAnalysisResult(score.Score(), score.IsExactMatch(), plugin.Precedence(), plugin.Name(), metadata)
	}

	return results
}

func getBestResult(results []analysisResult) analysisResult {
	sortAnalysisResults(results)
	return results[0]
}
