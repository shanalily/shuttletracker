package schedule

import (
    "encoding/json"
    "fmt"
    "github.com/360EntSecGroup-Skylar/excelize"
    "strings"
    "io/ioutil"
    "net/http"
    "os"
    "io"
)

// structure to be converted into JSON format
type Stop struct{
     Location string `json: location`
     Times []string `json: times `
}
type Operation struct{
     Name string `json: name`
     Stops []Stop `json: stops`
}
type Route struct {
     Name string `json: route_name`
     Operations []Operation `json: operations`
}
type Schedule struct {
     Routes []Route `json: routes`
}

var loc string;
var direction string;
var operation string;
var time string;


// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}

// helper function to read stops 
func readStops(matrix [][]string, row int, cl int, sch *Schedule, rt int, op int,){
    // traverses the matrix until we reach a nonempty row
    col_length := 0
    for row < len(matrix){
        // non blank signifies a row of all stops
        if matrix[row][cl] != ""{
          // reads until we have a blank to determine how many stops
          for col:=cl; col < len(matrix[row]); col++{
            if matrix[row][col] == ""{
              col_length = col
              break
            }
          }
          break
        }
        row++
    }
    // read's a stop and then their times
    for c:= cl; c < col_length; c++ {
      loc = matrix[row][c]
      if strings.Compare(loc, "X") != 0 {

        var tmpCat Stop
        found := false
        loc_index := 0
        doublespace := 0

        // check if a stop exists
        for loc_index < len(sch.Routes[rt].Operations[op].Stops){
          if sch.Routes[rt].Operations[op].Stops[loc_index].Location == loc{
             found = true
             break
          }
          loc_index++
        }
        // if the stop is a new stop then it is added to the schedule
        if !found {
          tmpCat = Stop{Times: []string{}}
          tmpCat.Location = loc
          sch.Routes[rt].Operations[op].Stops = append(sch.Routes[rt].Operations[op].Stops, tmpCat)
        }
        // read times by iterating down the column
        for r := row + 1; r < len(matrix); r++{
          // guard to prevent reading out of bounds
          if(len(matrix[r]) == 0){
            break
          }
          time = matrix[r][c]
          if strings.Contains(time, ":") && strings.Compare(time, "X") != 0 {
            doublespace = 0
            time_index :=0
            found = false
            for time_index < len(sch.Routes[rt].Operations[op].Stops[loc_index].Times){
              if sch.Routes[rt].Operations[op].Stops[loc_index].Times[time_index] == time{
                 found = true
                 break
              }
              time_index++
            }
            if !found {
              sch.Routes[rt].Operations[op].Stops[loc_index].Times = append(sch.Routes[rt].Operations[op].Stops[loc_index].Times, time)
            }
          } else  { // has a tolerance of 2 non times, if that occurs then we break
            doublespace++
            if doublespace == 2{
              break
            }
          }
        }
      }
    }
}
// helper function to read headers
func readHeader(name string, sched *Schedule, xlsx *excelize.File){
        // Store a sheet in a 2x2 array 
      var matrix [][]string
      trim := false
      rows, err := xlsx.Rows(name)
      if err != nil {
          fmt.Printf("Error Excell File does not contain rows!:", err)
      }
      // dumps the excel file in a 2x2 string matrix removes lines in the beginning
      for rows.Next(){
         col := rows.Columns()
          if !trim {
            for _, colCell := range col {
              if colCell != ""{
                trim = true
                matrix = append(matrix, col)
                break
              }
            }
          } else {
            matrix = append(matrix, col)
          }
        // matrix = append(matrix, rows.Columns())
      }
      // loop through string matrix
      for i:= 0; i < len(matrix); i++{
          for j:=0; j < len(matrix[i]); j++ {
              text := matrix[i][j]
              if text != "" && (strings.Contains(text, "Schedule") || strings.Contains(text, "Shuttle") || strings.Contains(text, "WEEKEND")){
                  if strings.Contains(text, "West") {
                      direction = "West"
                  }
                  if strings.Contains(text, "East") {
                      direction = "East"
                  }
                  if strings.Contains(text, "WEEKEND"){
                      direction = "Weekend/Late Night"
                      operation = "Weekend/Late Night"
                  }
                  if strings.Contains(text, "Express"){
                      operation = "Express" + (strings.Split(text, " " + direction))[0]
                  }
                  if !strings.Contains(text, "Express") && !strings.Contains(text, "WEEKEND"){
                      operation = (strings.Split(text, " " + direction))[0]
                  }

                  rt_index := 0
                  op_index := 0
                  found := false
                  var tmpCat Route
                  var tmpCat2 Operation

                  // search if the current route is unique
                  for rt_index < len(sched.Routes) {
                    if sched.Routes[rt_index].Name == direction{
                       found = true
                       break
                    }
                    rt_index++
                  }
                  // adds route to the struct 
                  if !found {
                    tmpCat = Route{Operations: []Operation{}}
                    tmpCat.Name = direction
                    sched.Routes = append(sched.Routes, tmpCat)
                  }

                  // search if the current operation is unique
                  found = false
                  for op_index < len(sched.Routes[rt_index].Operations) {
                    if sched.Routes[rt_index].Operations[op_index].Name == operation{
                       found = true
                       break
                    }
                    op_index++
                  }
                  // adds the operation into the struct
                  if !found {
                    tmpCat2 = Operation{Stops: []Stop{}}
                    tmpCat2.Name = operation
                    sched.Routes[rt_index].Operations = append(sched.Routes[rt_index].Operations, tmpCat2)
                  }
                  readStops(matrix, i + 1, j, sched, rt_index, op_index)
              }
            }
        }
}

// driver function to read stops without a specified link
func ReadDefault() {
    fileUrl := "https://rpi.box.com/shared/static/naf8gm1wjdor8tbebho5k0t28wksaygd.xlsx"
    fileName := "master_schedule.xlsx"

    fmt.Println("Downloading file...")
      err := DownloadFile(fileName, fileUrl)
      if err != nil {
          panic(err)
      }
      fmt.Println("File Downloaded Successfully")
      fmt.Println("Reading File...")
      xlsx, err := excelize.OpenFile(fileName)
      if err != nil {
          fmt.Println(err)
      }
      // initialize schedule data structure 
      sched := &Schedule{Routes: []Route{}}

      // loop through sheets
      for _, sheet_name := range xlsx.GetSheetMap(){
          readHeader(sheet_name, sched, xlsx)
      }
      fmt.Println("File Read Success!")
      stopsJson, err := json.Marshal(sched)
      if err != nil {
          fmt.Printf("error:", err) //changed from println to Printf
      }
      err = os.Remove(fileName)
      err = os.Remove("schedule.json")
      err = ioutil.WriteFile("schedule.json", stopsJson, 0644)
      fmt.Println("Json Created")
  }

// reads an excel file from a specified link
func ReadLink(FileLink string) {
    fileUrl := link 
    fileName := "master_schedule.xlsx"

    fmt.Println("Downloading file...")
      err := DownloadFile(FileLink, fileUrl)
      if err != nil {
          panic(err)
      }
      fmt.Println("File Downloaded Successfully")
      fmt.Println("Reading File...")
      xlsx, err := excelize.OpenFile(fileName)
      if err != nil {
          fmt.Println(err)
      }
      // initialize schedule data structure 
      sched := &Schedule{Routes: []Route{}}

      // loop through sheets
      for _, sheet_name := range xlsx.GetSheetMap(){
          readHeader(sheet_name, sched, xlsx)
      }
      fmt.Println("File Read Success!")
      stopsJson, err := json.Marshal(sched)
      if err != nil {
          fmt.Printf("error:", err) //changed from println to Printf
      }
      err = os.Remove(fileName)
      err = os.Remove("schedule.json")
      err = ioutil.WriteFile("schedule.json", stopsJson, 0644)
      fmt.Println("Json Created")
}