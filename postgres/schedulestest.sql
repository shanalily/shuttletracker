insert into stops (id, latitude, longitude, name, description) values
	(1, 42.73029109316892, -73.67655873298646, 'Student Union', 'Shuttle stop in front of the Student Union'),
	(19, 42.725010644975015, -73.67067933082582, 'Tibbits Ave.', 'Stop at Tibbits Ave.'),
	(14, 42.73637487312414, -73.67058759924475, 'Colonie Apartments', 'Stop at Colonie Apartments');

insert into schedules (id, name, weekend, west) values (1, 'Monday - Thursday East Route', 'f', 'f');

insert into schedule_times (schedule_id, stop_id, time) values
	(1, 1, 420), (1, 19, 423), (1, 14, 429), (1, 1, 440), (1, 1, 443);

SELECT
	st.id,
	r.name,
	s.name,
	st.time
FROM
	schedules r,
	stops s,
	schedule_times st
WHERE
	s.id = st.stop_id
	AND r.id = st.schedule_id
	AND r.id = 1;

SELECT
	st.id,
	r.name,
	s.name,
	st.time
FROM
	schedules r,
	stops s,
	schedule_times st
WHERE
	st.schedule_id = r.id
	and st.stop_id = s.id
	and st.schedule_id = 1
	and st.stop_id = 1
	and st.time >= 440
ORDER BY
	st.time ASC
LIMIT 1;