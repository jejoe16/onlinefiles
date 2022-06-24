package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/wes-kay/golang_asset_engine/core"
)

// swagger:model Certificate
type Certificate struct {
	ID              int
	UID             uuid.UUID
	Created         time.Time
	Updated         time.Time
	CourseID        int
	AccountID       int
	CertificateType int
	Type            string
	TypeID          int
	Issued          time.Time
	Activated       bool
	Accessed        int
}

// swagger:model CertificateDTO
type CertificateDTO struct {
	AccountEmail string    `json:"account_email" validate:"required,email"`
	Type         int       `json:"type" validate:"required,gte=0,numeric"`
	Issued       time.Time `json:"issued" validate:"required"`
	Activated    bool      `json:"activated" validate:"required"`
}

// swagger:model CertificateEditDTO
type CertificateEditDTO struct {
	Type      int       `json:"type" validate:"required,gte=0,numeric"`
	Issued    time.Time `json:"issued" validate:"required"`
	Activated bool      `json:"activated" validate:"required"`
}

// swagger:model CourseCertificate
type CourseCertificate struct {
	Course          Course
	UserAccount     Account
	ProviderAccount Account
	Certificate     Certificate
}

type AccountCertificate struct {
	Account     Account
	Certificate Certificate
}

type Certificates []Certificate

func (model Model) CreateCertificate(email string, fkCourseID int, providerName *string, certType int, issued time.Time, activated bool, core *core.Core) map[string]string {
	resp := map[string]string{}
	tx, err := model.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	account, _ := model.GetAccountByEmail(email)

	if account == (Account{}) {
		err = tx.QueryRow(context.Background(), "INSERT INTO account(email, password, temp) VALUES($1, $2, $3) RETURNING id", email, "", true).Scan(&account.ID)
		if err != nil {
			resp["system"] = err.Error()
			return resp
		}

		_, err = tx.Exec(context.TODO(), "INSERT INTO account_profile(fk_account_id) VALUES($1)", account.ID)
		if err != nil {
			resp["system"] = err.Error()
			return resp
		}

		uid := uuid.New()

		_, err = tx.Exec(context.Background(), "INSERT INTO email_verification(fk_account_id, code) VALUES($1, $2)", account.ID, uid.String())
		if err != nil {
			resp["system"] = err.Error()
			return resp
		}
		//TODO: After testing
		core.ActivationEmail(email, uid.String())
		if err != nil {
			resp["system"] = err.Error()
			return resp
		}
	}

	_, err = tx.Exec(context.TODO(), "INSERT INTO certificate (fk_account_id, fk_course_id, type, issued, activated) VALUES($1, $2, $3, $4, $5)", account.ID, fkCourseID, certType, issued, activated)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "INSERT INTO alert(fk_account_id, value) VALUES($1, $2)", account.ID, "You have a new certificate")
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = core.CertificateEmail(email)
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

func (model Model) GetCertificate(id int) (Certificate, map[string]string) {
	resp := map[string]string{}
	var cert Certificate

	err := model.Pool.QueryRow(context.TODO(), "SELECT id, fk_course_id, fk_account_id, created, type, issued, activated, accessed FROM certificate WHERE id = $1", id).Scan(&cert.ID, &cert.CourseID, &cert.AccountID, &cert.Created, &cert.CertificateType, &cert.Issued, &cert.Activated, &cert.Accessed)
	if err != nil {
		if err == pgx.ErrNoRows {
			resp["system"] = "no certificate found"
			return Certificate{}, resp
		}
		resp["system"] = err.Error()
		return Certificate{}, resp
	}

	return cert, resp
}

func (model Model) GetCertificatebyUid(uid uuid.UUID) (CourseCertificate, map[string]string) {
	resp := map[string]string{}

	var data CourseCertificate

	err := model.Pool.QueryRow(context.TODO(), "SELECT a.email, a.profile_image, a.phone, a.name, a.kin, a.rank, a.licence, a.ship_experience, a.nationality, a.date_of_birth, a.passport_number, cr.uid, cr.course_name, cr.certificate_name, cr.description, cr.additional_description, cr.expire, cr.image_source, ct.uid, ct.created, cc.name, ct.issued, ct.activated, cc.name, pa.name, pa.country, pa.address FROM certificate ct INNER JOIN course cr ON ct.fk_course_id = cr.id INNER JOIN account a ON ct.fk_account_id = a.id INNER JOIN course_category_lookup cc ON cr.fk_course_category_id = cc.id INNER JOIN account pa ON cr.fk_account_id = pa.id WHERE ct.uid = $1", uid).Scan(&data.UserAccount.Email, &data.UserAccount.ProfileImage, &data.UserAccount.Phone, &data.UserAccount.Name, &data.UserAccount.Kin, &data.UserAccount.Rank, &data.UserAccount.Licence, &data.UserAccount.ShipExperience, &data.UserAccount.Nationality, &data.UserAccount.DateOfBirth, &data.UserAccount.Passport, &data.Course.UID, &data.Course.CourseName, &data.Course.CertificateName, &data.Course.Description, &data.Course.AdditionalDescription, &data.Course.Expire, &data.Course.Image, &data.Certificate.UID, &data.Certificate.Created, &data.Certificate.Type, &data.Certificate.Issued, &data.Certificate.Activated, &data.Course.CourseCategory, &data.ProviderAccount.Name, &data.ProviderAccount.Country, &data.ProviderAccount.Address)

	if err != nil {
		if err == pgx.ErrNoRows {
			resp["system"] = "no certificate found"
			return CourseCertificate{}, resp
		}
		resp["system"] = err.Error()
		return CourseCertificate{}, resp
	}

	if !data.Certificate.Activated {
		resp["error"] = "Certificate not activated"
		return CourseCertificate{}, resp
	}

	_, err = model.Pool.Exec(context.TODO(), "UPDATE certificate SET accessed = accessed + 1 WHERE id = $1", data.Certificate.ID)
	if err != nil {
		resp["error"] = err.Error()
		return CourseCertificate{}, resp
	}

	return data, resp
}

func (model Model) GetAllCertificates() (Certificates, map[string]string) {
	resp := map[string]string{}
	rows, err := model.Pool.Query(context.TODO(), "SELECT a.id, a.email, a.kin, a.rank, a.licence, a.ship_experience, a.experience, a.nationality, c.id, c.fk_course_id, c.fk_account_id, c.created, c.type, c.issued, c.activated, c.accessed FROM certificate c INNER JOIN account a ON a.id = c.fk_account_id ORDER BY id")
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return Certificates{}, resp
	}
	defer rows.Close()

	var data Certificates
	for rows.Next() {
		var cert Certificate
		err = rows.Scan(&cert.ID, &cert.CourseID, &cert.AccountID, &cert.Created, &cert.Type, &cert.Issued, &cert.Activated, &cert.Accessed)

		if err != nil {
			continue
		}

		data = append(data, cert)
	}

	if err = rows.Err(); err != nil {
		resp["system"] = err.Error()
		return Certificates{}, resp
	}

	return data, resp
}

func (model Model) UpdateAccountCertificate(accountID int, certificateType int, issued time.Time, activated bool) map[string]string {
	resp := map[string]string{}
	tx, err := model.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	_, err = tx.Exec(context.TODO(), "UPDATE certificate SET type = $2, issued = $3, activated = $4 WHERE id = $1", accountID, certificateType, issued, activated)
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

//TODO: Update full certificate
func (model Model) UpdateCertificate(id int, certType int, issued string, activated bool) map[string]string {
	resp := map[string]string{}
	res, err := model.Pool.Exec(context.TODO(), "UPDATE certificate SET type = $2, issued = $3, activated = $4  WHERE id = $1", id, certType, issued, activated)
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

func (model Model) DeleteCertificate(certificateID int) map[string]string {
	resp := map[string]string{}
	res, err := model.Pool.Exec(context.TODO(), "DELETE FROM certificate WHERE id = $1", certificateID)
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

//TODO: get account Cert info then get cert
func (model Model) GetAccountCertificates(accountID int) ([]CourseCertificate, map[string]string) {
	resp := map[string]string{}

	rows, err := model.Pool.Query(context.TODO(), "SELECT a.email, a.date_of_birth, a.passport_number, a.profile_image, a.kin, a.rank, a.licence, a.ship_experience, a.nationality, cr.uid, cr.course_name, cr.certificate_name, cr.description, cr.additional_description, cr.expire, cr.image_source, ct.uid, ct.type, ct.created, cc.name, ct.issued, ct.activated FROM certificate ct INNER JOIN course cr ON ct.fk_course_id = cr.id INNER JOIN course_category_lookup cc ON cr.fk_course_category_id = cc.id INNER JOIN account a ON ct.fk_account_id = a.id WHERE ct.fk_account_id = $1 ORDER BY ct.type", accountID)
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return nil, resp
	}
	defer rows.Close()

	var data []CourseCertificate
	for rows.Next() {
		var cert CourseCertificate
		err = rows.Scan(&cert.UserAccount.Email, &cert.UserAccount.DateOfBirth, &cert.UserAccount.Passport, &cert.UserAccount.ProfileImage, &cert.UserAccount.Kin, &cert.UserAccount.Rank, &cert.UserAccount.Licence, &cert.UserAccount.ShipExperience, &cert.UserAccount.Nationality, &cert.Course.UID, &cert.Course.CourseName, &cert.Course.CertificateName, &cert.Course.Description, &cert.Course.AdditionalDescription, &cert.Course.Expire, &cert.Course.Image, &cert.Certificate.UID, &cert.Certificate.TypeID, &cert.Certificate.Created, &cert.Certificate.Type, &cert.Certificate.Issued, &cert.Certificate.Activated)
		if err != nil {
			continue
		}

		data = append(data, cert)
	}

	if err = rows.Err(); err != nil {
		resp["system"] = err.Error()
		return nil, resp
	}

	return data, resp
}
