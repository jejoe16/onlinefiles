package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (model Model) StoreUUID(token uuid.UUID, accountID int, ip string, expire time.Time) map[string]string {
	resp := map[string]string{}
	_, err := model.Pool.Exec(context.TODO(), "INSERT INTO session(fk_account_id, token, expire, ip) VALUES($1, $2, $3, $4) ON CONFLICT (fk_account_id) DO UPDATE set token = $2, expire = $3, ip = $4;", accountID, token, expire, ip)

	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	return nil
}

func (model Model) GetUUID(token uuid.UUID) (int, time.Time, string, map[string]string) {
	resp := map[string]string{}

	var ip string
	var expire time.Time
	var userID int

	err := model.Pool.QueryRow(context.TODO(), "SELECT fk_account_id, expire, ip FROM session WHERE token = $1", token).Scan(userID, expire, ip)

	if err != nil {
		resp["system"] = err.Error()
		return userID, expire, ip, resp
	}

	return userID, expire, ip, nil
}
