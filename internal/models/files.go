package models

//go:generate easyjson -all

type Image struct {
	Name string
	Data []byte
}

type S3File struct {
	Bucket string
	Name   string
	Data   []byte
}
