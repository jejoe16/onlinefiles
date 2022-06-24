package model

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// swagger:model CertificateType
type CertificateType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// swagger:model CertificateTypes
type CertificateTypes []CertificateType

// swagger:model AccountStatus
type AccountStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// swagger:model AccountStatuses
type AccountStatuses []AccountStatus

// swagger:model AccountRole
type AccountRole struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// swagger:model AccountRoles
type AccountRoles []AccountRole

// swagger:model CourseCategory
type CourseCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// swagger:model CourseCategories
type CourseCategories []CourseCategory

func (model Model) GetCertificateTypesLookup() (CertificateTypes, map[string]string) {
	resp := map[string]string{}
	rows, err := model.Pool.Query(context.TODO(), "SELECT id, name FROM certificate_type_lookup")
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return nil, resp
	}
	defer rows.Close()

	var data CertificateTypes
	for rows.Next() {
		var certType CertificateType
		err = rows.Scan(&certType.ID, &certType.Name)
		if err != nil {
			continue
		}

		data = append(data, certType)
	}

	if err = rows.Err(); err != nil {
		resp["system"] = err.Error()
		return nil, resp
	}

	return data, resp
}

func (model Model) GetAccountStatusLookups() (AccountStatuses, map[string]string) {
	resp := map[string]string{}
	rows, err := model.Pool.Query(context.TODO(), "SELECT id, name FROM account_status_lookup")
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return nil, resp
	}
	defer rows.Close()

	var datas AccountStatuses
	for rows.Next() {
		var data AccountStatus
		err = rows.Scan(&data.ID, &data.Name)
		if err != nil {
			continue
		}

		datas = append(datas, data)
	}

	if err = rows.Err(); err != nil {
		resp["system"] = err.Error()
		return nil, resp
	}

	return datas, resp
}

func (model Model) GetAccountRoleLookups() (AccountRoles, map[string]string) {
	resp := map[string]string{}
	rows, err := model.Pool.Query(context.TODO(), "SELECT id, name FROM account_role_lookup")
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return nil, resp
	}
	defer rows.Close()

	var datas AccountRoles
	for rows.Next() {
		var data AccountRole
		err = rows.Scan(&data.ID, &data.Name)
		if err != nil {
			continue
		}

		datas = append(datas, data)
	}

	if err = rows.Err(); err != nil {
		resp["system"] = err.Error()
		return nil, resp
	}

	return datas, resp
}

func (model Model) GetCourseCategoryLookups() (CourseCategories, map[string]string) {
	resp := map[string]string{}
	rows, err := model.Pool.Query(context.TODO(), "SELECT id, name FROM course_category_lookup")
	if err == pgx.ErrNoRows {
		resp["system"] = err.Error()
		return nil, resp
	}
	defer rows.Close()

	var datas CourseCategories
	for rows.Next() {
		var data CourseCategory
		err = rows.Scan(&data.ID, &data.Name)
		if err != nil {
			continue
		}

		datas = append(datas, data)
	}

	if err = rows.Err(); err != nil {
		resp["system"] = err.Error()
		return nil, resp
	}

	return datas, resp
}
