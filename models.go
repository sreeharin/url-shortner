package main

type URL struct {
	Original  string `json:"original"`
	Shortened string `json:"shortened"`
}

type UrlDB struct {
	ID uint `gorm:"primaryKey"`
	URL
}
