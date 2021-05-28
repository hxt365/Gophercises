package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	exitChan   chan struct{}
	answerChan chan string
)

func main() {
	csvFilename := flag.String("csv", "problem.csv", "A csv file that contains problems and answers")
	timeLimit := flag.Int("limit", 3, "Time limit for each problem")
	shuffle := flag.Bool("shuffle", true, "Shuffle the quiz order each time it is run")

	flag.Parse()

	lines := readFile(*csvFilename, *shuffle)
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	answerChan = make(chan string)
	exitChan = make(chan struct{})

	correct := 0

problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem %d: %s = ", i+1, p.q)

		timer.Reset(time.Duration(*timeLimit) * time.Second)
		go getAnswer()

		select {
		case <-timer.C:
			fmt.Printf("\nTimeout\n")
			break problemLoop
		case answer := <-answerChan:
			if p.a == answer {
				correct++
			}
		case <-exitChan:
			break problemLoop
		}
	}

	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return problems
}

func readFile(filename string, shuffle bool) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		exit(fmt.Sprintf("Cannot open the CSV file at %s\n", filename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the CSV file\n"))
	}

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
	}

	return lines
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}

func getAnswer() {
	var answer string
	_, err := fmt.Scanf("%s\n", &answer)
	if err != nil {
		fmt.Printf("Input failed\n")
		exitChan <- struct{}{}
		return
	}
	answerChan <- answer
}
