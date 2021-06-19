package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	duplicateVowel = true
	removeVowel    = false
	vowels         = "ueoaiUEOAI"
)

func randBool() bool {
	return rand.Intn(2) == 0
}

func isVowel(r rune) bool {
	return strings.ContainsRune(vowels, r)
}

func transform(s string, i int) string {
	var res string
	switch randBool() {
	case duplicateVowel:
		res = s[:i+1] + s[i:]
	case removeVowel:
		res = s[:i] + s[i+1:]
	}
	return res
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		text := s.Text()
		if randBool() {
			vI := -1
			for i, r := range text {
				if isVowel(r) {
					vI = i
				}
			}
			if vI >= 0 {
				fmt.Println(transform(text, vI))
			}
		}
	}
}
