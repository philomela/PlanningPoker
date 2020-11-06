package sessions

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

type SessionsTool struct {
	sessionStore   *sessions.CookieStore
	currentSession *sessions.Session
	currError      error
}

func InitSessionsTool() *SessionsTool {
	var s = SessionsTool{}
	s.sessionStore = sessions.NewCookieStore([]byte("philomelka"))
	s.sessionStore.Options.MaxAge = 3600
	return &s
}

func (s *SessionsTool) CreateNewSession(loginUser string, r *http.Request, w *http.ResponseWriter) {
	s.currentSession, _ = s.sessionStore.Get(r, "session")
	s.currentSession.Values["UserLogin"] = loginUser
	s.currentSession.Save(r, *w)
}

func (s *SessionsTool) CheckAndUpdateSession(r *http.Request, w *http.ResponseWriter) bool {
	s.currentSession, s.currError = s.sessionStore.Get(r, "session")

	fmt.Println(s.currentSession)
	untyped, ok := s.currentSession.Values["UserLogin"]
	if !ok {
		return false
	}
	userLogin, ok := untyped.(string)
	if !ok {
		return false
	}
	fmt.Println(userLogin)

	s.currentSession.Save(r, *w)
	return true
}

func (s *SessionsTool) GetUserLoginSession(r *http.Request) string {
	s.currentSession, s.currError = s.sessionStore.Get(r, "session")

	fmt.Println(s.currentSession)
	userLogin, ok := s.currentSession.Values["UserLogin"]
	if !ok {
		return ""
	}
	return userLogin.(string)
}
