package detect

import (
	"sort"
	"strings"
	"unicode/utf8"

	"sensitive-lexicon/internal/lexicon"
)

type FuzzyConfig struct {
	MinNgramLen int
	MaxNgramLen int
	MaxDistance int
}

type DetectRequest struct {
	Text string `json:"text"`
	// If true, enable fuzzy detection on n-grams within the text
	EnableFuzzy bool `json:"enable_fuzzy"`
}

type Match struct {
	Word     string `json:"word"`
	Type     string `json:"type"` // substring | fuzzy
	Distance int    `json:"distance,omitempty"`
}

type DetectResponse struct {
	Hits []Match `json:"hits"`
}

type ContainsRequest struct {
	Text string `json:"text"`
}

type ContainsResponse struct {
	Contains bool   `json:"contains"`
	Word     string `json:"word,omitempty"`
}

type Service struct {
	store    *lexicon.Store
	fuzzyCfg FuzzyConfig
}

func NewService(store *lexicon.Store) *Service {
	return &Service{store: store, fuzzyCfg: FuzzyConfig{MinNgramLen: 2, MaxNgramLen: 10, MaxDistance: 1}}
}

func (s *Service) SetFuzzyConfig(cfg FuzzyConfig) {
	s.fuzzyCfg = cfg
}

func (s *Service) Detect(req DetectRequest) DetectResponse {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return DetectResponse{}
	}

	unique := make(map[string]Match)

	// Substring hits: for each codepoint window from input, find lexicon entries containing it
	s.store.ForEachSubstringMatch(text, func(word string) bool {
		unique[word] = Match{Word: word, Type: "substring"}
		return true
	})

	if req.EnableFuzzy {
		for _, token := range generateNgrams(text, s.fuzzyCfg.MinNgramLen, s.fuzzyCfg.MaxNgramLen) {
			s.store.ForEachFuzzyMatch(token, s.fuzzyCfg.MaxDistance, func(word string, d int) bool {
				if old, ok := unique[word]; ok {
					if old.Type == "substring" && d == 0 {
						return true
					}
				}
				unique[word] = Match{Word: word, Type: ternary(d == 0, "substring", "fuzzy"), Distance: d}
				return true
			})
		}
	}

	res := DetectResponse{Hits: make([]Match, 0, len(unique))}
	for _, v := range unique {
		res.Hits = append(res.Hits, v)
	}
	sort.Slice(res.Hits, func(i, j int) bool {
		if res.Hits[i].Type == res.Hits[j].Type {
			if res.Hits[i].Distance == res.Hits[j].Distance {
				return res.Hits[i].Word < res.Hits[j].Word
			}
			return res.Hits[i].Distance < res.Hits[j].Distance
		}
		return res.Hits[i].Type < res.Hits[j].Type
	})
	return res
}

func (s *Service) Contains(req ContainsRequest) ContainsResponse {
	ok, w := s.store.HasAnyInText(strings.TrimSpace(req.Text))
	return ContainsResponse{Contains: ok, Word: w}
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func generateNgrams(text string, minLen, maxLen int) []string {
	if minLen < 1 {
		minLen = 1
	}
	if maxLen < minLen {
		maxLen = minLen
	}
	// Work on rune boundaries for CJK safety
	runes := []rune(text)
	n := len(runes)
	var out []string
	for i := 0; i < n; i++ {
		for l := minLen; l <= maxLen && i+l <= n; l++ {
			out = append(out, string(runes[i:i+l]))
		}
	}
	return dedupStrings(out)
}

func dedupStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// Guard for unused import warning if utf8 not referenced elsewhere
var _ = utf8.RuneCountInString
