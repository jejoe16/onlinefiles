## Understanding the URL params
### JSON
This is the actual name of the form field you must request

### Validate
These are the responses you will get back from the POST data if not correct

# REST API
All data is passed over a form post

## Sign In
Sets the session with the account user

### Request
    POST /v1/signin

### URL Params
	Email           string `json:"email"            validate:"required,gte=4,lte=320,email"`
	Password        string `json:"password"         validate:"required,gte=8,lte=156"`
    Remember        bool `json:"remember" 
### Success Response
    status: 200
    remember-token=2067ea1a-0bea-4164-b2c0-c5d250b56911; Expires=Thu, 23 Dec 2021 08:13:24 GMT; Secure      
    session-token=MTYzNzY0Nzg0NHxpb0gzUFVzc01hYWdmTWhmeXNDT2VMcE02ZEJob3NxUG1TOUNXbFoxc2VqS3VNQU1DQUlWWW9PX3QydXVUckV6Y0lQY1REeDdVdW9RfAEspDk1ZNO2PHhpC7uZJoZf8DgcDCc4FS7vmBhuqUXQ; Expires=Tue, 23 Nov 2021 06:25:44 GMT; Max-Age=900; Secure
---

## Sign Up
Creates a new Account user

### Request
    POST /v1/signup

### URL Params
	Email           string `json:"email"            validate:"required,gte=4,lte=320,email"`
	Password        string `json:"password"         validate:"required,gte=8,lte=156"`
	PasswordConfirm string `json:"passwordconfirm"  validate:"required,gte=8,lte=156"`


### Success Response
    status: 200
    {system:"success"}
---

## Validate Email
Once a user is signed up then they are sent an email, they need to validate inorder to use the account. 

### Request
    GET /v1/validate/email/{uid}

### URL Params
    UID uid `json:"uid" validate:"required,uid"`

### Success Response
    status: 200
---

## Account information
Returns the account info for keeping data

### Request
    GET /v1/account

### URL Params

### Success Response
    status: 200
	ID            int     
	Email         string  
	Name          *string 
	Status        *int      [Active, Banned]
	Role          *int      [Unverified, User, Certifier, Admin]
	EmailValidate bool    

    Unverified is if the email has not been verified, User means they can only SEE certificates, Certifier means they can make courses and certificates, Admin is the site wide admin.
---


## Get Account
Gets the account by id

### Request
    GET /v1/account/{id}

### Header
session cookie

### URL Params
    id              int `json:"id" validate:"required,numeric"`

### Success Response
    status: 200
	ID            int     `json:"id" validate:"required,numeric"`
	Email         string  `json:"email" validate:"required,alphanum"`
	Name          *string `json:"name" validate:"required,lte=100"`
	Status        *int    `json:"status" validate:"gte=0,required,numeric"`
	Role          *int    `json:"role" validate:"gte=0,required,numeric"`
	EmailValidate bool    `json:"email_validated" validate:"required,boolean"`      
---

## Delete account
Deletes account by ud

### Request
    DELETE /v1/account/{id}

### Header
session cookie

### URL Params
    id              int `json:"id" validate:"required,numeric"`

### Success Response
    status: 200
---

## Update account
Updates a account by id

### Request
    PUT /v1/account/{id}

### Header
session cookie

### URL Params
    id              int `json:"id" validate:"required,numeric"`

### Success Response
    status: 200
---


## Get Course
Gets the account course by id

### Request
    GET /v1/course/{id}

### Header
session cookie

### URL Params
    id              int `json:"id" validate:"required,numeric"`

### Success Response
    status: 200
	ID                    int            `json:"id"`
	Image                 sql.NullString `json:"image"`
	Description           sql.NullString `json:"description" validate:"required,gte=3,lte=500"`
	CourseName            string         `json:"course_name" validate:"required,gte=3,lte=100"`
	CertificateName       string         `json:"certificate_name" validate:"required,gte=3,lte=100"`
	AdditionalDescription sql.NullString `json:"additional_description" validate:"gte=3,lte=500"`
	Expire                int            `json:"expire" validate:"required,gte=0,numeric"`
	CourseID              int            `json:"course_id" validate:"required,gte=0,numeric"`         
---

## Create Course
Create a new account course (must be a certifier)

### Request
    POST /v1/course

### Header
session cookie

### URL Params
	ID                    int            `json:"id"`
	Image                 sql.NullString `json:"image"`
	Description           sql.NullString `json:"description" validate:"required,gte=3,lte=500"`
	CourseName            string         `json:"course_name" validate:"required,gte=3,lte=100"`
	CertificateName       string         `json:"certificate_name" validate:"required,gte=3,lte=100"`
	AdditionalDescription sql.NullString `json:"additional_description" validate:"gte=3,lte=500"`
	Expire                int            `json:"expire" validate:"required,gte=0,numeric"`
	CourseID              int            `json:"course_id" validate:"required,gte=0,numeric"`

### Success Response
    status: 200
    {system:"success"}
---

## Get Courses
Gets the account courses

### Request
    GET /v1/courses

### Header
session cookie

### URL Params

### Success Response
    status: 200
	ID                    int            `json:"id"`
	Image                 sql.NullString `json:"image"`
	Description           sql.NullString `json:"description" validate:"required,gte=3,lte=500"`
	CourseName            string         `json:"course_name" validate:"required,gte=3,lte=100"`
	CertificateName       string         `json:"certificate_name" validate:"required,gte=3,lte=100"`
	AdditionalDescription sql.NullString `json:"additional_description" validate:"gte=3,lte=500"`
	Expire                int            `json:"expire" validate:"required,gte=0,numeric"`
	CourseID              int            `json:"course_id" validate:"required,gte=0,numeric"`

    Courses[]
---

## Update Course
Updates a course by the account ID

### Request 
    PUT /v1/course/{id}

### Header
session cookie

### URL Params
	ID                    int    `json:"id"`
	Name                  string `json:"name"                   validate:"gte=3,lte=100,required"`
	Description           string `json:"description"            validate:"gte=3,lte=500,required"`
	CourseName            string `json:"course_name"            validate:"gte=3,lte=100,required"`
	AdditionalDescription string `json:"additional_description" validate:"gte=3,lte=500"`
	Expire                int    `json:"expire"                 validate:"gte=0,required,numeric"`
	Number                int    `json:"number"                 validate:"gte=0,required,numeric"`
	CourseID              int    `json:"course_id"              validate:"gte=0,required,numeric"`

### Success Response
    status: 200
    {system:"success"}
---

## Delete Course
Delete course by id

### Request 
    PUT /v1/course/{id}

### Header
session cookie

### URL Params
	ID          int    `json:"id" validate:"required,numeric"`

### Success Response
    status: 200
    {system:"success"}
---

## Create Certificate
Create a course certificate

### Request
    POST /v1/certificate/{id}

### Header
session cookie

### URL Params
	CourseID       int       `json:"fk_course_id" validate:"required,gte=0,numeric"`
	AccountEmail   string    `json:"account_email" validate:"required,email"`
	Type           int       `json:"type" validate:"gte=0,required,numeric"`
	Issued         time.Time `json:"issued" validate:"required,datetime"`
	Activated      bool      `json:"activated" validate:"required,boolean"`
	Phone          int       `json:"phone"`
	Kin            string    `json:"kin"`
	Passport       string    `json:"passport"`
	Rank           string    `json:"rank"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	Licence        string    `json:"licence"`
	ShipExpereince string    `json:"ship_experience"`
	Experience     int       `json:"experience"`
	Nationality    string    `json:"nationality"`
	Verified       bool      `json:"verified"`

### Success Response
    status: 200
    {system:"success"}
---

## Get course certificates
Returns all the certificates assigned from a course

### Request
    GET /v1/certificates/course/{id}

### Header
session cookie

### URL Params
	ID        int       `json:"id"`

### Success Response
    status: 200
	ID        int       `json:"id"`
	Course    int       `json:"fk_course_id" validate:"gte=0,required,numeric"`
	Type      int       `json:"type" validate:"gte=0,required,numeric"`
	Issued    time.Time `json:"issued" validate:"required,datetime"`
	Activated bool      `json:"activated" validate:"required,boolean"`
	Accessed  int       `json:"acessed" validate:"gte=0,required,numeric"`

    certificates[]
---
<!-- 
## Admin account get
Gets a admin account by ID

### Request
    GET /v1/admin/account/{id}

### Header
session cookie

### URL Params
    ID            int   `json:"id" validate:"required,numeric"`

### Success Response
    status: 200
	ID            int    `json:"id" validate:"required,numeric"`
	Email         string `json:"email" validate:"required,numeric"`
	Status        *int   `json:"status" validate:"gte=0,required,numeric"`
	Role          *int   `json:"role" validate:"gte=0,required,numeric"`
	EmailValidate bool   `json:"email_validated" validate:"required,boolean"`
--- -->

## Admin Get all accounts
Lists all of the site accounts

### Request
    GET /v1/admin/accounts

### Header
Must have an account role of 4 to access page
- session cookie

### URL Params

### Success Response
    status: 200
	ID            int    `json:"id" validate:"required,numeric"`
	Email         string `json:"email" validate:"required,numeric"`
	Status        *int   `json:"status" validate:"gte=0,required,numeric"`
	Role          *int   `json:"role" validate:"gte=0,required,numeric"`
	EmailValidate bool   `json:"email_validated" validate:"required,boolean"`

    Accounts[]
---

## Admin Update account
Updates an account with status and role

### Request
    POST /v1/admin/account/{id}

### URL Params
	ID            int    `json:"id" validate:"required,numeric"`
	Status        *int   `json:"status" validate:"gte=0,required,numeric"`    [Active, Banned]
	Role          *int   `json:"role" validate:"gte=0,required,numeric"`      [Unverified, User, Certifier, Admin]
### Success Response
    status: 200
    {system:"success"}
---

## Admin Get all courses
Lists all of the site courses

### Request
    GET /v1/admin/courses

### Header
Must have an account role of 4 to access page
- session cookie

### URL Params

### Success Response
    status: 200
	ID                    int    `json:"id"						validate:"required,gte=0,numeric"`
	Name                  string `json:"name" 					validate:"required,gte=3,lte=100"`
	Image                 string `json:"image"`
	Description           string `json:"description" 			validate:"required,gte=3,lte=500"`
	CourseName            string `json:"course_name" 			validate:"required,gte=3,lte=100"`
	AdditionalDescription string `json:"additional_description" validate:"gte=3,lte=500"`
	Expire                int    `json:"expire" 				validate:"required,gte=0,numeric"`
	CourseID              int    `json:"course_id" 				validate:"required,gte=0,numeric"`

    Courses[]
---

## Admin Update course
Updates a course

### Request
    POST /v1/admin/course/{id}

### URL Params
	ID            int    `json:"id" validate:"required,numeric"`
    Name                  string `json:"name" 					validate:"required,gte=3,lte=100"`
	Description           string `json:"description" 			validate:"required,gte=3,lte=500"`
	CourseName            string `json:"course_name" 			validate:"required,gte=3,lte=100"`
	AdditionalDescription string `json:"additional_description" validate:"gte=3,lte=500"`
	Expire                int    `json:"expire" 				validate:"required,gte=0,numeric"`
	CourseID              int    `json:"course_id" 				validate:"required,gte=0,numeric"`

### Success Response
    status: 200
    {system:"success"}

---


## Admin Get all certificates
Lists all of the site certificates

### Request
    GET /v1/admin/certificates

### Header
Must have an account role of 4 to access page
- session cookie

### URL Params

### Success Response
    status: 200
	ID             int       `json:"id" validate:"required,gte=0,numeric"`
	CourseID       int       `json:"fk_course_id" validate:"required,gte=0,numeric"`
	AccountID      int       `json:"fk_account_id" validate:"required,gte=0,numeric"`
	AccountEmail   string    `json:"account_email" validate:"required,email"`
	Type           int       `json:"type" validate:"required,gte=0,numeric"`
	Issued         time.Time `json:"issued" validate:"required,datetime"`
	Activated      bool      `json:"activated" validate:"required,boolean"`
	Accessed       int       `json:"accessed" validate:"required,gte=0,numeric"`
	Phone          int       `json:"phone" validate:"required,gte=0,numeric"`
	Kin            string    `json:"kin" validate:""`
	Passport       string    `json:"passport"  validate:"required"`
	Rank           string    `json:"rank" validate:"required"`
	DateOfBirth    time.Time `json:"date_of_birth" validate:"required"`
	Licence        string    `json:"licence" validate:"required"`
	ShipExpereince string    `json:"ship_experience" validate:"required"`
	Experience     int       `json:"experience" validate:"required"`
	Nationality    string    `json:"nationality" validate:"required"`
	Completed      bool      `json:"completed" validate:"required,bool"`
	Verified       bool      `json:"verified" validate:"required,bool"`

    Certificates[]
---

## Admin Update certificate
Updates a certificate

### Request
    POST /v1/admin/course/{id}

### URL Params
	ID             int       `json:"id" validate:"required,gte=0,numeric"`
	CourseID       int       `json:"fk_course_id" validate:"required,gte=0,numeric"`
	AccountID      int       `json:"fk_account_id" validate:"required,gte=0,numeric"`
	AccountEmail   string    `json:"account_email" validate:"required,email"`
	Type           int       `json:"type" validate:"required,gte=0,numeric"`
	Issued         time.Time `json:"issued" validate:"required,datetime"`
	Activated      bool      `json:"activated" validate:"required,boolean"`
	Accessed       int       `json:"accessed" validate:"required,gte=0,numeric"`
	Phone          int       `json:"phone" validate:"required,gte=0,numeric"`
	Kin            string    `json:"kin" validate:""`
	Passport       string    `json:"passport"  validate:"required"`
	Rank           string    `json:"rank" validate:"required"`
	DateOfBirth    time.Time `json:"date_of_birth" validate:"required"`
	Licence        string    `json:"licence" validate:"required"`
	ShipExpereince string    `json:"ship_experience" validate:"required"`
	Experience     int       `json:"experience" validate:"required"`
	Nationality    string    `json:"nationality" validate:"required"`
	Completed      bool      `json:"completed" validate:"required,bool"`
	Verified       bool      `json:"verified" validate:"required,bool"`

### Success Response
    status: 200
    {system:"success"}

---

## Lookup table for certificate types
Gets the list of records for the certificate dropdowns

### Request
   GET v1/lookup/certificate-type

### URL Params

### Success Response
    status: 200
	ID   int    `json:"id"`
	Name string `json:"name"`
---

## Lookup table for account status
Gets the list of records for the account status dropdowns

### Request
   GET /v1/lookup/account-status

### URL Params

### Success Response
    status: 200
	ID   int    `json:"id"`
	Name string `json:"name"`

---

## Lookup table for account roles
Gets the list of records for the account roles

### Request
   GET /v1/lookup/account-role

### URL Params

### Success Response
    status: 200
	ID   int    `json:"id"`
	Name string `json:"name"`

---
