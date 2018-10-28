package api

import (
	"net/http"

	"github.com/wtg/shuttletracker/log"
)

// can I specify which schedule I want in the url? like schedules?schedule_id=1
func (api *API) SchedulesHandler(w http.ResponseWriter, r *http.Request) {
	schedules, err := api.ms.ScheduleStops(1) // 1 for now
	if err != nil {
		log.WithError(err).Error("unable to get schedule")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, schedules)
}
