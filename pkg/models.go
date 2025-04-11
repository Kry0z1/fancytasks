package tasks

import "time"

type User struct {
	Username          string             `json:"username"`
	HashedPassword    string             `json:"-"`
	BaseTasks         []BaseTask         `json:"base_tasks"`
	Events            []Event            `json:"events"`
	TasksWithDeadline []TaskWithDeadline `json:"tasks_with_deadline"`
	RepeatingTasks    []RepeatingTask    `json:"repeating_tasks"`
}

type BaseTask struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"desciption"`
	Done        bool   `json:"done"`
	Owner       User   `json:"owner"`
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
	Period time.Duration `json:"period"`
	Loop   int           `json:"loop"`
	Except []int         `json:"except"`
}
