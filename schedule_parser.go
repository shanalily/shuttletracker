package main

import (
    "encoding/json"
    "fmt"
    "github.com/tealeg/xlsx"
    "io/ioutil"
    "strings"
)

func main() {
    excelFileName := "west_tester.xlsx"
    xlFile, err := xlsx.OpenFile(excelFileName)
    var stops []string //declare variable type for json
    var route string //route, ex.West
    var day string //days, ex.Weekend
    var stop_name string //stop_name, ex.Blitman
    var time string //time, ex. 7:00a

    //for loop goes through each cell in excel roqw by row
    for _, sheet := range xlFile.Sheets {         
        for _, row := range sheet.Rows { 
            for _, cell := range row.Cells{
                text := cell.String()

                //only consider non-empty cells
                if text != "" {

                    //west campus
                    if strings.Contains(text, "West--")  {
                        route = cell.String()
                        route = "West" //declare route
                        sub_str := text[len("West--"):len(text)]
                        day = cell.String() 
                        day = sub_str //declare day
                        //fmt.Printf("%s\n", day)
                    }

                    //east campus
                    if strings.Contains(text, "East--") {
                        route = cell.String()
                        route = "East" //declare route
                        sub_str := text[len("East--"):len(text)]
                        day = cell.String()
                        day = sub_str //declare day
                        //fmt.Printf("%s\n", day)
                    }

                    //obtain stop name by hardcoding so sorry
                    if strings.Contains(text, "--") == false && strings.Contains(text, ":") == false && strings.Contains(text, "X") == false {
                        stop_name = cell.String()
                        stop_name = text //declare stop name
                    }

                    //obtain time
                    if strings.Contains(text, ":")  {
                        time = cell.String()
                        time = text //declare time
                    }

                    //check again for no cells
                    if route != "" && day != "" && stop_name != "" && time != ""{
                        stops = append(stops, route, day, stop_name, time) //put into json
                    }
                }
            }
        }    
    }

    // export as json
    stopsJson, err := json.Marshal(stops)
    if err != nil {
        fmt.Printf("error:", err) //changed from println to Printf
    }
    err = ioutil.WriteFile("schedule.json", stopsJson, 0644)
}