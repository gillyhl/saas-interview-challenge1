package work

import (
	"encoding/json"
	"saas/exam"
	"saas/redis"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	client := redis.Client()
	client.LPush("exam#id:papers", exam.Paper{
		ID: "p1",
		Answers: exam.Answers{
			exam.Answer{
				QuestionNumber: 1,
				Answer:         "A",
			},
			exam.Answer{
				QuestionNumber: 2,
				Answer:         "C",
			},
			exam.Answer{
				QuestionNumber: 3,
				Answer:         "B",
			},
		},
	})
	client.Set("exam#id:answers", exam.Answers{
		exam.Answer{
			QuestionNumber: 1,
			Answer:         "A",
		},
		exam.Answer{
			QuestionNumber: 2,
			Answer:         "D",
		},
		exam.Answer{
			QuestionNumber: 3,
			Answer:         "B",
		},
	}, 0)

	go doWork(1, "exam#id")
	redisCh := redis.Subscribe(client, "paper-marked")
	markedPaperMessage := <-redisCh
	var markedPaper exam.MarkedPaper
	json.Unmarshal([]byte(markedPaperMessage.Payload), &markedPaper)
	assert.Equal(t, markedPaper.Mark, 2, "Final mark should be equal")
	client.FlushAll()
}
