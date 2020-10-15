package session

import (
	"crypto/md5"
	"encoding/hex"
)

type SessionData struct {
	login string
}

type Session struct {
	data map[string]*SessionData
}

func NewSession() *Session {
	session := new(Session)
	session.data = make(map[string]*SessionData)
	return session
}

func (s *Session) InitSession(login string) string {
	sessionId := GetMD5Hash(login)
	data := &SessionData{login: login}
	s.data[sessionId] = data
	return sessionId
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
