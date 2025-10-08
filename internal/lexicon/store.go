package lexicon

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ozeidan/fuzzy-patricia/patricia"
)

// Store holds the trie and statistics for the loaded lexicon.
type Store struct {
	mu   sync.RWMutex
	trie *patricia.Trie
	cnt  int
}

func NewStore() *Store {
	return &Store{trie: patricia.NewTrie()}
}

// LoadFromDir loads all .txt files from dir into the trie.
func (s *Store) LoadFromDir(dir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	newTrie := patricia.NewTrie()
	count := 0

	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.ToLower(filepath.Ext(info.Name())) != ".txt" {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		// Increase buffer for long lines
		buf := make([]byte, 0, 1024*64)
		scanner.Buffer(buf, 1024*1024)
		for scanner.Scan() {
			w := strings.TrimSpace(scanner.Text())
			if w == "" || strings.HasPrefix(w, "#") {
				continue
			}
			newTrie.Insert(patricia.Prefix(w), struct{}{})
			count++
		}
		return scanner.Err()
	})
	if walkErr != nil {
		return walkErr
	}
	if count == 0 {
		return errors.New("no entries loaded")
	}
	// Swap in
	s.trie = newTrie
	s.cnt = count
	return nil
}

func (s *Store) Stats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]interface{}{
		"count": s.cnt,
	}
}

// ForEachSubstringMatch visits any keys that contain the given substring.
// It uses the library's substring search.
func (s *Store) ForEachSubstringMatch(query string, visit func(word string) bool) {
	s.mu.RLock()
	tr := s.trie
	s.mu.RUnlock()
	if tr == nil || query == "" {
		return
	}
	// second argument is caseSensitive; we use false by default
	tr.VisitSubstring(patricia.Prefix(query), false, func(prefix patricia.Prefix, _ patricia.Item) error {
		// The library does not expose a public stop error in all versions; ignore early stop
		_ = visit(string(prefix))
		return nil
	})
}

// ForEachFuzzyMatch visits keys with fuzzy distance within maxDistance to query.
func (s *Store) ForEachFuzzyMatch(query string, maxDistance int, visit func(word string, distance int) bool) {
	s.mu.RLock()
	tr := s.trie
	s.mu.RUnlock()
	if tr == nil || query == "" {
		return
	}
	// signature in current lib: VisitFuzzy(prefix, caseSensitive bool, visitor)
	tr.VisitFuzzy(patricia.Prefix(query), false, func(prefix patricia.Prefix, _ patricia.Item, dist int) error {
		if dist <= maxDistance {
			_ = visit(string(prefix), dist)
		}
		return nil
	})
}

// HasAnyInText returns true if any lexicon word is a substring of the given text.
// It scans each rune offset and visits prefixes against the trie.
func (s *Store) HasAnyInText(text string) (bool, string) {
	s.mu.RLock()
	tr := s.trie
	s.mu.RUnlock()
	if tr == nil || text == "" {
		return false, ""
	}
	runes := []rune(text)
	n := len(runes)
	for i := 0; i < n; i++ {
		suffix := string(runes[i:])
		foundWord := ""
		tr.VisitPrefixes(patricia.Prefix(suffix), false, func(prefix patricia.Prefix, _ patricia.Item) error {
			if foundWord == "" {
				foundWord = string(prefix)
			}
			return nil
		})
		if foundWord != "" {
			return true, foundWord
		}
	}
	return false, ""
}
