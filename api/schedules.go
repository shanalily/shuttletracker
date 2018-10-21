package api

import (
	"net/http"

	// "github.com/wtg/shuttletracker"
	// "github.com/wtg/shuttletracker/log"
	// "github.com/wtg/shuttletracker/updater"
)

func (api *API) SchedulesHandler(w http.ResponseWriter, r *http.Request) {
	// get a specific schedule???

	// I should read from the spreadsheet instead of getting admin input

	// schedules, err := api.ms.Schedules()
	// if err != nil {
	// 	log.WithError(err).Error("unable to get schedules")
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// WriteJSON(w, schedules)
}