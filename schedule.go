package shuttletracker

import (
	"errors"
	"time"
)

// Schedule represents a set of times and places in a shuttle route schedule
type Schedule struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	StopIDs     []int64   `json:"stop_ids"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// I need to have the name/time of the stops
type ScheduleStop struct {
	ID          int64     `json:"id"`
	ScheduleID 	int64     `json:"schedule_id"`
	ArrivalTime time.Time `json:"arrival_time"`
	Name 		string    `json:"name"`
	Description *string   `json:"description"`
	Created   time.Time   `json:"created"`
	Updated   time.Time   `json:"updated"`
}

// ScheduleService is an interface for interacting with Schedules.
type ScheduleService interface {
	Schedule(id int64) (*Schedule, error)
	Schedules() ([]*Schedule, error)
	CreateSchedule(schedule *Schedule) error
	DeleteSchedule(id int64) error
	ModifySchedule(schedule *Schedule) error

	ScheduleStop(id int64) (*ScheduleStop, error)
	ScheduleStops() ([]*ScheduleStop, error)
	CreateScheduleStop(schedule *Schedule) error
	DeleteScheduleStop(id int64) error
	ModifyScheduleStop(schedule *Schedule) error
}

// ErrScheduleNotFound indicates that a Schedule is not in the service.
var ErrScheduleNotFound = errors.New("Schedule not found")
var ErrScheduleStopNotFound = errors.New("ScheduleStop not found")