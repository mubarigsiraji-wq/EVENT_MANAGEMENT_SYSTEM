package models

import "time"

type Event struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"EndTime"`
	Capacity    int       `json:"capacity"`
	Description string    `json:"Description"`
	ImgUrl      string    `json:"img_url"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}
