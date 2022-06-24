package model

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/wes-kay/golang_asset_engine/core"
)

// swagger:model Course
type Course struct {
	UID                   uuid.UUID `json:"uid"`
	ID                    int       `json:"id"`
	Created               time.Time `json:"created"`
	Image                 string    `json:"image"`
	Description           *string   `json:"description" validate:"required,lte=500"`
	CourseName            string    `json:"course_name" validate:"required,gte=3,lte=100"`
	CertificateName       string    `json:"certificate_name" validate:"required,gte=3,lte=100"`
	AdditionalDescription *string   `json:"additional_description" validate:"lte=500"`
	Expire                int       `json:"expire" validate:"required,gte=0,numeric"`
	CourseCategory        string    `json:"name"`
	CourseCategoryID      int       `json:"fk_course_category_id" validate:"required,gte=0,numeric"`
	// Address               string    `json:"address" validate:"required,lte=500"`
	// Country               string    `json:"country" validate:"required,lte=500"`
	// CourseID              int       `json:"course_id" validate:"required,gte=0,numeric"`
}

// swagger:model CourseDTO
type CourseDTO struct {
	Created               time.Time `json:"created"`
	Image                 *string   `json:"image"`
	Description           *string   `json:"description" validate:"required,lte=500"`
	CourseName            string    `json:"course_name" validate:"required,gte=3,lte=100"`
	CertificateName       string    `json:"certificate_name" validate:"required,gte=3,lte=100"`
	AdditionalDescription *string   `json:"additional_description" validate:"lte=500"`
	Expire                int       `json:"expire" validate:"required,gte=0,numeric"`
	CourseCategory        int       `json:"fk_course_category_id" validate:"required,gte=0,numeric"`
}

type Courses []Course

func (model Model) CreateCourse(sess *session.Session, accountID int, file multipart.File, handler *multipart.FileHeader, certificate_name string, description *string, courseName string, additionalDescription *string, expire int, courseCategory int) map[string]string {
	resp := map[string]string{}

	tx, err := model.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	fileName := "images/" + uuid.New().String() + filepath.Ext(handler.Filename)

	_, err = tx.Exec(context.TODO(), "INSERT INTO course (fk_account_id, certificate_name, description, course_name, additional_description, expire, image_source, fk_course_category_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8)", accountID, certificate_name, description, courseName, additionalDescription, expire, fileName, courseCategory)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = core.PutImage(sess, file, handler, fileName)
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

//TODO: Get count of cert users
func (model Model) GetCourse(id int) (Course, map[string]string) {
	resp := map[string]string{}
	var course Course
	err := model.Pool.QueryRow(context.TODO(), "SELECT c.id, c.certificate_name, c.created, c.description, c.course_name, c.additional_description, c.expire, c.image_source, cc.name FROM course c INNER JOIN course_category_lookup cc ON c.fk_course_category_id = cc.id WHERE c.id = $1", id).Scan(&course.ID, &course.CertificateName, &course.Created, &course.Description, &course.CourseName, &course.AdditionalDescription, &course.Expire, &course.Image, &course.CourseCategory)
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return Course{}, resp
	}
	return course, resp
}

func (model Model) UpdateCourse(sess *session.Session, accountID int, file multipart.File, handler *multipart.FileHeader, certicicateName string, description *string, courseName string, additionalDescription *string, expire int, courseCategory int) map[string]string {
	resp := map[string]string{}

	tx, err := model.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	fileName := "images/" + uuid.New().String() + filepath.Ext(handler.Filename)

	_, err = tx.Exec(context.TODO(), "UPDATE course SET certificate_name = $2, description = $3, course_name =$4, additional_description = $5, expire = $6, image_source = $7, fk_course_category_id = $8 WHERE id = $1", accountID, certicicateName, description, courseName, additionalDescription, expire, fileName, courseCategory)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	err = core.PutImage(sess, file, handler, fileName)
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

func (model Model) GetCoursesByAccountID(accountID int) (Courses, map[string]string) {
	resp := map[string]string{}
	rows, err := model.Pool.Query(context.TODO(), "SELECT c.id, c.certificate_name, c.created, c.description, c.course_name, c.additional_description, c.expire, c.image_source, cc.id, cc.name FROM course c INNER JOIN course_category_lookup cc ON c.fk_course_category_id = cc.id WHERE c.fk_account_id = $1", accountID)
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return Courses{}, resp
	}
	defer rows.Close()

	var data Courses
	for rows.Next() {
		var course Course
		err = rows.Scan(&course.ID, &course.CertificateName, &course.Created, &course.Description, &course.CourseName, &course.AdditionalDescription, &course.Expire, &course.Image, &course.CourseCategoryID, &course.CourseCategory)
		if err != nil {
			continue
		}

		data = append(data, course)
	}

	if err = rows.Err(); err != nil {
		resp["system"] = err.Error()
		return Courses{}, resp
	}

	return data, resp
}

func (model Model) GetAllCourses() (Courses, map[string]string) {
	resp := map[string]string{}
	rows, err := model.Pool.Query(context.TODO(), "SELECT c.id, c.certificate_name, c.created, c.description, c.course_name, c.additional_description, c.expire, c.image_source, cc.name FROM course c INNER JOIN course_category_lookup cc ON c.fk_course_category_id = c.id")
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return Courses{}, resp
	}
	defer rows.Close()

	var data Courses
	for rows.Next() {
		var course Course
		err = rows.Scan(&course.ID, &course.CertificateName, &course.Created, &course.Description, &course.CourseName, &course.AdditionalDescription, &course.Expire, &course.Image, &course.CourseCategory)
		if err != nil {
			continue
		}

		data = append(data, course)
	}

	if err = rows.Err(); err != nil {
		resp["system"] = err.Error()
		return Courses{}, resp
	}

	return data, resp
}

func (model Model) AccountHasCourse(accountID int, courseID int) (bool, map[string]string) {
	resp := map[string]string{}
	res, err := model.Pool.Exec(context.TODO(), "SELECT COUNT(1) FROM course WHERE fk_account_id = $1 AND id = $2", accountID, courseID)
	if err != nil {
		resp["system"] = err.Error()
		return false, resp
	}
	if res.RowsAffected() == 0 {
		resp["system"] = "no rows updated"
		return false, resp
	}

	return true, resp
}

func (model Model) DeleteCourse(courseID int) map[string]string {
	resp := map[string]string{}

	tx, err := model.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	res, err := tx.Exec(context.TODO(), "DELETE FROM certificate WHERE fk_course_id = $1", courseID)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	if res.RowsAffected() == 0 {
		resp["system"] = "no rows updated"
		return resp
	}

	res, err = tx.Exec(context.TODO(), "DELETE FROM course WHERE id = $1", courseID)
	if err != nil {
		resp["system"] = err.Error()
		return resp
	}

	if res.RowsAffected() == 0 {
		resp["system"] = "no rows updated"
		return resp
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		resp["system"] = err.Error()
	}

	return resp
}
