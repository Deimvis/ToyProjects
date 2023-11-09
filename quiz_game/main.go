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

func main() {
	csvFilepath := flag.String("problems", "problems.csv", "Path to csv file in the format of 'quesetion,answer'")
	timeLimit := flag.Int("time-limit", 30, "The time limit for the quiz (in seconds)")
	shuffle := flag.Bool("shuffle", false, "Shuffle problemset")
	flag.Parse()

	problems := extractProblems(*csvFilepath)
	if *shuffle {
		rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	}
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answer, ok := handleProblem(p, timer)
		if !ok {
			break
		}
		canonize := func(s string) string {
			return strings.ToLower(strings.TrimSpace(s))
		}
		if canonize(answer) == canonize(p.a) {
			correct++
		}
	}
	fmt.Printf("You scored %d/%d\n", correct, len(problems))
}

func handleProblem(p problem, timer *time.Timer) (string, bool) {
	answerCh := make(chan string)
	go func() {
		var answer string
		fmt.Scanf("%s\n", &answer)
		answerCh <- answer
	}()

	select {
	case <-timer.C:
		fmt.Println("\nYou ran out of time")
		return "", false
	case answer := <-answerCh:
		return answer, true
	}
}

func extractProblems(csvFilepath string) []problem {
	file, err := os.Open(csvFilepath)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", csvFilepath))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to read the CSV file: %s", csvFilepath))
	}
	return parseLines(lines)
}

func parseLines(lines [][]string) []problem {
	res := make([]problem, len(lines))
	for i, line := range lines {
		res[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return res
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

type problem struct {
	q string
	a string
}
