package models

type Image struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type S3File struct {
	Bucket string
	Name   string
	Data   []byte
}
