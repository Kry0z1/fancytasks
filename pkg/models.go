package tasks

import "time"

type UserBase struct {
	Username string
}

type UserRegister struct {
	UserBase
	Password string
}

type UserStored struct {
	UserBase
	HashedPassword    string
	Events            []Event
	TasksWithDeadline []TaskWithDeadline
	RepeatingTasks    []RepeatingTask
}

type BaseTask struct {
	ID          int
	Title       string
	Description string
	Done        bool
	Owner       UserBase
}

type Event struct {
	BaseTask
	StartsAt time.Time
	EndsAt   time.Time
}

type TaskWithDeadline struct {
	BaseTask
	Deadline time.Time
}

// For every `Loop` tasks those at places in Except are considered turned off
type RepeatingTask struct {
	Event
	Period time.Time
	Loop   int
	Except []int
}
