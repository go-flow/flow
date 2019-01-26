package flow

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Session wraps sesions in cleaner and easy to use sessions
type Session struct {
	Session *sessions.Session
	req     *http.Request
	res     http.ResponseWriter
}

// Save the current session.
func (s *Session) Save() error {
	return s.Session.Save(s.req, s.res)
}

// Get a value from the current session.
func (s *Session) Get(name interface{}) interface{} {
	return s.Session.Values[name]
}

// GetOnce gets a value from the current session and then deletes it.
func (s *Session) GetOnce(name interface{}) interface{} {
	if x, ok := s.Session.Values[name]; ok {
		s.Delete(name)
		return x
	}
	return nil
}

// Set a value onto the current session. If a value with that name
// already exists it will be overridden with the new value.
func (s *Session) Set(name, value interface{}) {
	s.Session.Values[name] = value
}

// Delete a value from the current session.
func (s *Session) Delete(name interface{}) {
	delete(s.Session.Values, name)
}

// Clear the current session
func (s *Session) Clear() {
	for k := range s.Session.Values {
		s.Delete(k)
	}
}

// Values returns all session values
func (s *Session) Values() map[interface{}]interface{} {
	return s.Session.Values
}
