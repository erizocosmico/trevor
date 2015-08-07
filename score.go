package trevor

type Score interface {
	// IsExactMatch returns a boolean value that will be true if the text was an exact match for a rule in the plugin.
	IsExactMatch() bool

	// Score returns the score
	Score() float64
}

type _score struct {
	exactMatch bool
	score      float64
}

// NewScore creates a new Score instance
func NewScore(score float64, exactMatch bool) Score {
	return &_score{exactMatch: exactMatch, score: score}
}

func (s *_score) IsExactMatch() bool {
	return s.exactMatch
}

func (s *_score) Score() float64 {
	return s.score
}
