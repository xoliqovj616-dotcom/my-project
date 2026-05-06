package model

type Todo struct {
	Id        int    `json:"id"`
	Work      string `json:"work"`
	Time      string `json:"time"`
	Completed bool   `json:"completed"`
	User_id   int    `json:"user_id"`
}
