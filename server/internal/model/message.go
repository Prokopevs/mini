package model

type Message struct {
	Id      int    `json:"id" db:"id"`
	Message string `json:"message" db:"message"`
	Status  string `json:"status" db:"status"`
}

type MessageCreate struct {
	Message string `json:"message"`
}
