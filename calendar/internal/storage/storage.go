package storage

import (
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"

	"calendar/internal/models"
)

type JWTGenerator interface {
	GenerateJWT(username string) (signedToken string, err error)
}

const timePattern = "2006-01-02 15:04:05.000000000"

type PostgresDB struct {
	db   *sql.DB
	log  *log.Logger
	auth JWTGenerator
}

func NewPostgresQL(l *log.Logger, auth JWTGenerator) *PostgresDB {
	return &PostgresDB{
		log:  l,
		auth: auth,
	}
}

func (s *PostgresDB) Connect(dbURL string) (err error) {
	s.db, err = sql.Open("postgres", dbURL)
	if err != nil {
		return
	}

	return s.db.Ping()
}

func (s *PostgresDB) Close() {
	s.db.Close()
}

// just generate new token on each login
func (s *PostgresDB) UserToken(usr, psw string) (token string, err error) {
	err = s.addUser(usr, psw)
	if err != nil {
		s.log.Println(err)
		return
	}

	return s.generateToken(usr, psw)
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
		return
	}

	//todo: postgres returns incorrect time pattern: 2021-10-22T03:36:56.363372Z -- result contains T and Z. TODO: clarify
	eventRow := s.db.QueryRow("SELECT * FROM events WHERE id=$1 AND user_id=$2", eventID, user.Login)
	var usrID string
	err = eventRow.Scan(
		&event.ID,
		&usrID,
		&event.Title,
		&event.Description,
		&event.Time,
		&event.Duration,
		pq.Array(&event.Notes),
	)

	event.Timezone = user.Timezone

	if err != nil {
		s.log.Println(err)
	}

	return
}

func (s *PostgresDB) UpdateEvent(login, eventID string, event models.Event) (err error) {
	sqlStatement := `UPDATE events
	SET title=$3, description=$4, datetime=$5, duration=$6, notes=$7
	WHERE id=$1 AND user_id=$2`

	_, err = s.db.Exec(
		sqlStatement,
		eventID,
		login,
		event.Title,
		event.Description,
		event.Time,
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
		t, err := time.Parse(timePattern, events[i].Time)
		if err != nil {
			s.log.Println(err)
			return err
		}

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
		return
	}

	eventRows, err := s.db.Query("SELECT id FROM events WHERE user_id=$1", user.Login)
	if err != nil {
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

func (s *PostgresDB) generateToken(usr, psw string) (token string, err error) {
	_ = psw

	token, err = s.auth.GenerateJWT(usr)
	if err != nil {
		s.log.Println(err)
		return
	}

	return
}

func (s *PostgresDB) user(login string) (user models.User, err error) {
	row := s.db.QueryRow("SELECT login, timezone FROM users WHERE login=$1", login)
	err = row.Scan(&user.Login, &user.Timezone)

	if err != nil {
		s.log.Println(err)
	}

	return
}
