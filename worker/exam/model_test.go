package exam

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"saas/redis"
	"encoding/json"
)

func TestMarkPaper(t *testing.T) {
	paperAnswers := Answers{
		Answer{
			QuestionNumber: 1,
			Answer:         "A",
		},
		Answer{
			QuestionNumber: 2,
			Answer:         "C",
		},
		Answer{
			QuestionNumber: 3,
			Answer:         "B",
		},
	}

	examAnswers := Answers{
		Answer{
			QuestionNumber: 1,
			Answer:         "A",
		},
		Answer{
			QuestionNumber: 2,
			Answer:         "D",
		},
		Answer{
			QuestionNumber: 3,
			Answer:         "B",
		},
	}

	paper := Paper{
		ID:      "paper0",
		Answers: paperAnswers,
	}

	markedPaper := paper.Mark(examAnswers, "exam-id")
	expectedMark := 2
	assert.Equal(t, markedPaper.Mark, expectedMark, "paper marks should be equal")
}

func TestMarkExam(t *testing.T) {
	exam := Exam{
		ID: "id",
		Name: "exam name",
		Answers: Answers{
			Answer{
				QuestionNumber: 1,
				Answer:         "A",
			},
			Answer{
				QuestionNumber: 2,
				Answer:         "D",
			},
			Answer{
				QuestionNumber: 3,
				Answer:         "B",
			},
		},
		Papers: Papers{
			Paper{
			ID: "paper1",
			},
			Paper{
				ID: "paper2",
			},
		},
	}
	client := redis.Client()
	exam.MarkExam(client)
	var answers Answers
	bytes, _ := client.Get("exam#id:answers").Bytes()
	json.Unmarshal(bytes, &answers)
	assert.Equal(t, exam.Answers, answers, "Answers should be the same")
	var paper Paper
	for i, _ := range exam.Papers {
		bytes, _ = client.RPop("exam#id:papers").Bytes()
		json.Unmarshal(bytes, &paper)
		assert.Equal(t, paper, exam.Papers[i], "Papers should be the same")
	}

	client.FlushAll()

}
