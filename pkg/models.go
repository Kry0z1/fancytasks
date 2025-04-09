package tasks

import "time"

type UserBase struct {
	Username string `json:"username"`
}

type UserRegister struct {
	UserBase
	Password string `json:"password"`
}

type UserStored struct {
	UserBase
	HashedPassword    string
	Events            []Event
	TasksWithDeadline []TaskWithDeadline
	RepeatingTasks    []RepeatingTask
}

type BaseTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
	Owner       UserBase
}

type Event struct {
	BaseTask
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

type TaskWithDeadline struct {
	BaseTask
	Deadline time.Time `json:"deadline"`
}

// For every `Loop` tasks those at places in Except are considered turned off
type RepeatingTask struct {
	Event
	Period time.Time `json:"period"`
	Loop   int       `json:"loop"`
	Except []int     `json:"except"`
}
