package util

import (
	"math/rand"
	"regexp"
	"time"
)

var tokenContainer = make(map[string]struct{})
var tokenLength = 10

func ValidateToken(token string) bool {
	if _, ok := tokenContainer[token]; ok {
		delete(tokenContainer, token)
		return true
	}
	return false
}

func GetToken() string {
	str := "123456789ABCDEFGHIJKLMNPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < tokenLength; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	ok1, _ := regexp.MatchString(".[1|2|3|4|5|6|7|8|9]", string(result))
	ok2, _ := regexp.MatchString(".[Z|X|C|V|B|N|M|A|S|D|F|G|H|J|K|L|Q|W|E|R|T|Y|U|I|P]", string(result))
	if ok1 && ok2 {
		tokenContainer[string(result)] = struct{}{}
		return string(result)
	} else {
		return GetToken()
	}
}
