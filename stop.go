package shuttletracker

import (
	"errors"
	"time"
)

// Stop is a place where vehicles frequently stop.
type Stop struct {
	ID        int64     `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`

	// Name and Description are pointers because they may be nil.
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type ScheduleStop struct {
	ID         int64 `json:"id"`
	ScheduleID int64 `json:"schedule_id"`
	StopID     int64 `json:"stop_id"`
	Time       int   `json:"time"`

	ScheduleName *string
	StopName     *string
}

// StopService is an interface for interacting with Stops.
type StopService interface {
	Stops() ([]*Stop, error)
	CreateStop(stop *Stop) error
	DeleteStop(id int64) error

	ScheduleStops(id int64) ([]*ScheduleStop, error)
	StopTimes(stop_id int64, schedule_id int64, departure_time int64, entries int64) ([]*ScheduleStop, error)
	CreateScheduleStop(stop *ScheduleStop) error
	DeleteScheduleStop(id int64) error
}

// ErrStopNotFound indicates that a Stop is not in the service.
var ErrStopNotFound = errors.New("Stop not found")
