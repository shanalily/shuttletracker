package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/mock"
)

func TestSchedulesHandlerNoQuery(t *testing.T) {
	ms := &mock.ModelService{}
	// not passing yet because calling the function without parameters
	ms.StopService.On("AllScheduleStops").Return([]*shuttletracker.ScheduleStop{}, nil)

	api := API{
		ms: ms,
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("unable to create HTTP request: %s", err)
		return
	}

	api.SchedulesHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("got status code %d, expected 200", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("got Content-Type \"%s\", expected \"application/json\"", resp.Header.Get("Content-Type"))
	}

	var returnedScheduleStops []*shuttletracker.ScheduleStop
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&returnedScheduleStops)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	// do stuff here

	ms.VehicleService.AssertExpectations(t)
	ms.VehicleService.AssertNumberOfCalls(t, "AllScheduleStops", 1)
}

func TestSchedulesHandlerScheduleID(t *testing.T) {

}