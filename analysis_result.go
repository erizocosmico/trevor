package trevor

import "sort"

type analysisResult struct {
	score        float64
	isExactMatch bool
	precedence   int
	name         string
}

func newAnalysisResult(score float64, isExactMatch bool, precedence int, name string) analysisResult {
	return analysisResult{score: score, isExactMatch: isExactMatch, precedence: precedence, name: name}
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

func getBestResult(results []analysisResult) analysisResult {
	sortAnalysisResults(results)
	return results[0]
}
