package models

import "time"

type Article struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Content     string    `db:"content" json:"content"`
	Category    string    `db:"category" json:"category"`
	CreatedDate time.Time `db:"created_date" json:"created_date"`
	UpdatedDate time.Time `db:"updated_date" json:"updated_date"`
	Status      string    `db:"status" json:"status"`
}