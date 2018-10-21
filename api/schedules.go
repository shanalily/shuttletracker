package api

import (
	// "encoding/json"
	"net/http"

	// "github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
)

func (api *API) SchedulesHandler(w http.ResponseWriter, r *http.Request) {
	schedules, err := api.ms.Routes()
	if err != nil {
		log.WithError(err).Error("unable to get schedules")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, schedules)
}
