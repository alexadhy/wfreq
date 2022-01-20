package model

// Pair is just a wrapper for normal KV store
type Pair struct {
	Key   string `json:"word"`
	Value int    `json:"occurence"`
}

// PairList contains a slice of Pair, it satisfies the Sort interface
type PairList []Pair

// Len satisfies Sort interface
func (p PairList) Len() int { return len(p) }

// Less satisfies Sort interface
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// Swap satisfies Sort interface
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// RequestBody is used for sending request body
type RequestBody struct {
	MinWordLength int    `json:"min_word_length"`
	Content       string `json:"content"`
}
