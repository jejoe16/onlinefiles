package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wes-kay/golang_asset_engine/model"
)

func (a *App) InitializeRoutes(router *mux.Router) {

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	router.HandleFunc("/v1/signin", a.SignIn).Methods("POST")

	// router.HandleFunc("/v1/signout", a.Signout).Methods("POST")

	router.HandleFunc("/v1/signup", a.Signup).Methods("POST")
	router.HandleFunc("/v1/signup/{uid}", a.VerifyEmail).Methods("POST")
	router.HandleFunc("/v1/signup/{email}/resend", a.VerifyEmailResend).Methods("POST")

	router.HandleFunc("/v1/view/certificate/{uid}", a.ViewCourseCertificate).Methods("GET")
	router.HandleFunc("/v1/download/certificate/{uid}", a.DownloadCourseCertificate).Methods("GET")

	// router.HandleFunc("/v1/view/qr/{id}", NotImplemented).Methods("GET")

	s := router.PathPrefix("/").Subrouter()
	s.Use(a.Verify)

	s.HandleFunc("/v1/account", a.GetAccountInfo).Methods("GET")
	// s.HandleFunc("/v1/dashboard", NotImplemented).Methods("GET")

	s.HandleFunc("/v1/account/{id}", a.GetAccount).Methods("GET")
	s.HandleFunc("/v1/account/{id}", a.UpdateAccount).Methods("PUT")
	s.HandleFunc("/v1/account/{id}", a.DeleteAccount).Methods("DELETE")

	s.HandleFunc("/v1/profile/{id}", a.GetAccountProfile).Methods("GET")
	s.HandleFunc("/v1/profile/{id}", a.UpdateAccountProfile).Methods("PUT")

	s.HandleFunc("/v1/profile/validate", a.ValidateAccountProfile).Methods("POST")
	s.HandleFunc("/v1/profile/reject", a.RejectAccountProfile).Methods("POST")

	s.HandleFunc("/v1/course/{id}", a.GetCourse).Methods("GET")
	s.HandleFunc("/v1/course", a.CreateCourse).Methods("POST")
	s.HandleFunc("/v1/course/{id}", a.UpdateCourse).Methods("PUT")
	s.HandleFunc("/v1/course/{id}", a.DeleteCourse).Methods("DELETE")
	//Change
	s.HandleFunc("/v1/courses/{id}", a.GetCourses).Methods("GET")

	s.HandleFunc("/v1/course-information/{id}", a.GetCourses).Methods("GET")

	s.HandleFunc("/v1/certificate/{id}", a.GetCertificate).Methods("GET")
	s.HandleFunc("/v1/certificate/{id}", a.CreateCertificate).Methods("POST")
	s.HandleFunc("/v1/certificate/{id}", a.UpdateCertificate).Methods("PUT")
	s.HandleFunc("/v1/certificate/{id}", a.DeleteCertificate).Methods("DELETE")

	//change
	s.HandleFunc("/v1/certificates/{id}", a.GetCertificates).Methods("GET")

	s.HandleFunc("/v1/certificates/account/{id}", a.GetAccountCertificates).Methods("GET")
	// s.HandleFunc("/v1/certificate/account/{id}", a.UpdateAccountCertificate).Methods("PUT")
	s.HandleFunc("/v1/certificates/course/{id}", a.GetCourseCertificates).Methods("GET")
	s.HandleFunc("/v1/certificates/course/{id}", NotImplemented).Methods("PUT")

	s.HandleFunc("/v1/admin/accounts", a.GetAllAccounts).Methods("GET")
	s.HandleFunc("/v1/admin/account/{id}", a.UpdateAdminAccount).Methods("PUT")
	s.HandleFunc("/v1/admin/account/{id}", a.DeleteAdminAccount).Methods("DELETE")
	s.HandleFunc("/v1/admin/courses", a.GetAllCourses).Methods("GET")
	s.HandleFunc("/v1/admin/course/{id}/analytics", a.GetAllAccountCoursesAnalytics).Methods("GET")

	s.HandleFunc("/v1/admin/certificates", a.GetAllCertificates).Methods("GET")
	// s.HandleFunc("/v1/admin/courses/{id}", NotImplemented).Methods("PUT")
	// s.HandleFunc("/v1/admin/certificate/{id}", NotImplemented).Methods("PUT")

	s.HandleFunc("/v1/lookup/certificate-type", a.GetCertificateTypes).Methods("Get")
	s.HandleFunc("/v1/lookup/account-status", a.GetAccountStatuses).Methods("Get")
	s.HandleFunc("/v1/lookup/account-role", a.GetAccountRoles).Methods("Get")
	s.HandleFunc("/v1/lookup/course-categories", a.GetCourseCategories).Methods("Get")
}

// swagger:operation POST /signin
// @Summary      Sign Into Account
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        object  formData model.SignIn  false  "Sign in"
// @Success      200  {object}  model.SignInResponse
// @Router       /signin [post]
func (a *App) SignIn(w http.ResponseWriter, r *http.Request) {
	var data model.SignIn
	err := DecodeForm(r, &data)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	id, resp := a.Model.ValidateSignin(data.Email, data.Password)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	token, resp := a.NewToken(id, GetIP(r))
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}
	respondWithJSON(w, http.StatusOK, token)
}

func (a *App) Signout(w http.ResponseWriter, r *http.Request) {
	// resp := map[string]string{}
	// session, err := a.Sessions.Get(r, "session-token")
	// if err != nil {
	// 	resp["system"] = err.Error()
	// 	respondWithJSON(w, http.StatusInternalServerError, resp)
	// 	return
	// }

	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "remember-token",
	// 	Value:   "",
	// 	Secure:  true,
	// 	Expires: time.Now(),
	// })

	// session.Options.MaxAge = 0
}

// swagger:operation POST /signup
// @Summary      Creates an account on the service
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        object  formData model.SignUp  false  "Sign up"
// @Success      200
// @Router       /signup [post]
func (a *App) Signup(w http.ResponseWriter, r *http.Request) {
	var data model.SignUp
	err := DecodeForm(r, &data)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusNotAcceptable, resp)
		return
	}

	if !data.TCAgreed {
		resp["tcAgreed"] = "Must confirm to Terms & Conditions"
		respondWithJSON(w, http.StatusNotAcceptable, resp)
		return
	}

	temp, resp := a.Model.ValidateSignup(data.Email, data.Password, data.PasswordConfirm)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusNotAcceptable, resp)
		return
	}

	if temp {
		resp = a.Model.UpdateAccountPassword(data.Email, data.Password)
		if len(resp) != 0 {
			respondWithJSON(w, http.StatusInternalServerError, resp)
			return
		}
	} else {
		resp = a.Model.CreateAccount(data.Email, data.Password, a.Core)
		if len(resp) != 0 {
			respondWithJSON(w, http.StatusInternalServerError, resp)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, "success")
}

// swagger:operation POST /v1/signup/{uid} verify
// @Summary      Verifies an email with a UID sent to the account email
// @Tags         verify
// @Accept       json
// @Produce      json
// @Param        uid query string  true  "uid from email"
// @Success      200
// @Router       /v1/signup/{uid} [post]
func (a *App) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := uuid.Parse(vars["uid"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data := PostUid{
		UID: uid,
	}

	resp := ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.EmailVerification(data.UID.String())
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	respondWithJSON(w, 200, "success")
}

// swagger:operation POST /v1/signup/{email}/resend
// @Summary  	 Loads certificate from uid
// @Tags         certificate
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        email query string  true  "email of user"
// @Success      200  {object}  model.CourseCertificate
// @Router       /v1/signup/{email}/resend [POST]
func (a *App) VerifyEmailResend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	data := PostEmail{
		Email: email,
	}

	resp := ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.AccountActivationResend(a.Core, email)
	if len(resp) != 0 {
		//TODO: handle error silent
	}

	respondWithJSON(w, 200, data)
}

// swagger:operation GET /v1/view/certificate/{uid}
// @Summary  	 Loads certificate from uid
// @Tags         certificate
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        uid query string  true  "uid from url"
// @Success      200  {object}  model.CourseCertificate
// @Router       /v1/view/certificate/{uid} [GET]
func (a *App) ViewCourseCertificate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := uuid.Parse(vars["uid"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	postData := PostUid{
		UID: uid,
	}

	resp := ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)

		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	data, resp := a.Model.GetCertificatebyUid(postData.UID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	respondWithJSON(w, 200, data)
}

// func GeneratePDF(w http.ResponseWriter, data model.CourseCertificate) {
// 	pdf := fpdf.New("P", "mm", "A4", "")
// 	pdf.AddPage()
// 	pdf.SetFont("Arial", "B", 16)
// 	pdf.Cell(40, 10, "Hello, world")
// 	pdf.Output(w)
// }

// func (a *App) Test(w http.ResponseWriter, r *http.Request) {

// }

// swagger:operation GET /v1/download/certificate/{uid}
// @Summary  	 Loads certificate from uid
// @Tags         certificate
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        uid query string  true  "uid from url"
// @Success      200  {object}  model.CourseCertificate
// @Router       /v1/download/certificate/{uid} [GET]
func (a *App) DownloadCourseCertificate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := uuid.Parse(vars["uid"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	postData := PostUid{
		UID: uid,
	}

	resp := ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)

		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	data, resp := a.Model.GetCertificatebyUid(postData.UID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	fileName := fmt.Sprintf("attachment; filename=%s.pdf", strings.ReplaceAll(data.Course.CertificateName, " ", "_"))

	w.Header().Set("Content-Disposition", fileName)
	err = a.Core.GeneratePDF(w, data.Course.Image, *data.UserAccount.ProfileImage, data.Course.CourseName, data.Course.CertificateName, data.UserAccount.Name, data.Course.CourseCategory, data.Certificate.Created, data.Course.Expire, data.Certificate.UID.String(), data.Certificate.ID, data.UserAccount.Passport, data.UserAccount.DateOfBirth)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}
}

// swagger:operation GET /v1/account get account info
// @Summary  	 Get's account info using the session token
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Success      200  {object}  model.CourseCertificate
// @Router       /v1/account [GET]
func (a *App) GetAccountInfo(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}

	id, ok := r.Context().Value(contextKeyRequestID).(int)
	if !ok {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	account, resp := a.Model.GetAccountById(id)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, account)
}

// func (a *App) DisplayCertificate(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	uid, err := uuid.FromBytes([]byte(vars["id"]))
// 	if err != nil {
// 		respondWithJSON(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	data := model.PostUid{
// 		UID: uid,
// 	}

// 	resp := ValidatePostData(&data)
// 	if len(resp) != 0 {
// 		respondWithJSON(w, http.StatusBadRequest, resp)
// 		return
// 	}

// 	cert, resp := a.Model.GetCertificate(data.UID)
// 	if len(resp) != 0 {
// 		respondWithJSON(w, http.StatusBadRequest, resp)
// 		return
// 	}
// 	respondWithJSON(w, 200, cert)
// }

// swagger:operation GET /v1/account/{id}
// @Summary  	 Gets an account by id
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "account id"
// @Success      200  {object}  model.Account
// @Router       /v1/account/{id} [GET]
func (a *App) GetAccount(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	data := PostID{
		ID: id,
	}

	resp = ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	account, resp := a.Model.GetAccountById(data.ID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, account)
}

// swagger:operation PUT /v1/account/{id}
// @Summary  	 Updates an account by ID
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "account id"
// @Success      200
// @Router       /v1/account/{id} [PUT]
func (a *App) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	postData := PostID{
		ID: id,
	}

	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	accountProfile, resp := a.Model.GetAccountProfileByID(postData.ID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.UpdateAccountByID(a.AWS, a.Core, postData.ID, accountProfile.PassportImageOne, accountProfile.Name, accountProfile.Address, accountProfile.Country, accountProfile.Phone, accountProfile.PassportNumber, accountProfile.DateOfBirth, accountProfile.Kin, accountProfile.Rank, accountProfile.Licence, accountProfile.Nationality)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "success")
}

// swagger:operation GET /v1/profile/{id}
// @Summary  	 Gets an account by id
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "account id"
// @Success      200  {object}  model.AccountProfile
// @Router       /v1/profile/{id} [GET]
func (a *App) GetAccountProfile(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	data := PostID{
		ID: id,
	}

	resp = ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	account, resp := a.Model.GetAccountProfileByFK(data.ID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, account)
}

// swagger:operation PUT /v1/profile/{id}
// @Summary  	 Updates an account by ID
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "account id"
// @Param        object  formData model.AccountProfileDTO  false  "Account data"
// @Success      200
// @Router      /v1/profile/{id} [PUT]
func (a *App) UpdateAccountProfile(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	postData := PostID{
		ID: id,
	}

	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	var data model.AccountProfileDTO

	data.ProfileImage.File, data.ProfileImage.FileHeader, err = r.FormFile("ProfileImage")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data.PassportImageOne.File, data.PassportImageOne.FileHeader, err = r.FormFile("PassportImageOne")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data.PassportImageTwo.File, data.PassportImageTwo.FileHeader, err = r.FormFile("PassportImageTwo")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = DecodeForm(r, &data)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	emptyStr := ""

	data.Phone = &emptyStr
	data.Kin = &emptyStr
	data.Licence = &emptyStr
	data.Rank = &emptyStr
	data.Address = &emptyStr
	data.Country = &emptyStr

	resp = ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.UpdateAccountProfileByID(a.AWS, postData.ID, data.ProfileImage, data.PassportImageOne, data.PassportImageTwo, data.Name, data.Address, data.Country, data.Phone, data.PassportNumber, data.DateOfBirth, data.Kin, data.Rank, data.Licence, data.Nationality)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "success")
}

// swagger:operation POST /v1/profile/validate
// @Summary  	 Updates an account by ID
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "account id"
// @Success      200
// @Router      /v1/profile/validate [POST]
func (a *App) ValidateAccountProfile(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}

	var data model.AccountProfileVerifyDTO
	err := DecodeForm(r, &data)
	if err != nil {
		resp["system"] = err.Error()
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	resp = ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	accountProfile, resp := a.Model.GetAccountProfileByID(data.ID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.UpdateAccountByID(a.AWS, a.Core, data.ID, accountProfile.ProfileImage, accountProfile.Name, accountProfile.Address, accountProfile.Country, accountProfile.Phone, accountProfile.PassportNumber, accountProfile.DateOfBirth, accountProfile.Kin, accountProfile.Rank, accountProfile.Licence, accountProfile.Nationality)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "success")
}

// swagger:operation POST /v1/profile/reject
// @Summary  	 Updates an account by ID
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "account id"
// @Success      200
// @Router      /v1/profile/reject [POST]
func (a *App) RejectAccountProfile(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	var data model.AccountProfileRejectDTO
	err := DecodeForm(r, &data)
	if err != nil {
		resp["system"] = err.Error()
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.ResetAccountProfileByID(a.AWS, a.Core, data.ID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "success")
}

// swagger:operation DELETE /v1/account/{id}
// @Summary  	 Deletes an account by ID
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "account id"
// @Success      200
// @Router       /v1/account/{id} [DELETE]
func (a *App) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	data := PostID{
		ID: id,
	}

	resp = ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.DeleteAccount(data.ID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "success")
}

// swagger:operation GET /v1/course/{id}
// @Summary  	 Gets a course by ID
// @Tags         course
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Course ID"
// @Success      200  {object}  model.Course
// @Router       /v1/course/{id} [GET]
func (a *App) GetCourse(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	data := PostID{
		ID: id,
	}

	resp = ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	course, resp := a.Model.GetCourse(data.ID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, course)
}

// swagger:operation POST /v1/course
// @Summary  	 Creates a new course using session ID
// @Tags         course
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        object formData model.CourseDTO  false  "Course"
// @Success      200
// @Router       /v1/course/{id} [POST]
func (a *App) CreateCourse(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("Image")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var data model.CourseDTO
	err = DecodeForm(r, &data)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	id, ok := r.Context().Value(contextKeyRequestID).(int)
	if !ok {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	resp = a.Model.CreateCourse(a.AWS, id, file, handler, data.CertificateName, data.Description, data.CourseName, data.AdditionalDescription, data.Expire, data.CourseCategory)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "success")
}

// swagger:operation DELETE /v1/course/{id}
// @Summary  	 Delete a course by ID
// @Tags         course
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Course ID"
// @Success      200
// @Router       /v1/course/{id} [DELETE]
func (a *App) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}
	postData := PostID{
		ID: id,
	}
	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.DeleteCourse(id)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "success")
}

// swagger:operation PUT /v1/course/{id}
// @Summary  	 Update a course by ID
// @Tags         course
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Course ID"
// @Param        object formData model.CourseDTO  false  "Course"
// @Success      200
// @Router       /v1/course/{id} [PUT]
func (a *App) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}
	postData := PostID{
		ID: id,
	}
	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	file, handler, err := r.FormFile("Image")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var data model.CourseDTO
	err = DecodeForm(r, &data)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	resp = ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.UpdateCourse(a.AWS, postData.ID, file, handler, data.CertificateName, data.Description, data.CourseName, data.AdditionalDescription, data.Expire, data.CourseCategory)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "success")
}

// swagger:operation GET /v1/courses/{id}
// @Summary  	 Returns all courses under account
// @Tags         course
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Course ID"
// @Success      200  {array}   model.Course
// @Router       /v1/courses/{id} [GET]
func (a *App) GetCourses(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}
	data := PostID{
		ID: id,
	}
	postResponse := ValidatePostData(&data)
	if len(postResponse) != 0 {
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	// id, ok := r.Context().Value(contextKeyRequestID).(int)
	// if !ok {
	// 	respondWithJSON(w, http.StatusInternalServerError, resp)
	// 	return
	// }

	course, resp := a.Model.GetCoursesByAccountID(id)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, course)
}

// func (a *App) GetCoursesinformation(w http.ResponseWriter, r *http.Request) {
// 	resp := map[string]string{}
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		resp["id"] = "must be a number"
// 		respondWithJSON(w, http.StatusBadRequest, resp)
// 		return
// 	}
// 	data := PostID{
// 		ID: id,
// 	}
// 	postResponse := ValidatePostData(&data)
// 	if len(postResponse) != 0 {
// 		respondWithJSON(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	// id, ok := r.Context().Value(contextKeyRequestID).(int)
// 	// if !ok {
// 	// 	respondWithJSON(w, http.StatusInternalServerError, resp)
// 	// 	return
// 	// }

// 	course, resp := a.Model.GetCourseCertificateInformation(id)
// 	if len(resp) != 0 {
// 		respondWithJSON(w, http.StatusInternalServerError, resp)
// 		return
// 	}

// 	respondWithJSON(w, http.StatusOK, course)
// }

// swagger:operation GET /v1/certificate/{id}
// @Summary  	 Gets certificate by ID
// @Tags         certificate
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Certificate ID"
// @Success      200  {object}  model.Certificate
// @Router       /v1/certificate/{id} [GET]
func (a *App) GetCertificate(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	postData := PostID{
		ID: id,
	}

	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	cert, resp := a.Model.GetCertificate(postData.ID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, cert)
}

// swagger:operation POST /v1/certificate
// @Summary  	 Creates a certificate by session ID
// @Tags         certificate
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Success      200
// @Router       /v1/certificate [POST]
func (a *App) CreateCertificate(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	fk_course_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	courseID := PostID{
		ID: fk_course_id,
	}

	resp = ValidatePostData(&courseID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	// email := model.PostEmail{
	// 	Email: vars["email"],
	// }

	// resp = ValidatePostData(&email)
	// if len(resp) != 0 {
	// 	respondWithJSON(w, http.StatusBadRequest, resp)
	// 	return
	// }

	var data model.CertificateDTO
	err = DecodeForm(r, &data)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	resp = ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	id, ok := r.Context().Value(contextKeyRequestID).(int)
	if !ok {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	ok, resp = a.Model.AccountHasCourse(id, courseID.ID)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, resp)
		return
	}

	provider, resp := a.Model.GetAccountById(id)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.CreateCertificate(data.AccountEmail, courseID.ID, provider.Name, data.Type, data.Issued, data.Activated, a.Core)

	//TODO: Remove after testing
	// account, resp := a.Model.GetAccountByEmail(email.Email)
	// if len(resp) != 0 {
	// 	respondWithJSON(w, http.StatusInternalServerError, resp)
	// 	return
	// }

	// if account != (model.Account{}) {
	// } else {
	// 	resp = a.Model.CreateTempAccountWithCertificate(data.AccountEmail, id, account.ID, accountID, id, data.Type, data.Issued, data.Activated, a.Core)
	// }

	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "sucess")
}

// swagger:operation PUT /v1/certificate/{id}
// @Summary  	 Updates a certificate by ID
// @Tags         certificate
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Certificate ID"
// @Param        object formData model.CertificateEditDTO  false  "Certificate"
// @Success      200
// @Router       /v1/certificate/{id} [PUT]
func (a *App) UpdateCertificate(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	postData := PostID{
		ID: id,
	}

	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	var data model.CertificateEditDTO
	err = DecodeForm(r, &data)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	resp = ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	//TODO REMOVE AFTER TESTING
	// accountID, ok := r.Context().Value(contextKeyRequestID).(int)
	// if !ok {
	// 	respondWithJSON(w, http.StatusInternalServerError, resp)
	// 	return
	// }
	// if accountID == id {
	// 	// if data.Verified {
	// 	// 	resp["error"] = "You cannot edit a certificate that has been verfied"
	// 	// 	respondWithJSON(w, http.StatusInternalServerError, resp)
	// 	// 	return
	// 	// }

	// 	resp = a.Model.UpdateAccountCertificate(accountID, data.Kin, data.Passport, data.Rank, data.DateOfBirth, data.Licence, data.ShipExpereince, data.Experience, data.Nationality)
	// 	if len(resp) != 0 {
	// 		respondWithJSON(w, http.StatusBadRequest, resp)
	// 		return
	// 	}

	// 	respondWithJSON(w, http.StatusOK, resp)
	// 	return
	// }

	resp = a.Model.UpdateAccountCertificate(id, data.Type, data.Issued, data.Activated)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}
	respondWithJSON(w, http.StatusOK, resp)
}

// swagger:operation DELTE /v1/certificate/{id}
// @Summary  	 Deletes a certificate by ID
// @Tags         certificate
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Certificate ID"
// @Success      200
// @Router       /v1/certificate/{id} [DELETE]
func (a *App) DeleteCertificate(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}
	postData := PostID{
		ID: id,
	}
	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.DeleteCertificate(id)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "success")
}

// swagger:operation GET /v1/certificates/{id}
// @Summary  	 Gets all certificate by account ID
// @Tags         certificate
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Account ID"
// @Success      200  {array}   model.Certificate
// @Router       /v1/certificates/{id} [GET]
func (a *App) GetCertificates(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	id, ok := r.Context().Value(contextKeyRequestID).(int)
	if !ok {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}
	certs, resp := a.Model.GetAccountCertificates(id)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, certs)
}

func (a *App) GetAccountCertificates(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	postData := PostID{
		ID: id,
	}

	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	id, ok := r.Context().Value(contextKeyRequestID).(int)
	if !ok {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}
	certs, resp := a.Model.GetAccountCertificates(id)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, certs)
}

// swagger:operation GET /v1/certificates/course/{id}
// @Summary  	 Gets all certificate by course ID
// @Tags         certificate
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Account ID"
// @Success      200  {array}   model.Certificate
// @Router       /v1/certificates/course/{id} [GET]
func (a *App) GetCourseCertificates(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	postData := PostID{
		ID: id,
	}

	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	// id, ok := r.Context().Value(contextKeyRequestID).(int)
	// if !ok {
	// 	respondWithJSON(w, http.StatusInternalServerError, resp)
	// 	return
	// }

	//TODO: get anylitics
	certs, resp := a.Model.GetCourseCertificatesByCourseID(id)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, certs)
}

// swagger:operation GET /v1/admin/accounts
// @Summary  	 Gets all accounts for Admin
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Success      200  {array}   model.Account
// @Router       /v1/admin/accounts [GET]
func (a *App) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, resp := a.Model.GetAllAccounts()
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, accounts)
}

// swagger:operation POST /v1/admin/account/{id}
// @Summary  	 Sets a user role from the Admin dashboard
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Account ID"
// @Param        object formData model.AccountAdmin  false  "Role"
// @Success      200
// @Router       /v1/admin/account/{id} [POST]
func (a *App) UpdateAdminAccount(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	postData := PostID{
		ID: id,
	}

	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	var data model.AccountAdmin
	err = DecodeForm(r, &data)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	resp = ValidatePostData(&data)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.UpdateAdminAccount(id, data.Status, data.Role)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "success")
}

// swagger:operation GET /v1/admin/courses
// @Summary  	 Deletes a user by ID
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Account ID"
// @Success      200
// @Router       /v1/admin/courses [GET]
func (a *App) DeleteAdminAccount(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	postData := PostID{
		ID: id,
	}

	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	resp = a.Model.DeleteAccount(id)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, "sucess")
}

// swagger:operation GET /v1/admin/courses
// @Summary  	Gets all user course
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        object formData model.Courses  false  "Courses"
// @Success      200
// @Router       /v1/admin/courses [GET]
func (a *App) GetAdminCourse(w http.ResponseWriter, r *http.Request) {
	data, resp := a.Model.GetAllCourses()
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

// swagger:operation GET /v1/admin/courses
// @Summary  	 Gets all courses for Admin
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Success      200  {array}   model.Course
// @Router       /v1/admin/courses [GET]
func (a *App) GetAllCourses(w http.ResponseWriter, r *http.Request) {
	data, resp := a.Model.GetAllCourses()
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

// swagger:operation GET /v1/admin/course/{id}/analytics
// @Summary  	 Gets all courses for Admin analytics
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Param        id query string  true  "Account ID"
// @Success      200  {array}   model.MontlyCourseAnalytics
// @Router       /v1/admin/course/{id}/analytics [GET]
func (a *App) GetAllAccountCoursesAnalytics(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["id"] = "must be a number"
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}

	postData := PostID{
		ID: id,
	}

	resp = ValidatePostData(&postData)
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}
	data, resp := a.Model.GetAccountCourseCertificatesCountInRange(id, time.Now())
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

// swagger:operation GET /v1/admin/certificate
// @Summary  	 Gets all certificate for Admin
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Success      200  {array}   model.Certificate
// @Router       /v1/admin/certificate [GET]
func (a *App) GetAllCertificates(w http.ResponseWriter, r *http.Request) {
	data, resp := a.Model.GetAllCertificates()
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

// swagger:operation GET /v1/lookup/certificate-type
// @Summary  	 Gets certificate types for form dropdown
// @Tags         Lookup table
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Success      200  {array}   model.CertificateTypes
// @Router       /v1/lookup/certificate-type [GET]
func (a *App) GetCertificateTypes(w http.ResponseWriter, r *http.Request) {
	data, resp := a.Model.GetCertificateTypesLookup()
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

// swagger:operation GET  /v1/lookup/account-status
// @Summary  	 Gets account status for form dropdown
// @Tags         Lookup table
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Success      200  {array}   model.AccountStatuses
// @Router       /v1/lookup/account-status [GET]
func (a *App) GetAccountStatuses(w http.ResponseWriter, r *http.Request) {
	data, resp := a.Model.GetAccountStatusLookups()
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

// swagger:operation GET /v1/lookup/account-role
// @Summary  	 Gets account roles for form dropdown
// @Tags         Lookup table
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Success      200  {array}   model.AccountRoles
// @Router       /v1/lookup/account-role [GET]
func (a *App) GetAccountRoles(w http.ResponseWriter, r *http.Request) {
	data, resp := a.Model.GetAccountRoleLookups()
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

// swagger:operation GET /v1/lookup/course-categories
// @Summary  	 Gets course category for form dropdown
// @Tags         Lookup table
// @Accept       json
// @Produce      json
// @Param        Auth  header    string  true  "Authentication header"
// @Success      200  {array}   model.CourseCategories
// @Router       /v1/lookup/course-categories [GET]
func (a *App) GetCourseCategories(w http.ResponseWriter, r *http.Request) {
	data, resp := a.Model.GetCourseCategoryLookups()
	if len(resp) != 0 {
		respondWithJSON(w, http.StatusInternalServerError, resp)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, "not implemented")
}
