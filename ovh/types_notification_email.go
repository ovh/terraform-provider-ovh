package ovh

type NotificationEmail struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Date    string `json:"date"`
	Id      int64  `json:"id"`
}
