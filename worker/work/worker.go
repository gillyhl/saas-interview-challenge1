package work

import (
	"encoding/json"
	"fmt"
	"log"
	"saas/exam"
	"saas/redis"
	"sync"
)

func doWork(id int, examKey string) {
	client := redis.Client()
	var paper exam.Paper
	var answers exam.Answers
	answersCh := make(chan exam.Answers)
	paperCh := make(chan exam.Paper)
	paperKey := fmt.Sprintf("%s:papers", examKey)
	var wg sync.WaitGroup
	wg.Add(2)
	go func(examKey, paperKey string, paperCh chan<- exam.Paper) {
		defer wg.Done()
		var paper exam.Paper
		log.Println(fmt.Sprintf("Worker %d marking paper from exam %s", id, examKey))
		bytes, e := client.RPop(paperKey).Bytes()
		if e != nil {
			panic(e)
		}
		json.Unmarshal(bytes, &paper)
		paperCh <- paper
	}(examKey, paperKey, paperCh)

	go func(examKey string, answersCh chan<- exam.Answers) {
		defer wg.Done()
		var answers exam.Answers
		bytes, e := client.Get(fmt.Sprintf("%s:answers", examKey)).Bytes()
		if e != nil {
			panic(e)
		}
		json.Unmarshal(bytes, &answers)
		answersCh <- answers
	}(examKey, answersCh)

	for i := 0; i < 2; i++ {
		select {
		case paper = <-paperCh:
		case answers = <-answersCh:
		}
	}
	wg.Wait()

	markedPaper := paper.Mark(answers, examKey)
	client.Publish("paper-marked", markedPaper)
}

// Worker function to mark exam paper
func Worker(id int, ch <-chan string) {
	for {
		examKey := <-ch
		doWork(id, examKey)
	}
}
