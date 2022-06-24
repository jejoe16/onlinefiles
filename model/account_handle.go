package model

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

// swagger:model SignIn
type SignIn struct {
	ID            int    `json:"id"`
	Email         string `json:"email" validate:"required,gte=4,lte=320,email,lowercase"`
	Password      string `json:"password" validate:"required,gte=8,lte=156"`
	Temp          bool   `json:"temp"`
	Role          int    `json:"role"`
	EmailValidate bool   `json:"email_validated"`
}

// swagger:model SignInResponse
type SignInResponse struct {
	Token uuid.UUID
}

// swagger:model SignUp
type SignUp struct {
	Email           string `json:"email" validate:"required,gte=4,lte=320,email,lowercase"`
	Password        string `json:"password" validate:"required,gte=8,lte=156"`
	PasswordConfirm string `json:"passwordconfirm" validate:"required,gte=8,lte=156"`
	Temp            bool   `json:"temp"`
	TCAgreed        bool   `json:"tcAgreed"`
}

func (s *SignIn) CheckPasswordHash(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password)) == nil
}

func (model Model) ValidateSignin(email string, password string) (int, map[string]string) {
	resp := map[string]string{}
	var signin SignIn

	err := model.Pool.QueryRow(context.Background(), "SELECT id, password, temp, email_validated, role FROM account WHERE email=$1", strings.ToLower(email)).Scan(&signin.ID, &signin.Password, &signin.Temp, &signin.EmailValidate, &signin.Role)
	if err == pgx.ErrNoRows {
		resp["Email"] = "Wrong email or password"
		return 0, resp
	} else if !signin.CheckPasswordHash(password) {
		resp["Email"] = "Wrong email or password"
		return 0, resp
	} else if signin.Temp {
		resp["Email"] = "You will need to sign up with your account email to continue"
		return 0, resp
	} else if signin.Role <= 1 {
		resp["Email"] = "Please validate your email before sigining in"
		return 0, resp
	}
	//TODO: uncomment once done testing
	// else if !signin.EmailValidate {
	// 	resp["Email"] = "You will need to verify your account to continue with the email sent to you"
	// 	return 0, resp
	// }

	return signin.ID, resp
}

func (model Model) ValidateSignup(email string, password string, passwordconfirm string) (bool, map[string]string) {
	resp := map[string]string{}
	var data SignUp
	if password != passwordconfirm {
		resp["password"] = "Password must match"
		return false, resp
	}

	row := model.Pool.QueryRow(context.Background(), "SELECT email, temp FROM account WHERE email = $1", strings.ToLower(email)).Scan(&data.Email, &data.Temp)

	if !data.Temp {
		if row != pgx.ErrNoRows {
			resp["email"] = "Email Exists"
			return false, resp
		}
	}

	return data.Temp, resp
}
