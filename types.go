package webhooks

import "time"

//LogPayload - schema for the payload.
type LogPayload struct {
	UserId int     `json:"user_id"`
	Total  float64 `json:"total"`
	Title  string  `json:"title"`
	Meta   struct {
		Logins []struct {
			Time time.Time `json:"time"`
			Ip   string    `json:"ip"`
		} `json:"logins"`
		PhoneNumbers struct {
			Home   string `json:"home"`
			Mobile string `json:"mobile"`
		} `json:"phone_numbers"`
	} `json:"meta"`
	Completed bool `json:"completed"`
}