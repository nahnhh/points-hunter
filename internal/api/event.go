package api

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type EventParams struct {
	Points uint8		`json:"points"`
	Name   string		`json:"ev_name"`
	Time   string		`json:"ev_time"`
	Place  string		`json:"ev_place"`
	Club   string		`json:"club"`
	Form	 Form		 `json:"form"`
}

type Form struct{
	Link     string `json:"link"`
	Deadline string `json:"deadline"`
}

type Error struct {
	Code         int
	ErrorMessage string
}

func GetEvent(token string, post string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  token,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf("Post này muốn đề cập tới sự kiện gì?\n%v", post)
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}