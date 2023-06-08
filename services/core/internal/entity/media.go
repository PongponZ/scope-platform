package entity

import (
	"mime/multipart"
)

type MediaFileMeta struct {
	FileName    string `json:"file_name"`
	Buffer      *multipart.File
	ContentType string
	Extension   string
	Type        string
	Size        int64
}

type MediaUploadInfo struct {
	ID                string
	Filename          string
	Path              string
	UserId            string
	Type              string
	OriginalExtension string
}

type MediaConvertPayload struct {
	ID         string `json:"id"`
	FileName   string `json:"filename"`
	Path       string `json:"path"`
	Output     string `json:"output"`
	OutputPath string `json:"outputPath"`
	Type       string `json:"type"`
}

type MediaConvertStatus struct {
	Status  string `json:"status"`
	Process string `json:"process"`
}

type MediaUploaded struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"time_stamp"`
	Endpoint  string `json:"endpoint"`
}
