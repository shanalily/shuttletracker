package postgres

import (
	"database/sql"
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
CREATE TABLE IF NOT EXISTS schedule_stops (
	id serial PRIMARY KEY,
	schedule_id integer REFERENCES schedules ON DELETE CASCADE NOT NULL,
	stop_id integer REFERENCES stops NOT NULL,
	arrival_time time,
	"order" integer NOT NULL,
	UNIQUE (schedule_id, "order")
);`
	_, err := ss.db.Exec(schema)
	return err
}

func (ss *ScheduleService) Schedules() ([]*shuttletracker.Schedule, error) {
	schedules := []*shuttletracker.Schedule{}
	query := "SELECT s.id, s.name, s.created, s.updated," +
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

//  Schedule returns the Schedule with the provided name.
func (ss *ScheduleService) Schedule(name string) (*shuttletracker.Schedule, error) {
	query := "SELECT s.name, s.created, s.updated," +
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
func (rs *ScheduleService) CreateSchedule(schedule *shuttletracker.Schedule) error {
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