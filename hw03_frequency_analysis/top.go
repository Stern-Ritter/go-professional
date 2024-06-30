package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var regexpPunctuation = regexp.MustCompile(`[!"#$%&'()*+,./:;<=>?@[\\\]^_{|}~]|(\s-\s)`)

type entry[K comparable, V comparable] struct {
	key   K
	value V
}

func Top10(str string) []string {
	topCount := 10
	wordsCounter := make(map[string]int)

	formatedStr := regexpPunctuation.ReplaceAllString(strings.ToLower(str), " ")

	for _, word := range strings.Fields(formatedStr) {
		wordsCounter[word]++
	}

	entries := make([]entry[string, int], 0, len(wordsCounter))
	for word, count := range wordsCounter {
		entries = append(entries, entry[string, int]{word, count})
	}

	sortByValueDescAndNameAsc := func(i, j int) bool {
		if entries[i].value == entries[j].value {
			return entries[i].key < entries[j].key
		}
		return entries[i].value > entries[j].value
	}

	sort.Slice(entries, sortByValueDescAndNameAsc)

	top := make([]string, 0, topCount)
	for i := 0; i < topCount && i < len(entries); i++ {
		top = append(top, entries[i].key)
	}

	return top
}
