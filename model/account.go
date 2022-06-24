package model

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/wes-kay/golang_asset_engine/core"
	"golang.org/x/crypto/bcrypt"
)

// swagger:model Account
type Account struct {
	ID               int        `json:"id" validate:"numeric"`
	Created          *time.Time `json:"created"`
	Email            string     `json:"email" validate:"required,alphanum"`
	Status           *int       `json:"status" validate:"gte=0,numeric"`
	Role             *int       `json:"role" validate:"gte=0,numeric"`
	Temp             bool       `json:"temp"`
	EmailValidate    bool       `json:"email_validated" validate:"boolean"`
	ProfileCompleted bool       `json:"profile_completed"`
	ProfileValidated bool       `json:"profile_validated"`
	ProfileImage     *string    `json:"profile_image"`
	Phone            *string    `json:"phone"`
	Name             *string    `json:"name" validate:"lte=100"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	Country          *string    `json:"country"`
	Passport         *string    `json:"passport"`
	Address          *string    `json:"address"`
	Kin              *string    `json:"kin"`
	Rank             *string    `json:"rank"`
	Licence          *string    `json:"licence"`
	ShipExperience   *int       `json:"ship_experience"`
	Experience       *int       `json:"experience"`
	Nationality      *string    `json:"nationality"`
}

// swagger:model AccountAdmin
type AccountAdmin struct {
	ID     int `json:"id" validate:"required,numeric"`
	Status int `json:"status" validate:"gte=0,required,numeric"`
	Role   int `json:"role" validate:"gte=0,required,numeric"`
}

type Accounts []Account

func (model *Model) CreateAccount(email string, password string, core *core.Core) map[string]string {
	resp := map[string]string{}
	tx, err := model.Pool.Begin(context.TODO())
	defer tx.Rollback(context.TODO())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	var accountID int
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = tx.QueryRow(context.TODO(), "INSERT INTO account(email, password, temp) VALUES($1, $2, $3) RETURNING id", strings.ToLower(email), string(bytes), false).Scan(&accountID)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}
	uid := uuid.New()

	core.ActivationEmail(email, uid.String())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "INSERT INTO email_verification(fk_account_id, code) VALUES($1, $2)", accountID, uid.String())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "INSERT INTO account_profile(fk_account_id) VALUES($1)", accountID)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	return nil
}

func (model Model) CreateTempAccount(email string, core *core.Core) map[string]string {
	resp := map[string]string{}
	tx, err := model.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	var accountID int

	err = tx.QueryRow(context.Background(), "INSERT INTO account(email, password, temp) VALUES($1, $2, $3) RETURNING id", email, nil, true).Scan(&accountID)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}
	uid := uuid.New()

	core.ActivationEmail(email, uid.String())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO email_verification(fk_account_id, code) VALUES($1, $2)", accountID, uid.String())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "INSERT INTO account_profile(fk_account_id, profile_completed) VALUES($1, false)", accountID)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = tx.Commit(context.Background())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	return nil
}

func (model Model) AccountActivationResend(core *core.Core, email string) map[string]string {
	resp := map[string]string{}

	tx, err := model.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	var id int
	tx.QueryRow(context.TODO(), "SELECT id FROM account WHERE email = $1", email).Scan(&id)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	var code uuid.UUID
	tx.QueryRow(context.TODO(), "SELECT code FROM email_verification WHERE fk_account_id = $1", id).Scan(code)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	core.ActivationEmail(email, code.String())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	return resp
}

func (model Model) UpdateAccountByID(sess *session.Session, core *core.Core, accountID int, fileName string, name string, address *string, country *string, phone *string, passportNumber string, dateOfBirth *time.Time, kin *string, rank *string, licence *string, nationality string) map[string]string {
	resp := map[string]string{}

	tx, err := model.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	// fileName := "images/" + uuid.New().String() + filepath.Ext(profileImage.FileHeader.Filename)

	_, err = tx.Exec(context.TODO(), "UPDATE account SET profile_image = $2, name = $3, address = $4, country = $5, phone = $6, passport_number = $7, date_of_birth = $8, kin = $9, rank = $10, licence = $11, nationality = $12, profile_completed = false, profile_validated = true WHERE id = $1", accountID, fileName, name, address, country, phone, passportNumber, dateOfBirth, kin, rank, licence, nationality)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "UPDATE account_profile SET profile_image = null, passport_image_one = null, passport_image_two = null, name = null, passport_number = null, date_of_birth = null, phone = null, country = null, address = null, kin = null, rank = null, licence = null, ship_experience = null, experience = null, nationality = null, completed = false WHERE id = $1", accountID)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	//TODO: Delete aws
	//TODO: set completed
	var accountEmail string
	err = tx.QueryRow(context.TODO(), "SELECT email FROM account WHERE id = $1", accountID).Scan(&accountEmail)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = core.ApprovedEmail(accountEmail)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	return resp
}

func (model Model) UpdateAdminAccount(id int, status int, role int) map[string]string {
	resp := map[string]string{}
	res, err := model.Pool.Exec(context.TODO(), "UPDATE account SET status = $2, role = $3 WHERE id = $1", id, status, role)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	if res.RowsAffected() == 0 {
		resp["system"] = "no rows updated"
		return resp
	}
	return resp
}

func (model Model) UpdateAccountPassword(email string, password string) map[string]string {
	resp := map[string]string{}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	res, err := model.Pool.Exec(context.TODO(), "UPDATE account SET password = $2, temp = $3 WHERE email = $1", email, string(bytes), false)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	if res.RowsAffected() == 0 {
		resp["system"] = "no rows updated"
		return resp
	}
	return resp
}

func (model Model) GetAccountById(id int) (Account, map[string]string) {
	resp := map[string]string{}
	var account Account
	err := model.Pool.QueryRow(context.TODO(), "SELECT id, created, email, status, role, temp, email_validated, profile_completed, profile_validated, profile_image, name, address, country, phone, passport_number, date_of_birth, kin, rank, licence, ship_experience, experience, nationality from account WHERE id = $1", id).Scan(
		&account.ID, &account.Created, &account.Email, &account.Status, &account.Role, &account.Temp, &account.EmailValidate, &account.ProfileCompleted, &account.ProfileValidated, &account.ProfileImage, &account.Name, &account.Address, &account.Country, &account.Phone, &account.Passport, &account.DateOfBirth, &account.Kin, &account.Rank, &account.Licence, &account.ShipExperience, &account.Experience, &account.Nationality)
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return Account{}, resp
	}
	return account, resp
}

func (model Model) GetAccountByEmail(email string) (Account, map[string]string) {
	resp := map[string]string{}
	var account Account
	err := model.Pool.QueryRow(context.TODO(), "SELECT id, created, email, status, role, temp, email_validated, profile_validated, profile_image, name, address, country, phone, passport_number, date_of_birth, kin, rank, licence, ship_experience, experience, nationality FROM account WHERE email = $1", email).Scan(&account.ID, &account.Created, &account.Email, &account.Status, &account.Role, &account.Temp, &account.EmailValidate, &account.EmailValidate, &account.ProfileImage, &account.Name, &account.Address, &account.Country, &account.Phone, &account.Passport, &account.DateOfBirth, &account.Kin, &account.Rank, &account.Licence, &account.ShipExperience, &account.Experience, &account.Nationality)

	// err := model.Pool.QueryRow(context.TODO(), "SELECT id, name, email, status, role, temp, email_validated, image, phone, country, passport_number, address FROM account WHERE email = $1", email).Scan(&account.ID, &account.Name, &account.Email, &account.Status, &account.Role, &account.Temp, &account.EmailValidate, &account.ProfileImage, &account.Phone, &account.Country, &account.Passport, &account.Address)
	if err != nil {
		resp["system"] = err.Error()
		return Account{}, resp
	}
	return account, resp
}

func (model Model) DeleteAccount(id int) map[string]string {
	resp := map[string]string{}
	tx, err := model.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "DELETE FROM account_profile WHERE fk_account_id = $1", id)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	//TODO: Get all certificates
	rows, err := model.Pool.Query(context.TODO(), "SELECT id FROM certificate WHERE fk_account_id = $1", id)
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return resp
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			resp["system"] = err.Error()
			return resp
		}
		_, err = tx.Exec(context.TODO(), "DELETE FROM certificate_accessed WHERE fk_certificate_id = $1", id)
		if err != nil {
			resp["system"] = err.Error()
			return resp
		}
	}

	_, err = tx.Exec(context.TODO(), "DELETE FROM certificate WHERE fk_account_id = $1", id)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "DELETE FROM course WHERE fk_account_id = $1", id)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "DELETE FROM email_verification WHERE fk_account_id = $1", id)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "DELETE FROM session WHERE fk_account_id = $1", id)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "DELETE FROM alert WHERE fk_account_id = $1", id)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "DELETE FROM terms WHERE fk_account_id = $1", id)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "DELETE FROM account WHERE id = $1", id)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	return resp
}

func (model Model) GetAllAccounts() (Accounts, map[string]string) {
	resp := map[string]string{}

	rows, err := model.Pool.Query(context.TODO(), "SELECT id, created, email, status, role, temp, email_validated, profile_completed, profile_validated, profile_image, name, address, country, phone, passport_number, date_of_birth, kin, rank, licence, ship_experience, experience, nationality FROM account ORDER BY id")
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return nil, resp
	}
	defer rows.Close()

	var data Accounts
	for rows.Next() {
		var account Account
		err = rows.Scan(&account.ID, &account.Created, &account.Email, &account.Status, &account.Role, &account.Temp, &account.EmailValidate, &account.ProfileCompleted, &account.ProfileValidated, &account.ProfileImage, &account.Name, &account.Address, &account.Country, &account.Phone, &account.Passport, &account.DateOfBirth, &account.Kin, &account.Rank, &account.Licence, &account.ShipExperience, &account.Experience, &account.Nationality)
		if err != nil {
			continue
		}

		data = append(data, account)
	}

	if err = rows.Err(); err != nil {
		resp["system"] = err.Error()
		return nil, resp
	}

	return data, resp
}
