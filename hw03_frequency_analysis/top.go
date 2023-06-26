package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reDashLeft = regexp.MustCompile(`- `)
var reDashRight = regexp.MustCompile(` -`)
var reSpace = regexp.MustCompile(`[ \n\t,.:;"]+`)

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
	var sStruct []words

	for w, cnt := range m {
		sStruct = append(sStruct, words{w, cnt})
	}

	sort.Slice(sStruct, func(i, j int) bool {
		if sStruct[i].count == sStruct[j].count {
			return sStruct[i].word < sStruct[j].word
		} else {
			return sStruct[i].count > sStruct[j].count
		}
	})

	var result []string
	for _, wrd := range sStruct {
		result = append(result, wrd.word)
	}
	l := len(result)

	if l > 9 {
		return result[:10]
	} else {
		return result[:0]
	}
}
