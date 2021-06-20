package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const toReplace = "*"

var transforms = []string{
	toReplace,
	toReplace + "app",
	toReplace + "site",
	toReplace + "time",
	toReplace + "hq",
	"get" + toReplace,
	"go" + toReplace,
	"lets" + toReplace,
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, toReplace, s.Text(), -1))
	}
}
