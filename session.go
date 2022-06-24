package main

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserID  int
	Expires time.Time
	IP      string
}

type contextKey int

const (
	contextKeyRequestID contextKey = iota
)

func (a *App) NewToken(userID int, ip string) (uuid.UUID, map[string]string) {
	token := uuid.New()
	expire := time.Now().Add(time.Hour * 24 * 30)
	resp := a.Model.StoreUUID(token, userID, ip, expire)
	a.Sessions[token] = Session{
		UserID:  userID,
		Expires: expire,
		IP:      ip,
	}
	return token, resp
}

func (a *App) Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := uuid.Parse(r.Header.Get("Auth"))
		if err != nil {
			respondWithJSON(w, http.StatusForbidden, "Invalid token")
			return
		}
		session, ok := a.Sessions[token]

		if !ok {
			userID, expire, ip, err := a.Model.GetUUID(token)
			if err != nil || ip == "" || expire.IsZero() {
				respondWithJSON(w, http.StatusForbidden, "In order to continue please login")
				return
			}

			session = Session{
				UserID:  userID,
				Expires: expire,
				IP:      ip,
			}
			a.Sessions[token] = session

		}

		if session.Expires.Before(time.Now()) {
			respondWithJSON(w, http.StatusForbidden, "Session expired")
			return
		}

		// if GetIP(r) != session.IP {
		// 	respondWithJSON(w, http.StatusForbidden, "Wrong user")
		// 	return
		// }
		ctx := context.WithValue(r.Context(), contextKeyRequestID, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// func (a *App) NewToken(userID int, ip string) (uuid.UUID, map[string]string) {
// 	token := uuid.NewV4()
// 	expire := time.Now().Add(time.Hour * 24 * 30) // Default token life set to 30 days
// 	resp := a.Model.StoreUUID(token, userID, ip, expire)
// 	a.Sessions[token] = Session{
// 		UserID:  userID,
// 		Expires: expire,
// 		IP:      ip,
// 	}

// func (a *App) CreatStoredSession(rw http.ResponseWriter, accountId int) map[string]string {
// 	token := uuid.New()

// 	exp := time.Now().Add(30 * (24 * time.Hour)) // Stored session length default to 30 days

// 	http.SetCookie(rw, &http.Cookie{
// 		Name:    "remember-token",
// 		Value:   token.String(),
// 		Secure:  true,
// 		Expires: exp,
// 	})

// 	err := a.Model.StoreSession(accountId, token.String(), exp)
// 	return err
// }

// func (a *App) CreateSession(w *http.Request, r http.ResponseWriter, accountID int) {
// 	session, err := a.Sessions.Get(w, "session-token")
// 	if err != nil {
// 		respondWithJSON(r, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	session.Values["id"] = accountID
// 	err = session.Save(w, r)
// 	if err != nil {
// 		respondWithJSON(r, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// }

// func GetSessionID(s *sessions.Session) (int64, bool) {
// 	val := s.Values["id"]
// 	id, ok := val.(int64)
// 	if !ok {
// 		return id, false
// 	}
// 	return id, true
// }

// func (a App) IsAuthenticated(w http.ResponseWriter, r *http.Request) (bool, int64, map[string]string) {
// 	session, err := a.Sessions.Get(r, "session-token")
// 	if err == nil {
// 		id, authenitcated := GetSessionID(session)
// 		if authenitcated {
// 			return true, id, nil
// 		}
// 	}

// 	cookie, err := r.Cookie("remember-token")
// 	if err == nil {
// 		token := cookie.Value
// 		storedSession, err := a.Model.GetStoredSession(token)
// 		if err != nil {
// 			return false, 0, err
// 		}

// 		if storedSession.Expire.After(time.Now()) {
// 			a.CreateSession(r, w, storedSession.AccountID)
// 			return true, storedSession.ID, nil
// 		}
// 	}
// 	return false, 0, nil
// }

// func (a *App) Verify(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authed, id, err := a.IsAuthenticated(w, r)
// 		if err != nil {
// 			http.Redirect(w, r, "signin", http.StatusTemporaryRedirect)
// 			return
// 		}
// 		if !authed {
// 			http.Redirect(w, r, "signin", http.StatusTemporaryRedirect)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), contextKeyRequestID, id)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
