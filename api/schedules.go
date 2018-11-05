package api

import (
	// "encoding/json"
	"net/http"
	"fmt"

  "github.com/Luxurioust/excelize"
	// "github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
	"testing"
)

func (api *API) SchedulesHandler(w http.ResponseWriter, r *http.Request) {
	schedules, err := api.ms.ScheduleStops()
	if err != nil {
		log.WithError(err).Error("unable to get schedules")
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
