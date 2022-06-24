package model

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/wes-kay/golang_asset_engine/core"
)

// swagger:model AccountProfile
type AccountProfile struct {
	ID               int        `json:"id"`
	Created          *time.Time `json:"created"`
	Updated          *time.Time `json:"updated"`
	ProfileImage     string     `json:"profile_image"`
	PassportImageOne string     `json:"passport_image_one"`
	PassportImageTwo string     `json:"passport_image_two"`
	Name             string     `json:"name"`
	PassportNumber   string     `json:"passport_number"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	Phone            *string    `json:"phone"`
	Country          *string    `json:"country"`
	Address          *string    `json:"address"`
	Kin              *string    `json:"kin"`
	Rank             *string    `json:"rank"`
	Licence          *string    `json:"licence"`
	// ShipExperience   *int       `json:"ship_experience"`
	// Experience       *int       `json:"experience"`
	Nationality      string `json:"nationality"`
	ProfileCompleted bool   `json:"profile_completed"`
	ProfileValidated bool   `json:"profile_validated"`
}

// swagger:model AccountProfileDTO
type AccountProfileDTO struct {
	ProfileImage     MultipartFile `json:"ProfileImage" validate:"required"`
	PassportImageOne MultipartFile `json:"PassportImageOne" validate:"required"`
	PassportImageTwo MultipartFile `json:"PassportImageTwo" validate:"required"`
	Name             string        `json:"name" validate:"required,lte=255"`
	Phone            *string       `json:"phone" validate:"lte=255"`
	PassportNumber   string        `json:"passport_number" validate:"lte=255"`
	DateOfBirth      *time.Time    `json:"date_of_birth" validate:"required"`
	Country          *string       `json:"country" validate:"lte=255"`
	Address          *string       `json:"address" validate:"lte=255"`
	Kin              *string       `json:"kin" validate:"lte=255"`
	Rank             *string       `json:"rank" validate:"lte=255"`
	Licence          *string       `json:"licence" validate:"lte=255"`
	// ShipExperience   *int          `json:"ship_experience" validate:"lte=255"`
	// Experience       *int          `json:"experience" validate:"lte=255"`
	Nationality string `json:"nationality" validate:"required,lte=255"`
}

// swagger:model AccountProfileVerifyDTO
type AccountProfileVerifyDTO struct {
	ID int `json:"id"`
}

// swagger:model AccountProfileRejectDTO
type AccountProfileRejectDTO struct {
	ID int `json:"id"`
}

func (model Model) UpdateAccountProfileByID(sess *session.Session, accountID int, profileImage MultipartFile, passportImageOne MultipartFile, passportImageTwo MultipartFile, profileName string, address *string, country *string, phone *string, passportNumber string, dateOfBirth *time.Time, kin *string, rank *string, licence *string, nationality string) map[string]string {
	resp := map[string]string{}

	tx, err := model.Pool.Begin(context.Background())
	defer tx.Rollback(context.TODO())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	profileImagePath := fmt.Sprintf("images/profile_images/%s%s", uuid.New().String(), filepath.Ext(profileImage.FileHeader.Filename))
	passportImageOnePath := fmt.Sprintf("images/passport_images/%s%s", uuid.New().String(), filepath.Ext(passportImageOne.FileHeader.Filename))
	passportImageTwoPath := fmt.Sprintf("images/passport_images/%s%s", uuid.New().String(), filepath.Ext(passportImageTwo.FileHeader.Filename))

	_, err = tx.Exec(context.TODO(), "UPDATE account_profile SET name = $2, profile_image = $3, passport_image_one = $4, passport_image_two = $5, address = $6, phone = $7, country = $8, passport_number = $9, date_of_birth = $10, kin = $11, rank = $12, licence = $13, nationality = $14, completed = true WHERE fk_account_id = $1", accountID, profileName, profileImagePath, passportImageOnePath, passportImageTwoPath, address, phone, country, passportNumber, dateOfBirth, kin, rank, licence, nationality)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "UPDATE account SET profile_completed = true WHERE id = $1", accountID)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = core.PutImage(sess, profileImage.File, profileImage.FileHeader, profileImagePath)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = core.PutImage(sess, passportImageOne.File, passportImageOne.FileHeader, passportImageOnePath)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = core.PutImage(sess, passportImageTwo.File, passportImageTwo.FileHeader, passportImageTwoPath)
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

func (model Model) GetAccountProfileByID(accountID int) (AccountProfile, map[string]string) {
	resp := map[string]string{}
	var account AccountProfile
	err := model.Pool.QueryRow(context.TODO(), "SELECT id, created, Updated, profile_image, passport_image_one, passport_image_two, name, passport_number, date_of_birth, phone, country, address, kin, rank, licence, nationality, completed FROM account_profile WHERE id = $1", accountID).Scan(
		&account.ID, &account.Created, &account.Updated, &account.ProfileImage, &account.PassportImageOne, &account.PassportImageTwo, &account.Name, &account.PassportNumber, &account.DateOfBirth, &account.Phone, &account.Country, &account.Address, &account.Kin, &account.Rank, &account.Licence, &account.Nationality, &account.ProfileCompleted)
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return AccountProfile{}, resp
	}
	return account, resp
}

func (model Model) GetAccountProfileByFK(accountID int) (AccountProfile, map[string]string) {
	resp := map[string]string{}
	var account AccountProfile
	err := model.Pool.QueryRow(context.TODO(), "SELECT id, created, Updated, profile_image, passport_image_one, passport_image_two, name, passport_number, date_of_birth, phone, country, address, kin, rank, licence, nationality, completed FROM account_profile WHERE fk_account_id = $1", accountID).Scan(
		&account.ID, &account.Created, &account.Updated, &account.ProfileImage, &account.PassportImageOne, &account.PassportImageTwo, &account.Name, &account.PassportNumber, &account.DateOfBirth, &account.Phone, &account.Country, &account.Address, &account.Kin, &account.Rank, &account.Licence, &account.Nationality, &account.ProfileCompleted)
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return AccountProfile{}, resp
	}
	return account, resp
}

func (model Model) ResetAccountProfileByID(awsSession *session.Session, core *core.Core, accountID int) map[string]string {
	resp := map[string]string{}
	tx, err := model.Pool.Begin(context.Background())
	defer tx.Rollback(context.TODO())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "UPDATE account_profile SET profile_image = null, passport_image_one = null, passport_image_two = null, name = null, passport_number = null, date_of_birth = null, phone = null, country = null, address = null, kin = null, rank = null, licence = null, ship_experience = null, experience = null, nationality = null, completed = false WHERE id = $1", accountID)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "UPDATE account SET profile_completed = false, profile_validated = false WHERE id = $1", accountID)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	//TODO: delete aws images.

	var accountEmail string
	err = tx.QueryRow(context.TODO(), "SELECT email FROM account WHERE id = $1", accountID).Scan(&accountEmail)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = core.RejecctedEmail(accountEmail)
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
