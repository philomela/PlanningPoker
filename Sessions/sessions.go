package Sessions

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

type SessionsTool struct {
	sessionStore *sessions.CookieStore
}

func InitSessionsTool() *SessionsTool {
	var s = SessionsTool{}
	s.sessionStore = sessions.NewCookieStore([]byte("philomelka"))
	s.sessionStore.Options.MaxAge = 3600
	return &s
}

func (s *SessionsTool) CreateNewSession(loginUser string, r *http.Request, w *http.ResponseWriter) {
	currentSession, _ := s.sessionStore.Get(r, "session")
	currentSession.Values["UserLogin"] = loginUser
	currentSession.Save(r, *w)
}

func (s *SessionsTool) CheckAndUpdateSession(r *http.Request, w *http.ResponseWriter) bool {
	currentSession, err := s.sessionStore.Get(r, "session")
	if err != nil {
		fmt.Println(err)
	}

	untyped, ok := currentSession.Values["UserLogin"]
	if !ok {
		return false
	}
	_, ok = untyped.(string)
	if !ok {
		return false
	}

	currentSession.Save(r, *w)
	return true
}

func (s *SessionsTool) GetUserLoginSession(r *http.Request) string {
	currentSession, _ := s.sessionStore.Get(r, "session")

	userLogin, ok := currentSession.Values["UserLogin"]
	if !ok {
		return ""
	}
	return userLogin.(string)
}
