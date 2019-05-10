package exam

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

// MarkedPaperAnswer data type
type MarkedPaperAnswer struct {
	Answer  `json:"answer"`
	Correct bool `json:"correct"`
}

// MarkedPaper struct
type MarkedPaper struct {
	ExamID             string              `json:"examId"`
	PaperID            string              `json:"paperId"`
	MarkedPaperAnswers []MarkedPaperAnswer `json:"answers"`
	Mark               int                 `json:"mark"`
}

// Answer type with question number and multiple option answer
type Answer struct {
	QuestionNumber int    `json:"questionNumber"`
	Answer         string `json:"answer"`
}

// Answers slice
type Answers []Answer

// Paper type
type Paper struct {
	ID      string  `json:"id"`
	Answers Answers `json:"answers"`
}

// Papers type
type Papers []Paper

// Exam with correct answers
type Exam struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Answers Answers `json:"answers"`
	Papers  Papers  `json:"papers"`
}

// RedisKey for a given exam for storage
func (e Exam) RedisKey() string {
	return fmt.Sprintf("exam#%s", e.ID)
}

// MarshalBinary of exam
func (e Exam) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

// MarshalBinary of answers
func (a Answers) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

// MarshalBinary of papers
func (p Paper) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

// MarshalBinary of marked paper
func (mp MarkedPaper) MarshalBinary() ([]byte, error) {
	return json.Marshal(mp)
}

// Mark a paper against a slice of answers
func (p Paper) Mark(answers Answers, examID string) *MarkedPaper {
	markedAnswers := make([]MarkedPaperAnswer, len(answers))
	mark := 0
	for i, answer := range answers {
		correct := answer.Answer == p.Answers[i].Answer
		if correct {
			mark++
		}
		markedAnswers[i] = MarkedPaperAnswer{
			Answer:  answer,
			Correct: correct,
		}
	}

	return &MarkedPaper{
		ExamID:             examID,
		PaperID:            p.ID,
		Mark:               mark,
		MarkedPaperAnswers: markedAnswers,
	}
}

// MarkExam by setting answers in redis and pushing papers
func (e *Exam) MarkExam(client *redis.Client) {
	client.Set(fmt.Sprintf("%s:answers", e.RedisKey()), e.Answers, 0)
	for _, paper := range e.Papers {
		key := fmt.Sprintf("%s:papers", e.RedisKey())
		client.LPush(key, paper)
		client.Publish("mark-paper", e.RedisKey())
	}
}
