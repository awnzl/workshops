package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"sync"
	"time"
)

// SessionManager keeps track of all sessions from creation, updating
// to destroying.
type SessionManager struct {
	sessions map[string]Session
	mu       sync.Mutex
}

// Session stores the session's data
type Session struct {
	Data map[string]interface{}
	time time.Time
}

func (s *Session) Expired() bool {
	return time.Since(s.time) > 5*time.Second
}

// NewSessionManager creates a new sessionManager
func NewSessionManager() *SessionManager {
	m := &SessionManager{
		sessions: make(map[string]Session),
	}

	// TODO run cleaner job here in a goroutine
	go m.Cleaner()

	return m
}

// TODO implement cleaner job function
// you may use tickers from time package to execute job every second\2nd second etc.
func (m *SessionManager) Cleaner() {
	tk := time.NewTicker(time.Second)
	defer tk.Stop()

	for {
		<-tk.C

		m.mu.Lock()
		for k, s := range m.sessions {
			if s.Expired() {
				delete(m.sessions, k)
			}
		}
		m.mu.Unlock()
	}
}

// CreateSession creates a new session and returns the sessionID
func (m *SessionManager) CreateSession() (string, error) {
	sessionID, err := MakeSessionID()
	if err != nil {
		return "", err
	}

	// TODO use mutex
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[sessionID] = Session{
		Data: make(map[string]interface{}),
		time: time.Now(),
	}

	return sessionID, nil
}

// ErrSessionNotFound returned when sessionID not listed in
// SessionManager
var ErrSessionNotFound = errors.New("SessionID does not exists")

// GetSessionData returns data related to session if sessionID is
// found, errors otherwise
func (m *SessionManager) GetSessionData(sessionID string) (map[string]interface{}, error) {
	// TODO use mutex
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session.Data, nil
}

// UpdateSessionData overwrites the old session data with the new one
func (m *SessionManager) UpdateSessionData(sessionID string, data map[string]interface{}) error {
	// TODO use mutex
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	// Hint: you should renew expiry of the session here
	m.sessions[sessionID] = Session{
		Data: data,
		time: time.Now(),
	}

	return nil
}

func main() {
	// Create new sessionManager and new session
	m := NewSessionManager()
	sID, err := m.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created new session with ID", sID)

	// Update session data
	data := make(map[string]interface{})
	data["website"] = "longhoang.de"

	if err = m.UpdateSessionData(sID, data); err != nil {
		log.Fatal(err)
	}

	log.Println("Update session data, set website to longhoang.de")

	// Retrieve data from manager again
	updatedData, err := m.GetSessionData(sID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get session data:", updatedData)
}

// MakeSessionID is used to generate a random dummy sessionID
func MakeSessionID() (string, error) {
	buf := make([]byte, 26)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}
