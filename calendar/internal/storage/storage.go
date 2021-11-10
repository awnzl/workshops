package storage

import (
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"

	"calendar/internal/models"
)

const timePattern = "2006-01-02 15:04:05.000000000 -0700"

type PostgresDB struct {
	db   *sql.DB
	log  *log.Logger
}

func NewPostgresQL(l *log.Logger) *PostgresDB {
	return &PostgresDB{
		log:  l,
	}
}

func (s *PostgresDB) Connect(dbURL string) (err error) {
	s.db, err = sql.Open("postgres", dbURL)
	if err != nil {
		s.log.Println(err)
		return
	}

	return s.db.Ping()
}

func (s *PostgresDB) Close() {
	s.db.Close()
}

func (s *PostgresDB) LoginUser(usr, psw string) error {
	return s.addUser(usr, psw)
}

func (s *PostgresDB) User(user models.User) (err error) {
	sqlStatement := `UPDATE users SET timezone=$2 WHERE login=$1`

	_, err = s.db.Exec(
		sqlStatement,
		user.Login,
		user.Timezone,
	)

	if err != nil {
		s.log.Println(err)
	}

	return
}

func (s *PostgresDB) Event(userID, eventID string) (event models.Event, err error) {
	user, err := s.user(userID)
	if err != nil {
		s.log.Println(err)
		return
	}

	eventRow := s.db.QueryRow("SELECT * FROM events WHERE id=$1 AND user_id=$2", eventID, user.Login)
	err = eventRow.Scan(
		&event.ID,
		&userID,
		&event.Title,
		&event.Description,
		&event.Time,
		&event.Duration,
		pq.Array(&event.Notes),
	)

	t, err := time.Parse(time.RFC3339Nano, event.Time)
	if err != nil {
		s.log.Println(err, "arg:", event.Time)
		return
	}

	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		s.log.Println(err)
		return
	}
	t = t.In(loc)

	event.Time = t.String()
	event.Timezone = user.Timezone

	return
}

func (s *PostgresDB) UpdateEvent(login, eventID string, event models.Event) (err error) {
	sqlStatement := `UPDATE events
	SET title=$3, description=$4, datetime=$5, duration=$6, notes=$7
	WHERE id=$1 AND user_id=$2`

	t, err := s.userTime(login, event.Time)

	_, err = s.db.Exec(
		sqlStatement,
		eventID,
		login,
		event.Title,
		event.Description,
		t,
		event.Duration,
		pq.Array(event.Notes),
	)

	if err != nil {
		s.log.Println(err)
	}

	return
}

func (s *PostgresDB) AddEvents(login string, events []models.Event) error {
	sqlStatement := `INSERT INTO events (id, user_id, title, description, datetime, duration, notes)
	VALUES($1, $2, $3, $4, $5, $6, $7)
	`

	for i := range events {
		t, err := s.userTime(login, events[i].Time)

		_, err = s.db.Exec(
			sqlStatement,
			events[i].ID,
			login,
			events[i].Title,
			events[i].Description,
			t,
			events[i].Duration,
			pq.Array(events[i].Notes),
		)

		if err != nil {
			s.log.Println(err)
			return err
		}
	}

	return nil
}

func (s *PostgresDB) Events(usr string) (events []models.Event, err error) {
	user, err := s.user(usr)
	if err != nil {
		s.log.Println(err)
		return
	}

	eventRows, err := s.db.Query("SELECT id FROM events WHERE user_id=$1", user.Login)
	if err != nil {
		s.log.Println(err)
		return
	}

	defer func() {
		_ = eventRows.Close()
		err = eventRows.Err()
	}()

	for eventRows.Next() {
		id := ""
		err = eventRows.Scan(&id)
		if err != nil {
			return
		}

		event, err := s.Event(user.Login, id)
		if err != nil {
			return events, err
		}

		events = append(events, event)
	}

	return
}

func (s *PostgresDB) addUser(usr, psw string) error {
	_ = psw

	_, err := s.user(usr)
	if err == sql.ErrNoRows {
		_, err := s.db.Exec(
			"INSERT INTO users VALUES ($1, $2)",
			usr,
			"Local",
		)

		return err
	}

	return err
}

func (s *PostgresDB) user(login string) (user models.User, err error) {
	row := s.db.QueryRow("SELECT login, timezone FROM users WHERE login=$1", login)
	err = row.Scan(&user.Login, &user.Timezone)

	if err != nil {
		s.log.Println(err)
	}

	return
}

func (s *PostgresDB) userTime(login, tStr string) (tz time.Time, err error) {
	t, err := time.Parse(timePattern, tStr)
	if err != nil {
		s.log.Println(err)
		return
	}

	user, err := s.user(login)
	if err != nil {
		s.log.Println(err)
		return
	}

	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		s.log.Println(err)
		return
	}

	tz = t.In(loc)

	return
}
