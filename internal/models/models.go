package models

import "time"

type Movie struct {
	ID          int
	Title       string
	Description string
	Director    string
	ReleaseYear int
	Genre       string
	PosterURL   string
	TrailerURL  string
}

type Comment struct {
	Content   string
	CreatedAt time.Time
	Username  string
}

type Role struct {
	ID   int
	Name string
}

type User struct {
	ID       int
	Username string
	Password string
	RoleID   int
}
