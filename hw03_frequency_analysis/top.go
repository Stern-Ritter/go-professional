package hw03frequencyanalysis

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const topWordsLimit = 10

var punctuationMarks = string([]rune{
	'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
	':', ';', '<', '=', '>', '?', '@', '[', '\\', ']', '^', '_', '`', '{', '|', '}', '~',
})

var replacePattern = regexp.MustCompile(fmt.Sprintf(`(^|\s)[%s]|[%s]($|\s)|(^|\s)[%s]($|\s)`, punctuationMarks,
	punctuationMarks, punctuationMarks))

type entry[K comparable, V comparable] struct {
	key   K
	value V
}

func Top10(str string) []string {
	dictionary := make(map[string]int)

	formatedStr := replacePattern.ReplaceAllString(strings.ToLower(str), " ")
	for _, word := range strings.Fields(formatedStr) {
		dictionary[word]++
	}

	entries := make([]entry[string, int], 0, len(dictionary))

	for word, count := range dictionary {
		entries = append(entries, entry[string, int]{word, count})
	}

	sortByValueDescAndNameAsc := func(i, j int) bool {
		if entries[i].value == entries[j].value {
			return entries[i].key < entries[j].key
		}
		return entries[i].value > entries[j].value
	}

	sort.Slice(entries, sortByValueDescAndNameAsc)

	limit := min(len(entries), topWordsLimit)
	result := make([]string, limit)
	for i := range limit {
		result[i] = entries[i].key
	}

	return result
}
