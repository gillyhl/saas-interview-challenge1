package work

import (
	"encoding/json"
	"fmt"
	"log"
	"saas/exam"
	"saas/redis"
)

// Dispatcher to delegate marking when exam comes in
func Dispatcher(workChannel chan<- string) {
	go func() {
		client := redis.Client()
		ch := redis.Subscribe(client, "exam")
		log.Println("Listening for exams to mark from Redis")
		for msg := range ch {
			var exam exam.Exam
			json.Unmarshal([]byte(msg.Payload), &exam)
			exam.MarkExam(client)
		}
	}()

	go func() {
		client := redis.Client()
		ch := redis.Subscribe(client, "mark-paper")
		for msg := range ch {
			workChannel <- msg.Payload
		}
	}()

	go func() {
		client := redis.Client()
		ch := redis.Subscribe(client, "paper-marked")
		var markedPaper exam.MarkedPaper
		for msg := range ch {
			json.Unmarshal([]byte(msg.Payload), &markedPaper)
			key := fmt.Sprintf("%s:marked-papers#%s", markedPaper.ExamID, markedPaper.PaperID)
			log.Println(fmt.Sprintf("Paper %s from exam %s marked with a score of %d", markedPaper.PaperID, markedPaper.ExamID, markedPaper.Mark))
			client.Set(key, markedPaper, 0)
			client.RPush(fmt.Sprintf("%s:marked-papers", markedPaper.ExamID), markedPaper)
			log.Println(fmt.Sprintf("Paper %s from exam %s stored in redis with key %s", markedPaper.PaperID, markedPaper.ExamID, key))
		}
	}()
}
