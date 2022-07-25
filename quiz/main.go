package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func exit(msg string, err error) {
	log.Fatal(msg, err)
	os.Exit(1)
}

type problem struct {
	question string
	answer   string
}

func mkproblem(question string, answer string) problem {
	return problem{
		question: strings.TrimSpace(question),
		answer:   strings.TrimSpace(answer),
	}
}

func problemsFromCSV(filepath *string) []problem {
	f, err := os.Open(*filepath)
	if err != nil {
		exit("Unable to read the file "+*filepath, err)
	}
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		exit("", err)
	}
	problems := make([]problem, len(records))
	for i, record := range records {
		question, answer := record[0], record[1]
		problems[i] = mkproblem(question, answer)
	}
	return problems
}

func quiz(problems []problem, timeLimit int) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	reader := bufio.NewReader(os.Stdin)
	var score int = 0
problemloop:
	for index, problem := range problems {
		fmt.Printf("Question(%d/%d): %s?\n", index+1, len(problems), problem.question)
		textCh := make(chan string)
		go func() {
			fmt.Print(">> ")
			text, err := reader.ReadString('\n')
			if err != nil {
				exit("", err)
			}
			text = strings.TrimSpace(text)
			textCh <- text
		}()
		select {
		case <-timer.C:
			fmt.Println("\nTime Expired!")
			break problemloop
		case text := <-textCh:
			if text == problem.answer {
				score++
			}
		}
	}
	fmt.Printf("You scored: %d out of %d\n", score, len(problems))
	if score == len(problems) {
		fmt.Println("Congratulations 🎉🎊, you have solved all the problems")
	}
}

func main() {
	filePath := flag.String("filePath", "problems.csv", "csv path of the form 'question, answer'")
	timeLimit := flag.Int("limit", 2, "Time limit for our quiz")
	flag.Parse()
	problems := problemsFromCSV(filePath)
	quiz(problems, *timeLimit)
}
