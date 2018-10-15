package postgres

import (
	"database/sql"

	"github.com/wtg/shuttletracker"
)

// ScheduleService implements shuttletracker.ScheduleService.
type ScheduleService struct {
	db *sql.DB
}

// I need to be able to have the order multiple times...
// should not have stop only occur once for route, it should be time as well
// maybe use time instead of "order" to order the stops
// should I number each individual route within 1 schedule? for instance, starts at Union and stops at Union
func (ss *ScheduleService) initializeSchema(db *sql.DB) error {
	ss.db = db
	schema := `
CREATE TABLE IF NOT EXISTS schedules (
    id serial PRIMARY KEY,
	name text NOT NULL,
	created timestamp with time zone NOT NULL DEFAULT now(),
	updated timestamp with time zone NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS schedules_stops (
	id serial PRIMARY KEY,
	schedule_id integer REFERENCES schedules ON DELETE CASCADE NOT NULL,
	name text NOT NULL,
	arrival_time time NOT NULL,
	description text,
	created timestamp with time zone NOT NULL DEFAULT now(),
	updated timestamp with time zone NOT NULL DEFAULT now(),
	"order" integer NOT NULL,
	UNIQUE (schedule_id, "order")
);`
	_, err := ss.db.Exec(schema)
	return err
}

// might get rid of left join because I need ss.arrival time
// Schedules returns all Schedules in the database
func (ss *ScheduleService) Schedules() ([]*shuttletracker.Schedule, error) {
	schedules := []*shuttletracker.Schedule{}
	query := "SELECT s.id, s.name, s.created, s.updated, ss.arrival_time," +
		" array_remove(array_agg(ss.stop_id ORDER BY ss.order ASC), NULL) as stop_ids" +
		" FROM schedules s LEFT JOIN schedules_stops ss" +
		" ON s.id = ss.schedule_id GROUP BY s.id;"
	rows, err := ss.db.Query(query)
	if err != nil {
		return nil, err
	}
	// I need to grab the times
	for rows.Next() {
		s := &shuttletracker.Schedule{}
		err := rows.Scan(&s.ID, &s.Name, &s.Created, &s.Updated)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

// it's a left join but I want to be able to get ss.arrival_time...
// Schedule returns the Schedule with the provided ID (might change to name).
func (ss *ScheduleService) Schedule(id int64) (*shuttletracker.Schedule, error) {
	query := "SELECT s.name, s.created, s.updated, ss.arrival_time," +
		" array_remove(array_agg(rs.stop_id ORDER BY rs.order ASC), NULL) as stop_ids" +
		" FROM schedules s LEFT JOIN schedules_stops ss" +
		" ON s.id = ss.schedule_id WHERE s.id = $1 GROUP BY s.id;"
	row := ss.db.QueryRow(query, id)
	s := &shuttletracker.Schedule{
		ID: id,
	}
	err := row.Scan(&s.Name, &s.Created, &s.Updated)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// CreateSchedule creates a Schedule.
func (ss *ScheduleService) CreateSchedule(schedule *shuttletracker.Schedule) error {
	tx, err := ss.db.Begin()
	if err != nil {
		return err
	}
	// We can't really do anything if rolling back a transaction fails.
	// nolint: errcheck
	defer tx.Rollback()

	// insert schedule
	statement := "INSERT INTO schedules (name)" +
		" VALUES ($1) RETURNING id, created, updated;"
	row := tx.QueryRow(statement, schedule.Name)
	err = row.Scan(&schedule.ID, &schedule.Created, &schedule.Updated)
	if err != nil {
		return err
	}

	// to do
	// insert stop ordering
	// statement = "INSERT INTO schedules_stops (schedule_id, stop_id, \"order\")" +
	// 	" SELECT $1, stop_id, \"order\" - 1 AS \"order\" FROM" +
	// 	" unnest($2::integer[]) WITH ORDINALITY AS s(stop_id, \"order\");"
	// _, err = tx.Exec(statement, schedule.ID, pq.Array(schedule.StopIDs))
	// if err != nil {
	// 	return err
	// }

	return tx.Commit()
}

// DeleteSchedule deletes a Schedule.
func (ss *ScheduleService) DeleteSchedule(id int64) error {
	statement := "DELETE FROM schedules WHERE id = $1;"
	result, err := ss.db.Exec(statement, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return shuttletracker.ErrScheduleNotFound
	}

	return nil
}

// ModifySchedule modifies an existing Schedule.
func (ss *ScheduleService) ModifySchedule(schedule *shuttletracker.Schedule) error {
	tx, err := ss.db.Begin()
	if err != nil {
		return err
	}
	// We can't really do anything if rolling back a transaction fails.
	// nolint: errcheck
	defer tx.Rollback()

	// update schedule
	statement := "UPDATE schedules SET name = $1, updated = now()" +
		" WHERE id = $2 RETURNING updated;"
	row := tx.QueryRow(statement, schedule.Name, schedule.ID)
	err = row.Scan(&schedule.Updated)
	if err != nil {
		return err
	}

	// remove existing stop ordering
	_, err = tx.Exec("DELETE FROM schedules_stops WHERE schedule_id = $1;", schedule.ID)
	if err != nil {
		return err
	}

	// to do
	// insert stop ordering
	// statement = "INSERT INTO schedules_stops (schedule_id, stop_id, \"order\")" +
	// 	" SELECT $1, stop_id, \"order\" - 1 AS \"order\" FROM" +
	// 	" unnest($2::integer[]) WITH ORDINALITY AS s(stop_id, \"order\");"
	// _, err = tx.Exec(statement, schedule.ID, pq.Array(schedule.StopIDs))
	// if err != nil {
	// 	return err
	// }

	return tx.Commit()
}

// create schedule stops?

func (ss *ScheduleService) CreateScheduleStop(stop *shuttletracker.ScheduleStop) error {
	statement := "INSERT INTO schedules_stops (name, schedule_id, arrival_time, description) VALUES" +
		" ($1, $2, $3, $4) RETURNING id, created, updated;"
	row := ss.db.QueryRow(statement, stop.Name, stop.Description, stop.ArrivalTime)
	return row.Scan(&stop.ID, &stop.Created, &stop.Updated)
}

func (ss *ScheduleService) ScheduleStops() ([]*shuttletracker.ScheduleStop, error) {
	stops := []*shuttletracker.ScheduleStop{}
	query := "SELECT s.id, s.schedule_id, s.name, s.arrival_time, s.created, s.updated" +
		" FROM schedules_stops s;"
	rows, err := ss.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &shuttletracker.ScheduleStop{}
		err := rows.Scan(&s.ID, &s.ScheduleID, &s.Name, &s.ArrivalTime, &s.Created, &s.Updated)
		if err != nil {
			return nil, err
		}
		stops = append(stops, s)
	}
	return stops, nil
}

func (ss *ScheduleService) DeleteScheduleStop(id int64) error {
	statement := "DELETE FROM schedules_stops WHERE id = $1;"
	result, err := ss.db.Exec(statement, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return shuttletracker.ErrScheduleStopNotFound
	}

	return nil
}