// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson3f98fba7DecodeGithubComGoParkMailRu20231SeekersInternalModels(in *jlexer.Lexer, out *S3File) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Bucket":
			out.Bucket = string(in.String())
		case "Name":
			out.Name = string(in.String())
		case "Data":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				out.Data = in.Bytes()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3f98fba7EncodeGithubComGoParkMailRu20231SeekersInternalModels(out *jwriter.Writer, in S3File) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Bucket\":"
		out.RawString(prefix[1:])
		out.String(string(in.Bucket))
	}
	{
		const prefix string = ",\"Name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"Data\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Data)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v S3File) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3f98fba7EncodeGithubComGoParkMailRu20231SeekersInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v S3File) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3f98fba7EncodeGithubComGoParkMailRu20231SeekersInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *S3File) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3f98fba7DecodeGithubComGoParkMailRu20231SeekersInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *S3File) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3f98fba7DecodeGithubComGoParkMailRu20231SeekersInternalModels(l, v)
}
func easyjson3f98fba7DecodeGithubComGoParkMailRu20231SeekersInternalModels1(in *jlexer.Lexer, out *Image) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Name":
			out.Name = string(in.String())
		case "Data":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				out.Data = in.Bytes()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3f98fba7EncodeGithubComGoParkMailRu20231SeekersInternalModels1(out *jwriter.Writer, in Image) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"Data\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Data)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Image) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3f98fba7EncodeGithubComGoParkMailRu20231SeekersInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Image) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3f98fba7EncodeGithubComGoParkMailRu20231SeekersInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Image) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3f98fba7DecodeGithubComGoParkMailRu20231SeekersInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Image) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3f98fba7DecodeGithubComGoParkMailRu20231SeekersInternalModels1(l, v)
}
