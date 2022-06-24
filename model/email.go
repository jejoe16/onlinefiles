package model

import (
	"context"
)

type EmailValidate struct {
	ID            int  `json:"id"`
	UserAccountId int  `json:"fk_account_id"`
	Used          bool `json:"used"`
}

func (model *Model) EmailVerification(code string) map[string]string {
	responseErrors := map[string]string{}
	var emailValidate EmailValidate

	tx, err := model.Pool.Begin(context.TODO())
	defer tx.Rollback(context.Background())
	if err != nil {
		responseErrors["system"] = err.Error()
		return responseErrors
	}

	err = tx.QueryRow(context.TODO(), "SELECT id, fk_account_id, used FROM email_verification WHERE code = $1", code).Scan(&emailValidate.ID, &emailValidate.UserAccountId, &emailValidate.Used)
	if err != nil {
		responseErrors["system"] = err.Error()
		return responseErrors
	}

	if emailValidate.Used {
		responseErrors["system"] = "email already verfied"
		return responseErrors
	}

	_, err = tx.Exec(context.TODO(), "UPDATE account SET role = 2, email_validated = true WHERE id = $1", emailValidate.UserAccountId)
	if err != nil {
		responseErrors["system"] = err.Error()
		return responseErrors
	}

	_, err = tx.Exec(context.TODO(), "UPDATE email_verification SET used = true WHERE id = $1", emailValidate.ID)
	if err != nil {
		responseErrors["system"] = err.Error()
		return responseErrors
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		responseErrors["system"] = err.Error()
		return responseErrors
	}

	return responseErrors
}
