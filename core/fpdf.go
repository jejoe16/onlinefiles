package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

func (c *Core) GeneratePDF(w http.ResponseWriter, headerImage string, profileImage string, courseName string, certificateName string, userAccountName *string, courseCategory string, created time.Time, expire int, certificateUID string, id int, passportNumber *string, dateOfBirth *time.Time) error {

	var formatRect = func(pdf *fpdf.Fpdf, font string, style string, size float64, height float64, text string, align string) {
		pdf.SetAutoPageBreak(false, 0)
		pdf.SetXY(20, 20)
		pdf.SetFont(font, style, size)
		pdf.CellFormat(170, height, text, "", 0, align, false, 0, "")
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	headerImageStream, err := http.Get(fmt.Sprintf("https://did2.s3.us-west-2.amazonaws.com/%s", headerImage))
	if err != nil {
		return err
	}

	defer headerImageStream.Body.Close()
	if headerImageStream.StatusCode != 200 {
		return err
	}

	buf1 := &bytes.Buffer{}

	tee := io.TeeReader(headerImageStream.Body, buf1)

	headerBody, err := ioutil.ReadAll(tee)
	if err != nil {
		return err
	}

	headerType := strings.TrimPrefix(http.DetectContentType(headerBody), "image/")

	pdf.RegisterImageOptionsReader("header_image", fpdf.ImageOptions{ImageType: headerType}, buf1)
	pdf.Image("header_image", 50, 0, 100, 70, false, "", 0, "")
	pdf.Ln(50)

	title := fmt.Sprintf("%s %s", courseName, certificateName)
	pdf.SetFont("Arial", "", 25)
	pdf.CellFormat(200, 10, title, "", 1, fpdf.AlignCenter, false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 10, fmt.Sprintf("This is to certify that: %s", *userAccountName), "", 0, fpdf.AlignBaseline, false, 0, "")
	pdf.Ln(10)
	pdf.CellFormat(0, 10, fmt.Sprintf("Has been found duly qualified to be in charge of: %s", courseCategory), "", 0, fpdf.AlignBaseline, false, 0, "")
	profileImageStream, err := http.Get(fmt.Sprintf("https://did2.s3.us-west-2.amazonaws.com/%s", profileImage))
	if err != nil {
		return err
	}

	defer profileImageStream.Body.Close()
	if headerImageStream.StatusCode != 200 {
		return err
	}
	buf2 := &bytes.Buffer{}
	tee2 := io.TeeReader(profileImageStream.Body, buf2)

	profileBody, err := ioutil.ReadAll(tee2)
	if err != nil {
		return err
	}

	profileType := strings.TrimPrefix(http.DetectContentType(profileBody), "image/")
	if profileType == "jpeg" {
		profileType = "jpg"
	}

	pdf.RegisterImageOptionsReader("profile_image", fpdf.ImageOptions{ImageType: profileType}, buf2)
	pdf.Image("profile_image", 175, 60, 25, 25, false, "", 0, "")

	formatRect(pdf, "Arial", "B", 18, 45, courseName, "C")
	formatRect(pdf, "Arial", "B", 11, 128, fmt.Sprintf("This is to certify that: %s", *userAccountName), "L")
	formatRect(pdf, "Arial", "B", 11, 143, "Has been found buly qualified to be in charge of:", "L")
	formatRect(pdf, "Arial", "", 12, 143, GetStringWithLeftSpace(courseCategory, 78), "L")

	formatRect(pdf, "Arial", "B", 11, 200, "Issued on: ", "L")
	formatRect(pdf, "Arial", "", 12, 200, GetStringWithLeftSpace(time.Time(created).String(), 17), "L")
	formatRect(pdf, "Arial", "B", 11, 200, GetStringWithRightSpace("Valid untill:", 20), "R")
	formatRect(pdf, "Arial", "", 11, 200, GetStringWithLeftSpace(created.AddDate(0, expire, 0).String(), 20), "R")

	formatRect(pdf, "Arial", "B", 11, 219, "Certificate no. ", "L")
	formatRect(pdf, "Arial", "", 12, 219, GetStringWithLeftSpace(certificateUID, 23), "L")
	formatRect(pdf, "Arial", "B", 11, 219, GetStringWithRightSpace("Course ID:", 44), "R")
	formatRect(pdf, "Arial", "", 11, 219, GetStringWithLeftSpace(strconv.FormatInt(int64(id), 10), 44), "R")

	formatRect(pdf, "Arial", "B", 11, 238, "Passport number of the certificate holder ", "L")
	formatRect(pdf, "Arial", "", 12, 238, GetStringWithLeftSpace(*passportNumber, 66), "L")
	if dateOfBirth != nil {
		formatRect(pdf, "Arial", "B", 11, 238, GetStringWithRightSpace("Date of birth:", 20), "R")
		formatRect(pdf, "Arial", "", 11, 238, GetStringWithLeftSpace(dateOfBirth.String(), 20), "R")
	}

	formatRect(pdf, "Arial", "", 12, 335, "The issuers organisation is approved by XYZ", "C")
	formatRect(pdf, "Arial", "B", 16, 353, "Certificate QR:", "C")

	return pdf.Output(w)
}

func GetStringWithLeftSpace(str string, space int) string {
	return (strings.Repeat(" ", space) + str)
}
func GetStringWithRightSpace(str string, space int) string {
	return (str + strings.Repeat(" ", space))
}
