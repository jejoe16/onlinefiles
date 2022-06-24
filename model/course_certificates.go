package model

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

// swagger:model CourseCertificates
type CourseCertificates struct {
	Course              Course
	AccountCertificates []AccountCertificate
}

// swagger:model CourseCertificateDTO
type CourseCertificateDTO struct {
	DateTime time.Time
}

// swagger:model MontlyCourseAnalytics
type MontlyCourseAnalytics struct {
	Issued    uint32
	Edits     uint32
	Revoked   uint32
	Validated uint32
}

type CoursesCertificates []CourseCertificates

func (model Model) GetCourseCertificatesByCourseID(courseID int) (CourseCertificates, map[string]string) {

	course, resp := model.GetCourse(courseID)
	if len(resp) != 0 {
		return CourseCertificates{}, resp
	}

	rows, err := model.Pool.Query(context.TODO(), "SELECT a.id, a.email, a.profile_image, a.name, a.date_of_birth, a.kin, a.rank, a.licence, a.ship_experience, a.experience, a.nationality,  a.passport_number, a.status, a.temp, a.email_validated, c.uid, c.fk_course_id, c.fk_account_id, c.id, c.created, c.type, c.issued, c.activated, c.accessed FROM certificate c INNER JOIN account a ON c.fk_account_id = a.id WHERE fk_course_id = $1", courseID)
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return CourseCertificates{}, resp
	}
	defer rows.Close()

	var data []AccountCertificate
	for rows.Next() {
		var ac AccountCertificate
		//TODO: find a clean way to have the certificate ID returned
		err = rows.Scan(&ac.Account.ID, &ac.Account.Email, &ac.Account.ProfileImage, &ac.Account.Name, &ac.Account.DateOfBirth, &ac.Account.Kin, &ac.Account.Rank, &ac.Account.Licence, &ac.Account.ShipExperience, &ac.Account.Experience, &ac.Account.Nationality, &ac.Account.Passport, &ac.Account.Status, &ac.Account.Temp, &ac.Account.EmailValidate, &ac.Certificate.UID, &ac.Certificate.CourseID, &ac.Certificate.AccountID, &ac.Certificate.ID, &ac.Certificate.Created, &ac.Certificate.TypeID, &ac.Certificate.Issued, &ac.Certificate.Activated, &ac.Certificate.Accessed)
		if err != nil {
			continue
		}

		data = append(data, ac)
	}

	if err = rows.Err(); err != nil {
		resp["system"] = err.Error()
		return CourseCertificates{}, resp
	}

	returnStruct := CourseCertificates{
		Course:              course,
		AccountCertificates: data,
	}

	return returnStruct, resp
}

func (model Model) GetCourseCertificatesInRange(courseID int, dateTime time.Time) (CourseCertificates, map[string]string) {

	course, resp := model.GetCourse(courseID)
	if len(resp) != 0 {
		return CourseCertificates{}, resp
	}

	rows, err := model.Pool.Query(context.TODO(), "SELECT a.id, a.email, a.profile_imabge, a.name, a.date_of_birth, a.kin, a.rank, a.licence, a.ship_experience, a.experience, a.nationality, a.passport_number, a.status, a.temp, a.email_validated, c.uid, c.fk_course_id, c.fk_account_id, c.id, c.created, c.type, c.issued, c.activated, c.accessed FROM certificate c INNER JOIN account a ON c.fk_account_id = a.id WHERE created >= date_trunc('month', current_date - interval '1' month) and created < date_trunc('month', current_date) AND fk_course_id = $1", courseID)
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return CourseCertificates{}, resp
	}
	defer rows.Close()

	var data []AccountCertificate
	for rows.Next() {
		var ac AccountCertificate
		//TODO: find a clean way to have the certificate ID returned
		err = rows.Scan(&ac.Account.ID, &ac.Account.Email, &ac.Account.ProfileImage, &ac.Account.Name, &ac.Account.DateOfBirth, &ac.Account.Kin, &ac.Account.Rank, &ac.Account.Licence, &ac.Account.ShipExperience, &ac.Account.Experience, &ac.Account.Nationality, &ac.Account.Passport, &ac.Account.Status, &ac.Account.Temp, &ac.Account.EmailValidate, &ac.Certificate.UID, &ac.Certificate.CourseID, &ac.Certificate.AccountID, &ac.Certificate.ID, &ac.Certificate.Created, &ac.Certificate.Type, &ac.Certificate.Issued, &ac.Certificate.Activated, &ac.Certificate.Accessed)
		if err != nil {
			continue
		}

		data = append(data, ac)
	}

	if err = rows.Err(); err != nil {
		resp["system"] = err.Error()
		return CourseCertificates{}, resp
	}

	returnStruct := CourseCertificates{
		Course:              course,
		AccountCertificates: data,
	}

	return returnStruct, resp
}

func (model Model) GetCoursesCertificatesByID(accountID int) (CoursesCertificates, map[string]string) {
	courses, resp := model.GetCoursesByAccountID(accountID)
	var courseCertificates CoursesCertificates
	if len(resp) != 0 {
		return courseCertificates, resp
	}

	for _, course := range courses {
		temp, resp := model.GetCourseCertificatesByCourseID(course.ID)
		if len(resp) != 0 {
			return courseCertificates, resp
		}
		courseCertificates = append(courseCertificates, temp)
	}
	return courseCertificates, resp
}

func (model Model) GetCoursesCertificatesInRange(accountID int, dateTime time.Time) (CoursesCertificates, map[string]string) {
	courses, resp := model.GetCoursesByAccountID(accountID)
	var courseCertificates CoursesCertificates
	if len(resp) != 0 {
		return courseCertificates, resp
	}

	for _, course := range courses {
		temp, resp := model.GetCourseCertificatesInRange(course.ID, dateTime)
		if len(resp) != 0 {
			return courseCertificates, resp
		}
		courseCertificates = append(courseCertificates, temp)
	}
	return courseCertificates, resp
}

func (model Model) GetAllCoursesCertificatesInRange(dateTime time.Time) (CoursesCertificates, map[string]string) {
	courses, resp := model.GetAllCourses()
	var courseCertificates CoursesCertificates
	if len(resp) != 0 {
		return courseCertificates, resp
	}

	for _, course := range courses {
		temp, resp := model.GetCourseCertificatesInRange(course.ID, dateTime)
		if len(resp) != 0 {
			return courseCertificates, resp
		}
		courseCertificates = append(courseCertificates, temp)
	}
	return courseCertificates, resp
}

func (model Model) GetAccountCourseCertificatesCountInRange(accountID int, dateTime time.Time) (MontlyCourseAnalytics, map[string]string) {
	var montlyCourseAnalytics MontlyCourseAnalytics
	courses, resp := model.GetCoursesByAccountID(accountID)
	if len(resp) != 0 {
		return MontlyCourseAnalytics{}, resp
	}

	for _, course := range courses {
		var temp MontlyCourseAnalytics
		err := model.Pool.QueryRow(context.TODO(), "SELECT COUNT(*) FROM certificate WHERE created > date_trunc('month', current_date) AND fk_course_id = $1", course.ID).Scan(&temp.Issued)
		if err != nil {
			resp["system"] = err.Error()
			return MontlyCourseAnalytics{}, resp
		}
		montlyCourseAnalytics.Issued += temp.Issued

		err = model.Pool.QueryRow(context.TODO(), "SELECT COUNT(*) FROM certificate WHERE updated > date_trunc('month', current_date) AND fk_course_id = $1", course.ID).Scan(&temp.Edits)
		if err != nil {
			resp["system"] = err.Error()
			return MontlyCourseAnalytics{}, resp
		}
		montlyCourseAnalytics.Edits += temp.Edits

		// err = model.Pool.QueryRow(context.TODO(), "SELECT COUNT(*) FROM certificate_accessed WHERE created >= date_trunc('month', current_date - interval '1' month) and created < date_trunc('month', current_date) AND fk_course_id = $1", course.ID).Scan(&temp.Edits)
		// if err != nil {
		// 	resp["system"] = err.Error()
		// 	return MontlyCourseAnalytics{}, resp
		// }
		montlyCourseAnalytics.Validated += temp.Validated

	}
	return montlyCourseAnalytics, resp
}
