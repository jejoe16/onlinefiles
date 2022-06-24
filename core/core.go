package core

import "github.com/sendgrid/sendgrid-go"

type Core struct {
	SendgridClient *sendgrid.Client
}
