package main // for now

import (
	"fmt"
	"io/ioutil"
	"log"
	"encoding/json"

	"github.com/wtg/shuttletracker"
)

func main() {
	filename := "postgres/testschedulestops.json"
	data := parseFile(filename)

	scheduleStops := formatForDB(data)
	fmt.Println(scheduleStops)
}

func parseFile(filename string) []map[string]string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	var data []map[string]string
	if err := json.Unmarshal(b, &data); err != nil {
		log.Fatal(err.Error())
	}

	return data
}

func formatForDB(data []map[string]string) []shuttletracker.ScheduleStop {
	var scheduleStops []shuttletracker.ScheduleStop
	for _, stop := range data {
		for key, val := range stop {
			fmt.Println(key, ":", val)
		}
		// find the correct stop_id
		// find the correct schedule_id - if schedule doesn't exist then create schedule?
		// make new 
		stop := shuttletracker.ScheduleStop{
			ScheduleID: 0,
			StopID: 0,
			Time: 0,
		}
		scheduleStops = append(scheduleStops, stop)
	}
	return scheduleStops
}