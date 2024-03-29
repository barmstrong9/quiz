package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

// to change time limit, go run main.go -limit=int
func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv in the format 'question, answer'")
	timeLimit := flag.Int("limit", 30, "the time left for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed To Open The CSV File: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed To Parse The Given CSV File")
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nYou got %d out of %d correct\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.answer {
				correct++
			}

		}
	}
	fmt.Printf("You got %d out of %d correct", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	returnValue := make([]problem, len(lines))
	for i, line := range lines {
		returnValue[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return returnValue
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
