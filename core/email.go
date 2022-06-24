package core

import (
	"fmt"
	"net/http"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email struct {
	To          string
	ToEmail     string
	Subject     string
	Content     string
	HTMLContent string
}

func (c *Core) SendEmail(e Email) error {
	from := mail.NewEmail("Digital ID", "no-reply@idsure.io")
	to := mail.NewEmail(e.To, e.ToEmail)

	message := mail.NewSingleEmail(from, e.Subject, to, e.Content, e.HTMLContent)

	response, err := c.SendgridClient.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusOK {
		return nil
	}
	return err
}

func (c *Core) ActivationEmail(to string, uid string) error {

	temp := fmt.Sprintf(`<h1>Thank you for signing up</h1>
	<p>To receive and access your new digital certificate in your personal IDsure App you need to <a href='https://idsure.io/signup/%s'>Confirm Here</a></p>
	<p>IDsure ApS, Kranvejen 59, 5000 Odense C. Denmark. CVR: 40223584</p>
	<p>Web: www.IDsure.io Support: support@idsure.io</p> 
	<p>IDsure © Digital ID and Certificates</p>`, uid)

	// filledString := fmt.Sprintf(temp, uid)
	email := Email{
		To:          to,
		ToEmail:     to,
		Subject:     "Action required to access IDsure",
		Content:     "Let's confirm your email address.",
		HTMLContent: temp,
	}

	return c.SendEmail(email)
}

func (c *Core) CertificateEmail(to string) error {
	temp := `<h1>Congratulations on your new achievement!</h1>
	<p>You have received an eCertificate from IDsure</p>
	<p>To view and access your new eCertificate please register <a href='https://idsure.io/login'>Here</a> by downloading the IDsure App and join the new digital platform for your certificates.</p>
	<p>IDsure ApS, Kranvejen 59, 5000 Odense C. Denmark. CVR: 40223584</p>
	<p>Web: www.IDsure.io Support: support@idsure.io</p>
	<p>IDsure © Digital ID and Certificates</p>`

	// filledString := fmt.Sprintf(temp, providerName)
	email := Email{
		To:          to,
		ToEmail:     to,
		Subject:     "You have recieved a certificate from IDsure",
		Content:     "",
		HTMLContent: temp,
	}

	return c.SendEmail(email)
}

func (c *Core) ApprovedEmail(to string) error {
	temp := `<h1>Your account profile has been approved</h1>
	<p>Please log into your IDsure account to access your certificates`
	email := Email{
		To:          to,
		ToEmail:     to,
		Subject:     "Account approved",
		Content:     "",
		HTMLContent: temp,
	}
	return c.SendEmail(email)
}

func (c *Core) RejecctedEmail(to string) error {
	temp := `<h1>Your account profile has been rejected</h1>
	<p>There was an issue with your approval, please login to re submit a approval`
	email := Email{
		To:          to,
		ToEmail:     to,
		Subject:     "Account profile rejected",
		Content:     "",
		HTMLContent: temp,
	}
	return c.SendEmail(email)
}
