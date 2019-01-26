package sessions

import (
	"context"
	"fmt"
	"net/http"
)

// sessionInfo stores a session tracked by the registry.
type sessionInfo struct {
	s *Session
	e error
}

// contextKey is the type used to store the registry in the context.
type contextKey int

// registryKey is the key used to store the registry in the context.
const registryKey contextKey = 0

// GetRegistry returns a registry instance for the current request.
func GetRegistry(r *http.Request) *Registry {
	var ctx = r.Context()
	registry := ctx.Value(registryKey)
	if registry != nil {
		return registry.(*Registry)
	}
	newRegistry := &Registry{
		request:  r,
		sessions: make(map[string]sessionInfo),
	}
	*r = *r.WithContext(context.WithValue(ctx, registryKey, newRegistry))
	return newRegistry
}

// Registry stores sessions used during a request.
type Registry struct {
	request  *http.Request
	sessions map[string]sessionInfo
}

// Get registers and returns a session for the given name and session store.
//
// It returns a new session if there are no sessions registered for the name.
func (s *Registry) Get(store Store, name string) (session *Session, err error) {
	if !isCookieNameValid(name) {
		return nil, fmt.Errorf("sessions: invalid character in cookie name: %s", name)
	}
	if info, ok := s.sessions[name]; ok {
		session, err = info.s, info.e
	} else {
		session, err = store.New(s.request, name)
		session.name = name
		s.sessions[name] = sessionInfo{s: session, e: err}
	}
	session.store = store
	return
}

// Save saves all sessions registered for the current request.
func (s *Registry) Save(w http.ResponseWriter) error {
	var errMulti MultiError
	for name, info := range s.sessions {
		session := info.s
		if session.store == nil {
			errMulti = append(errMulti, fmt.Errorf(
				"sessions: missing store for session %q", name))
		} else if err := session.store.Save(s.request, w, session); err != nil {
			errMulti = append(errMulti, fmt.Errorf(
				"sessions: error saving session %q -- %v", name, err))
		}
	}
	if errMulti != nil {
		return errMulti
	}
	return nil
}

// Error ----------------------------------------------------------------------

// MultiError stores multiple errors.
//
// Borrowed from the App Engine SDK.
type MultiError []error

func (m MultiError) Error() string {
	s, n := "", 0
	for _, e := range m {
		if e != nil {
			if n == 0 {
				s = e.Error()
			}
			n++
		}
	}
	switch n {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", s, n-1)
}
