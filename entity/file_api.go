package entity

import "mime/multipart"

type UploadedFileReq struct {
	File *multipart.FileHeader
}

type UploadedFileRps struct {
	FileId string `json:"fileId"`
	Path   string `json:"path"`
	Size   int    `json:"size"`
}
