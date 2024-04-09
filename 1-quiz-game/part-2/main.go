package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	shufflePtr := flag.Bool("shuffle", false, "Shuffle the questions")
	timeLimitSecondsPtr := flag.Int("timelimit", 30, "Time limit in seconds")
	flag.Parse()

	var filename string
	if len(flag.Args()) > 0 {
		filename = flag.Args()[0]
	} else {
		filename = "problems.csv"
	}

	err := startQuiz(filename, *shufflePtr, *timeLimitSecondsPtr)
	if err != nil {
		fmt.Println(err)
	}
}

func startQuiz(file string, shuffle bool, timeLimitSeconds int) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("open quiz: %w", err)
	}
	defer f.Close()

	q, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return fmt.Errorf("read questions: %w", err)
	}

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(q), func(i, j int) { q[i], q[j] = q[j], q[i] })
	}

	fmt.Printf("Once you press enter, you will have %d seconds to answer %d questions.\n", timeLimitSeconds, len(q))
	fmt.Scanln()

	var (
		correctCount = 0
		errc         = make(chan error)
		isCorrectc   = make(chan bool)
		timer        = time.NewTimer(time.Duration(timeLimitSeconds) * time.Second)
	)

Quiz:
	for _, question := range q {
		go askQuestion(question, isCorrectc, errc)

		select {
		case <-timer.C:
			fmt.Println("out of time!")
			break Quiz
		case isCorrect := <-isCorrectc:
			if isCorrect {
				fmt.Println("correct!")
				correctCount++
			} else {
				fmt.Println("incorrect😔")
			}
		case err := <-errc:
			return err
		}
	}
	fmt.Printf("Quiz complete! Score: %d/%d\n", correctCount, len(q))
	return nil
}

func askQuestion(question []string, isCorrectc chan bool, errc chan error) {
	if len(question) < 2 {
		errc <- errors.New("question %v: missing field")
	}
	fmt.Println(question[0])

	var answer string
	fmt.Scanln(&answer)

	isCorrectc <- strings.TrimSpace(strings.ToLower(answer)) == strings.TrimSpace(strings.ToLower(question[1]))
}
