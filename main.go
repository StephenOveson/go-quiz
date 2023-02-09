package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type QuizQuestion struct {
	question string
	answer   string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "A csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds.")
	flag.Parse()
	f, err := os.Open(*csvFilename)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	quizQuestions := createQuizQuestion(data)
	score := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

quizLoop:
	for _, problem := range quizQuestions {

		fmt.Println(problem.question)

		answerChannel := make(chan string)

		go func() {
			var answerFromUser string
			fmt.Scanln(&answerFromUser)
			answerChannel <- answerFromUser
		}()

		select {
		case <-timer.C:
			fmt.Println("Time is up!")
			break quizLoop
		case answer := <-answerChannel:
			if answer == problem.answer {
				fmt.Println("Correct!")
				score++
			} else {
				fmt.Println("incorrect...")
			}
		}
	}
	fmt.Printf("Final score: %d", score)
}

func createQuizQuestion(data [][]string) []QuizQuestion {
	var quizQuestions []QuizQuestion
	for _, line := range data {
		quizQuestions = append(quizQuestions, QuizQuestion{question: line[0], answer: line[1]})
	}

	return quizQuestions
}
