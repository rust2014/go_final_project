package models

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"` // 20060102
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
