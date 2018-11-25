package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/wtg/shuttletracker"
)

// StopService implements a mock of shuttletracker.StopService.
type StopService struct {
	mock.Mock
}

// CreateStop creates a Stop.
func (ss *StopService) CreateStop(stop *shuttletracker.Stop) error {
	args := ss.Called(stop)
	return args.Error(0)
}

// DeleteStop deletes a Stop.
func (ss *StopService) DeleteStop(id int64) error {
	args := ss.Called(id)
	return args.Error(0)
}

// Stops gets all stops.
func (ss *StopService) Stops() ([]*shuttletracker.Stop, error) {
	args := ss.Called()
	return args.Get(0).([]*shuttletracker.Stop), args.Error(1)
}

// CreateScheduleStop creates a ScheduleStop.
func (ss *StopService) CreateScheduleStop(stop *shuttletracker.ScheduleStop) error {
	args := ss.Called(stop)
	return args.Error(0)
}

// DeleteStop deletes a Stop.
func (ss *StopService) DeleteScheduleStop(id int64) error {
	args := ss.Called(id)
	return args.Error(0)
}

// AllScheduleStops gets all schedule stops.
func (ss *StopService) AllScheduleStops() ([]*shuttletracker.ScheduleStop, error) {
	args := ss.Called()
	return args.Get(0).([]*shuttletracker.ScheduleStop), args.Error(1)
}

// ScheduleStops gets all stops.
func (ss *StopService) ScheduleStops(schedule_id int64) ([]*shuttletracker.ScheduleStop, error) {
	args := ss.Called(schedule_id)
	return args.Get(0).([]*shuttletracker.ScheduleStop), args.Error(1)
}

// StopTimes gets all stop times.
// holding off on this until I'm sure the interface is right
func (ss *StopService) StopTimes(stop_id int64, schedule_id int64, departure_time int64, entries int64) ([]*shuttletracker.ScheduleStop, error) {
	args := ss.Called(stop_id, schedule_id, departure_time, entries)
	return args.Get(0).([]*shuttletracker.ScheduleStop), args.Error(1)
}

