package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	reDashLeft  = regexp.MustCompile(`- `)
	reDashRight = regexp.MustCompile(` -`)
	reSpace     = regexp.MustCompile(`[ \n\t,.:;"]+`)
)

type words struct {
	word  string
	count int
}

func Top10(s string) []string {
	s0 := reDashLeft.ReplaceAll([]byte(s), []byte(" "))
	s1 := reDashRight.ReplaceAll(s0, []byte(" "))
	s2 := reSpace.ReplaceAll(s1, []byte(" "))
	s3 := strings.Split(string(s2), " ")

	m := make(map[string]int)
	for _, v := range s3 {
		m[strings.ToLower(v)]++
	}
	sStruct := make([]words, 0, len(m))

	for w, cnt := range m {
		sStruct = append(sStruct, words{w, cnt})
	}

	sort.Slice(sStruct, func(i, j int) bool {
		if sStruct[i].count == sStruct[j].count {
			return sStruct[i].word < sStruct[j].word
		}
		return sStruct[i].count > sStruct[j].count
	})

	result := make([]string, len(sStruct))
	for i, wrd := range sStruct {
		result[i] = wrd.word
	}
	l := len(result)

	if l > 9 {
		return result[:10]
	}
	return result[:0]
}
