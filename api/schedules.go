package api

import (
	"net/http"
	"fmt"

	"github.com/wtg/shuttletracker/log"
	"github.com/360EntSecGroup-Skylar/excelize"
	"testing"
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
