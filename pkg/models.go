package tasks

import "time"

type User struct {
	Username          string
	HashedPassword    string
	BaseTasks         []BaseTask
	Events            []Event
	TasksWithDeadline []TaskWithDeadline
	RepeatingTasks    []RepeatingTask
}

type BaseTask struct {
	ID          int
	Title       string
	Description string
	Done        bool
	Owner       User
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
	Period time.Duration
	Loop   int
	Except []int
}
