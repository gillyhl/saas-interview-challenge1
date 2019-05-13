package work

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"saas/exam"
	"saas/redis"
	"saas/services"
)

func doWork(id int, examKey string) {
	client := redis.Client()
	defer client.Close()
	answersCh := make(chan interface{})
	paperCh := make(chan interface{})
	paperKey := fmt.Sprintf("%s:papers", examKey)
	waits := []services.Wait{
		services.Wait{
			ValueType: reflect.TypeOf(exam.Answers{}),
			Channel:   answersCh,
			Function: func(ch chan<- interface{}) {
				var answers exam.Answers
				bytes, e := client.Get(fmt.Sprintf("%s:answers", examKey)).Bytes()
				if e != nil {
					panic(e)
				}
				json.Unmarshal(bytes, &answers)
				ch <- answers
			},
		},
		services.Wait{
			ValueType: reflect.TypeOf(exam.Paper{}),
			Channel:   paperCh,
			Function: func(ch chan<- interface{}) {
				var paper exam.Paper
				log.Println(fmt.Sprintf("Worker %d marking paper from exam %s", id, examKey))
				bytes, e := client.RPop(paperKey).Bytes()
				if e != nil {
					panic(e)
				}
				json.Unmarshal(bytes, &paper)
				ch <- paper
			},
		},
	}

	values := services.WaitForAll(waits)
	fmt.Println(reflect.TypeOf(values))
	markedPaper := values[1].(exam.Paper).Mark(values[0].(exam.Answers), examKey)
	client.Publish("paper-marked", markedPaper)
}

// Worker function to mark exam paper
func Worker(id int, ch <-chan string) {
	for {
		examKey := <-ch
		doWork(id, examKey)
	}
}
