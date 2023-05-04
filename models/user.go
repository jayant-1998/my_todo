package models

import (
	"time"
)

//var JwtKey = []byte("golang_my_todo")

type Register struct {
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"` //newly added for login system
}

type RetrieveUserInfo struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Task struct {
	Task        string `json:"task" db:"task"`
	Description string `json:"description" db:"description"`
}

type AllTask struct {
	Id          int       `json:"taskId" db:"id"`
	Task        string    `json:"taskName" db:"task"`
	Description string    `json:"description" db:"description"`
	Complete    bool      `json:"completed" db:"is_completed"`
	DueDate     time.Time `json:"toBeCompletedBy" db:"due_date"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type UpdateTask struct {
	Id          int    `json:"id" db:"id"`
	Task        string `json:"task" db:"task"`
	Description string `json:"description" db:"description"`
}

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPass struct {
	UserID   int    `db:"user_id"`
	Password string `db:"password"`
}

//type JwtClaims struct {
//	Session  string    `json:"session"`
//	Auth     bool      `json:"auth"`
//	Username string    `json:"username"`
//	Exp      time.Time `json:"exp"`
//	jwt.StandardClaims
//}
