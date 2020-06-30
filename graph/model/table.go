package model

import "time"

type TablePost struct {
	Email string
	Text string    `json:"text"`
	Time time.Time `json:"time"`

}
