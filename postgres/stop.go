package postgres

import (
	"database/sql"
	"fmt"

	"github.com/wtg/shuttletracker"
)

// StopService is an implementation of shuttletracker.StopService.
type StopService struct {
	db *sql.DB
}

func (ss *StopService) initializeSchema(db *sql.DB) error {
	ss.db = db
	schema := `
CREATE TABLE IF NOT EXISTS stops (
	id serial PRIMARY KEY,
	name text,
	description text,
	latitude double precision NOT NULL,
	longitude double precision NOT NULL,
	created timestamp with time zone NOT NULL DEFAULT now(),
	updated timestamp with time zone NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS schedules (
	id serial PRIMARY KEY,
	name text NOT NULL,
	weekend boolean NOT NULL,
	west boolean NOT NULL
);
CREATE TABLE IF NOT EXISTS schedule_times (
	id serial PRIMARY KEY,
	schedule_id integer REFERENCES schedules NOT NULL,
	stop_id integer REFERENCES stops NOT NULL,
	time integer NOT NULL,
	UNIQUE (schedule_id, stop_id, time)
);`
	_, err := ss.db.Exec(schema)
	return err
}

// CreateStop creates a Stop.
func (ss *StopService) CreateStop(stop *shuttletracker.Stop) error {
	statement := "INSERT INTO stops (name, description, latitude, longitude) VALUES" +
		" ($1, $2, $3, $4) RETURNING id, created, updated;"
	row := ss.db.QueryRow(statement, stop.Name, stop.Description, stop.Latitude, stop.Longitude)
	return row.Scan(&stop.ID, &stop.Created, &stop.Updated)
}

// Stops returns all Stops.
func (ss *StopService) Stops() ([]*shuttletracker.Stop, error) {
	stops := []*shuttletracker.Stop{}
	query := "SELECT s.id, s.name, s.created, s.updated, s.description, s.latitude, s.longitude" +
		" FROM stops s;"
	rows, err := ss.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &shuttletracker.Stop{}
		err := rows.Scan(&s.ID, &s.Name, &s.Created, &s.Updated, &s.Description, &s.Latitude, &s.Longitude)
		if err != nil {
			return nil, err
		}
		stops = append(stops, s)
	}
	return stops, nil
}

// DeleteStop deletes a Stop.
func (ss *StopService) DeleteStop(id int64) error {
	statement := "DELETE FROM stops WHERE id = $1;"
	result, err := ss.db.Exec(statement, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return shuttletracker.ErrStopNotFound
	}

	return nil
}



// CreateStop creates a ScheduleStop.
func (ss *StopService) CreateScheduleStop(stop *shuttletracker.ScheduleStop) error {
	// assume ScheduleStop has stop name and schedule name set instead of schedule_id and stop_id?
	statement := "INSERT INTO schedule_times (schedule_id, stop_id, time) VALUES" +
		" ($1, $2, $3) RETURNING id;"
	row := ss.db.QueryRow(statement, stop.ScheduleID, stop.StopID, stop.Time)
	return row.Scan(&stop.ID)
}

// ScheduleStops returns all ScheduleStops associated with one schedule.
func (ss *StopService) ScheduleStops(schedule_id int64) ([]*shuttletracker.ScheduleStop, error) {
	// what about finding all schedule stop ids, then for each id find all stops associated so that
	// I can group them together by schedule in JSON output?
	stops := []*shuttletracker.ScheduleStop{}
	query := "SELECT st.id, r.name, s.name, st.time " +
		"FROM schedules r, stops s, schedule_times st " +
		"WHERE s.id = st.stop_id AND r.id = st.schedule_id and r.id = $1;"
	rows, err := ss.db.Query(query, schedule_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &shuttletracker.ScheduleStop{}

		err := rows.Scan(&s.ID, &s.ScheduleName, &s.StopName, &s.Time)
		// fmt.Println(convertTime(s.Time))
		if err != nil {
			return nil, err
		}
		stops = append(stops, s)
	}
	return stops, nil
}

// StopTimes returns all ScheduleStops associated with one stop.
// This should actually take a stop id, schedule?, time, and number of entries wanted
func (ss *StopService) StopTimes(stop_id int64, schedule_id int64, departure_time int64, entries int64) ([]*shuttletracker.ScheduleStop, error) {
	// specify current day/time and get next few stops? It depends on what the frontend people need
	stops := []*shuttletracker.ScheduleStop{}
	query := "SELECT st.id, r.name, s.name, st.time " +
		"FROM schedules r, stops s, schedule_times st " +
		"WHERE st.schedule_id = r.id and st.stop_id = s.id and st.stop_id = $1 and st.schedule_id = $2 and st.time >= $3" + 
		"ORDER BY st.time ASC " +
		"LIMIT $4;"
	rows, err := ss.db.Query(query, stop_id, schedule_id, departure_time, entries)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &shuttletracker.ScheduleStop{}

		err := rows.Scan(&s.ID, &s.ScheduleName, &s.StopName, &s.Time)
		fmt.Println(convertTime(s.Time))
		if err != nil {
			return nil, err
		}
		stops = append(stops, s)
	}
	return stops, nil
}

// DeleteStop deletes a ScheduleStop.
func (ss *StopService) DeleteScheduleStop(id int64) error {
	statement := "DELETE FROM schedule_times WHERE id = $1;"
	result, err := ss.db.Exec(statement, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return shuttletracker.ErrStopNotFound
	}

	return nil
}

// not needed here, should probably be somewhere else?
func convertTime(time int) string {
	hour := time / 60
	minute := time % 60
	return fmt.Sprintf("%d:%d", hour, minute)
}
