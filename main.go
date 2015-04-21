package main

import (
  "os"
  "fmt"
  "time"
  "strings"
  "net/http"
  "io/ioutil"
  "regexp"
  "encoding/json"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/mux"
)

/**
 *  Shuttle Tracker 
 *   Auto Updater - send request to iTrak API,
 *                  get updated shuttle info,
 *                  store updated records in db
 */

type VehicleUpdate struct {
  Id          bson.ObjectId                 `bson:"_id,omitempty"`
  VehicleId   string     `json:"vehicleId"   bson:"vehicleId,omitempty"`
  Lat         string     `json:"lat"         bson:"lat"`
  Lng         string     `json:"lng"         bson:"lng"`
  Heading     string     `json:"heading"     bson:"heading"`
  Speed       string     `json:"speed"       bson:"speed"`
  Lock        string     `json:"lock"        bson:"lock"`
  Time        string     `json:"time"        bson:"time"`
  Date        string     `json:"date"        bson:"date"`
  Status      string     `json:"status"      bson:"status"`
  Created     time.Time  `json:"created"     bson:"created"`
}

func (Shuttles *Shuttles) UpdateShuttles(dataFeed string, updateInterval int) {
  for {
    // Reference updates collection and close db session upon exit
    UpdatesCollection := Shuttles.Session.DB("shuttle_tracking").C("updates")

    // Make request to our tracking data feed
    resp, err := http.Get(dataFeed)
    if err != nil {
      continue;
    }
    defer resp.Body.Close()

    // Read response body content
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      continue;
    }

    delim := "eof"
    // Iterate through all vehicles returned by data feed
    vehicles_data := strings.Split(string(body), delim)
    for i := 1; i < len(vehicles_data)-1; i++ {

      // Match eatch API field with any number (+)
      //   of the previous expressions (\d digit, \. escaped period, - negative number)
      //   Specify named capturing groups to store each field from data feed
      re := regexp.MustCompile(`(?P<id>Vehicle ID:([\d\.]+)) (?P<lat>lat:([\d\.-]+)) (?P<lng>lon:([\d\.-]+)) (?P<heading>dir:([\d\.-]+)) (?P<speed>spd:([\d\.-]+)) (?P<lock>lck:([\d\.-]+)) (?P<time>time:([\d]+)) (?P<date>date:([\d]+)) (?P<status>trig:([\d]+))`)
      n := re.SubexpNames()
      match := re.FindAllStringSubmatch(vehicles_data[i], -1)[0]

      // Store named capturing group and matching expression as a key value pair
      result := map[string]string{}
      for i, item := range match {
        result[n[i]] = item
      }

      // Create new vehicle update & insert update into database
      update := VehicleUpdate { 
        Id:        bson.NewObjectId(),
        VehicleId: strings.Replace(result["id"], "Vehicle ID:", "", -1),
        Lat:       strings.Replace(result["lat"], "lat:", "", -1),
        Lng:       strings.Replace(result["lng"], "lon:", "", -1),
        Heading:   strings.Replace(result["heading"], "dir:", "", -1),
        Speed:     strings.Replace(result["speed"], "spd:", "", -1),
        Lock:      strings.Replace(result["lock"], "lck:", "", -1),
        Time:      strings.Replace(result["time"], "time:", "", -1),
        Date:      strings.Replace(result["date"], "date:", "", -1),
        Status:    strings.Replace(result["status"], "trig:", "", -1),
        Created:   time.Now()}

      err := UpdatesCollection.Insert(&update)

      if err != nil {
        fmt.Println(err)
      } else {
        fmt.Println(update)
      }
    }

    // Sleep for n seconds before updating again
    time.Sleep(time.Duration(updateInterval) * time.Second)
  }
}

/**
 *  Route handlers - API requests,
 *                   serve view files
 */

func IndexHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "index.html")
}

func VehiclesHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Vehicles")
}

func UpdatesHandler(w http.ResponseWriter, r *http.Request) {

}

/**
 *  Main - connect to database, 
 *         handle routing,
 *         start tracker go routine,
 *         serve requests
 */

type Configuration struct {
  DataFeed        string
  UpdateInterval  int
  MongoUrl        string
  MongoPort       string
}

type Shuttles struct {
  Session *mgo.Session
}

func ReadConfiguration(fileName string) Configuration { 
  // Open config file and decode JSON to Configuration struct
  file, _ := os.Open(fileName)
  decoder := json.NewDecoder(file)
  config := Configuration{}
  err := decoder.Decode(&config)
  if err != nil {
    fmt.Print("Unable to read config file: ")
    panic(err)
  }
  return config
}

func main() {
  // Read app configuration file 
  config := ReadConfiguration("conf.json")

  // Connect to MongoDB
  session, err := mgo.Dial(config.MongoUrl + ":" + config.MongoPort)
  if err != nil {
    panic(err)
  }
  // close Mongo session when server terminates
  defer session.Close()

  Shuttles := &Shuttles{session}

  // Start auto updater 
  go Shuttles.UpdateShuttles(config.DataFeed, config.UpdateInterval)

  // Routing 
  r := mux.NewRouter()
  r.HandleFunc("/", IndexHandler).Methods("GET")
  r.HandleFunc("/admin", IndexHandler).Methods("GET")
  r.HandleFunc("/vehicles", VehiclesHandler).Methods("GET")
  r.HandleFunc("/updates", UpdatesHandler).Methods("GET")
  // Static files
  r.PathPrefix("/bower_components/").Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components/"))))
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
  // Serve requests
  http.Handle("/", r)
  http.ListenAndServe(":8080", r)
}