package model

import (
	"context"
	"net/http"
	"time"
)

type StoredSession struct {
	ID        int64
	Token     string
	Expire    time.Time
	AccountID int
}

type StoredSessions []StoredSession

func (model Model) StoreSession(accountID int, token string, expire time.Time) map[string]string {
	resp := map[string]string{}
	_, err := model.Pool.Exec(context.TODO(), "INSERT INTO session(id, token, expire) VALUES($1, $2, $3)", accountID, token, expire)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}
	return resp
}

func (model Model) GetStoredSession(token string) (StoredSession, map[string]string) {
	resp := map[string]string{}
	var session StoredSession

	err := model.Pool.QueryRow(context.TODO(), "SELECT id, expire FROM session WHERE token = $1", token).Scan(session.AccountID, session)

	if err != nil {
		resp["system"] = err.Error()
		return session, resp
	}

	return session, nil
}

func GetIP(r *http.Request) string {
	return r.RemoteAddr
}
