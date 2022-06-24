package model

import (
	"mime/multipart"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Model struct {
	Pool *pgxpool.Pool
}

// swagger:model MultipartFile
type MultipartFile struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}
