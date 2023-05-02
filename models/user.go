package models

import (
	"time"
)

type UsersBody struct {
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"` //newly added for login system
}

type UsersGetBody struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type TodoBody struct {
	//UserId      int    `json:"userId" db:"user_id"`
	Task        string `json:"task" db:"task"`
	Description string `json:"description" db:"description"`
}

type TodoGetBody struct {
	Id          int       `json:"id" db:"id"`
	UserId      int       `json:"userId" db:"user_id"`
	Task        string    `json:"task" db:"task"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_At" db:"created_at"`
	DueDate     time.Time `json:"dueDate" db:"due_date"`
}

type TodoUpdateBody struct {
	Task        string `json:"task" db:"task"`
	Description string `json:"description" db:"description"`
}

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
