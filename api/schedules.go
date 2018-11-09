package api

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/wtg/shuttletracker/log"
	"github.com/360EntSecGroup-Skylar/excelize"
	"testing"
)

// can I specify which schedule I want in the url? like schedules?schedule_id=1
func (api *API) SchedulesHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	base := 10
	int_type := 64
	var err error
	var schedule_id, stop_id, departure_time, entries int64
	if val, ok := params["schedule_id"]; ok {
		schedule_id, err = strconv.ParseInt(val[0], base, int_type)
		if err != nil {
			log.WithError(err).Error("invalid schedule_id")
			return
		}
		fmt.Println(schedule_id)
	}

	if val, ok := params["stop_id"]; ok {
		// check that departure time and number of entries are also set
		stop_id, err = strconv.ParseInt(val[0], base, int_type)
		if err != nil {
			log.WithError(err).Error("invalid stop_id")
			return
		}
		// should these parameters be optional?
		// default departure time could be now, default limit could 1 or everything (somehow)
		if val, ok := params["departure_time"]; ok {
			departure_time, err = strconv.ParseInt(val[0], base, int_type)
		} else {
			log.WithError(err).Error("invalid departure_time")
			return
		}
		if val, ok := params["entries"]; ok {
			entries, err = strconv.ParseInt(val[0], base, int_type)
		} else {
			log.WithError(err).Error("invalid entries")
			return
		}
		stops, err := api.ms.StopTimes(stop_id, schedule_id, departure_time, entries)
		if err != nil {
			log.WithError(err).Error("unable to get stops")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		WriteJSON(w, stops)
	} else {
		schedules, err := api.ms.ScheduleStops(schedule_id) // 1 for now
		if err != nil {
			log.WithError(err).Error("unable to get schedule")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		WriteJSON(w, schedules)
	}
}

func testParsing(t *testing.T){
	xlsx, err := excelize.OpenFile("2018-19 M-Fri Campus Schedule East-West shuttle schedule.xlsx")
	    if err != nil {
	        fmt.Println(err)
	        return
	    }

	    // Get all the rows in the Sheet1.
	    rows := xlsx.GetRows("East Master")
	    for _, row := range rows {
	        for _, colCell := range row {
	            fmt.Print(colCell, "\t")
	        }
	        fmt.Println()
	    }
}
