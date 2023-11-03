package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"quiz/model"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type OpenAiInterface interface {
	GenerateQuestions(newPromt model.OpenAiReq) []model.Questions
}

type OpenAi struct {
	Key string
}

func NewOpenAi(key string) OpenAiInterface {
	return &OpenAi{
		Key: key,
	}
}

func (op *OpenAi) GenerateQuestions(newPromt model.OpenAiReq) []model.Questions {

	format := `[{"question":"apa nama ibu kota Indonesia?","options":[{"value":"semarang","is_right":false},{"value":"bekasi","is_right":false},{"value":"balikpapan","is_right":false},{"value":"jakarta","is_right":true}},{"question":"apa nama presiden Indonesia sekarang?","options":[{"value":"Prabowo","is_right":false},{"value":"Anis","is_right":false},{"value":"Jokowi","is_right":true},{"value":"Ganjar","is_right":false}]}]`

	client := openai.NewClient(op.Key) 
	resp, err := client.CreateCompletion(
		context.Background(),
		openai.CompletionRequest{
			Model: openai.GPT3TextDavinci003,
			Prompt: fmt.Sprintf("kamu akan diperintahakan untuk membuat beberapa soal minimal 10 soal sesuai dengan description yang diberikan, dan kamu harus mengikuti format json seperti ini %s\n dan kamu hanya menampilkan format question dan options seperti dicontohkan dan kamu akan membuat soal yang semua quiz_id : %d, dan sesuai dengan deskirsi berikut: %s", format, newPromt.Quiz_id, newPromt.Description),
			MaxTokens: 3000,
		},
	)

	if err != nil {
		logrus.Error("CreateCompletion error: ", err.Error())
		return nil
	}

	recommendation := resp.Choices[0].Text

    cleanedJSON := strings.ReplaceAll(recommendation, "\n\n", "")
    cleanedJSON = strings.ReplaceAll(cleanedJSON, "\n\t", "")
	cleanedJSON = strings.ReplaceAll(cleanedJSON, "\n]]", "")
	cleanedJSON = strings.ReplaceAll(cleanedJSON, "[[\n", "")
    cleanedJSON = strings.TrimSpace(cleanedJSON)

	var questions []model.Questions
    if err := json.Unmarshal([]byte(cleanedJSON), &questions); err != nil {
        return nil
    }

	return questions
}

