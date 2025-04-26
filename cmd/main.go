package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	google "github.com/nahnhh/points-hunter/google"
	api "github.com/nahnhh/points-hunter/internal/api"
	preprocess "github.com/nahnhh/points-hunter/preprocess"
)

var	testMessage string = `Sự kiện: Workshop Lập Trình
Thời gian: 14h00 ngày 20/12/2023
Địa điểm: Phòng A1.01
CLB: CLB Lập Trình
Quyền lợi: Nhận 3 điểm đoàn
Link đăng ký: https://forms.gle/abc123
Hạn đăng ký: 26/04/2025`

var stripped = preprocess.FilterText(testMessage)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	email := os.Getenv("EMAIL")
	pass := os.Getenv("PASS")
	token := os.Getenv("GEMINI_API_KEY")
	refreshToken := os.Getenv("GOOGLE_REFRESH_TOKEN")

	fmt.Printf("Email: %v, Pass: %v\n", email, pass)
	fmt.Printf("\nPost: %v\n", stripped)

	// Get event details from Gemini
	response, err := api.GetEvent(token, stripped)
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Printf("Response: %v\n", response)

	// Initialize Google Calendar service
	calendarService := google.NewGoogleCalendarService()
	
	// Get access token
	accessToken, err := calendarService.GetAccessToken(refreshToken)
	if err != nil {
		log.Fatal("Error getting access token:", err)
	}

	// Get events for next 14 days
	now := time.Now()
	twoWeeksLater := now.Add(14 * 24 * time.Hour)
	events, err := calendarService.GetEvents(context.Background(), accessToken, now, twoWeeksLater)
	if err != nil {
		log.Fatal("Error getting events:", err)
	}

	// Print events
	fmt.Println("\nUpcoming events:")
	for _, event := range events {
		fmt.Printf("- %s (%s)\n", event.Summary, event.Start.DateTime)
	}
}