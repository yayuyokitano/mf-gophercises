package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func main() {
	shufflePtr := flag.Bool("shuffle", false, "Shuffle the questions")
	flag.Parse()

	var filename string
	if len(flag.Args()) > 0 {
		filename = flag.Args()[0]
	} else {
		filename = "problems.csv"
	}

	err := startQuiz(filename, *shufflePtr)
	if err != nil {
		fmt.Println(err)
	}
}

func startQuiz(file string, shuffle bool) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("open quiz: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	q, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return fmt.Errorf("read questions: %w", err)
	}

	if shuffle {
		rand.Shuffle(len(q), func(i, j int) { q[i], q[j] = q[j], q[i] })
	}

	var (
		questionCount = 0
		correctCount  = 0
	)

	for _, question := range q {
		isCorrect, err := askQuestion(question)
		if err != nil {
			return err
		}

		questionCount++
		if isCorrect {
			fmt.Println("correct!")
			correctCount++
		} else {
			fmt.Println("incorrect😔")
		}
	}
	fmt.Printf("Quiz complete! Score: %d/%d\n", correctCount, questionCount)
	return nil
}

func askQuestion(question []string) (bool, error) {
	if len(question) < 2 {
		return false, errors.New("question %v: missing field")
	}
	fmt.Println(question[0])

	var answer string
	_, err := fmt.Scanln(&answer)
	if err != nil {
		return false, fmt.Errorf("scan answer: %w", err)
	}
	return strings.TrimSpace(strings.ToLower(answer)) == strings.TrimSpace(strings.ToLower(question[1])), nil
}
