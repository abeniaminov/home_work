package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

func GetDomainStatFastAndFrugal(r io.Reader, domain string) (DomainStat, error) {
	var key string
	var email string
	searchStr := "." + domain
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)

	for scanner.Scan() {
		email = fastjson.GetString(scanner.Bytes(), "Email")
		if strings.HasSuffix(email, searchStr) {
			key = strings.ToLower(strings.SplitN(email, "@", 2)[1])
			result[key]++
		}
	}

	return result, nil
}
