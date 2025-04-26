package api

import (
	"context"
	"encoding/json"

	"google.golang.org/genai"
)

type EventParams struct {
	Points uint8		`json:"points"`
	Name   string		`json:"ev_name"`
	Time   string		`json:"ev_time"`
	Place  string		`json:"ev_place"`
	Club   string		`json:"club"`
	Form	Form		 `json:"form"`
}

type Form struct{
	Link string `json:"link"`
	Deadline string `json:"deadline"`
}

type Error struct {
	Code         int
	ErrorMessage string
}

var	testMessage = `Sự kiện: Workshop Lập Trình
Thời gian: 14h00 ngày 20/12/2023
Địa điểm: Phòng A1.01
CLB: CLB Lập Trình
Quyền lợi: Nhận 3 điểm đoàn
Link đăng ký: https://forms.gle/abc123
Hạn đăng ký: 26/04/2025`

func GetEvent(token string) (*EventParams, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  token,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(testMessage),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Parse Gemini's response into EventParams
	var event EventParams
	err = json.Unmarshal([]byte(result.Text()), &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}