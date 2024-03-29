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

func easyjson4a0f95aaDecodeGithubComGoParkMailRu20231SeekersInternalModels(in *jlexer.Lexer, out *EditPasswordRequest) {
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
		case "passwordOld":
			out.PasswordOld = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "repeatPw":
			out.RepeatPw = string(in.String())
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
func easyjson4a0f95aaEncodeGithubComGoParkMailRu20231SeekersInternalModels(out *jwriter.Writer, in EditPasswordRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"passwordOld\":"
		out.RawString(prefix[1:])
		out.String(string(in.PasswordOld))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	{
		const prefix string = ",\"repeatPw\":"
		out.RawString(prefix)
		out.String(string(in.RepeatPw))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EditPasswordRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncodeGithubComGoParkMailRu20231SeekersInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EditPasswordRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncodeGithubComGoParkMailRu20231SeekersInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EditPasswordRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecodeGithubComGoParkMailRu20231SeekersInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EditPasswordRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecodeGithubComGoParkMailRu20231SeekersInternalModels(l, v)
}
func easyjson4a0f95aaDecodeGithubComGoParkMailRu20231SeekersInternalModels1(in *jlexer.Lexer, out *AuthResponse) {
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
		case "email":
			out.Email = string(in.String())
		case "firstName":
			out.FirstName = string(in.String())
		case "lastName":
			out.LastName = string(in.String())
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
func easyjson4a0f95aaEncodeGithubComGoParkMailRu20231SeekersInternalModels1(out *jwriter.Writer, in AuthResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix[1:])
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"firstName\":"
		out.RawString(prefix)
		out.String(string(in.FirstName))
	}
	{
		const prefix string = ",\"lastName\":"
		out.RawString(prefix)
		out.String(string(in.LastName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AuthResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncodeGithubComGoParkMailRu20231SeekersInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AuthResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncodeGithubComGoParkMailRu20231SeekersInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AuthResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecodeGithubComGoParkMailRu20231SeekersInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AuthResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecodeGithubComGoParkMailRu20231SeekersInternalModels1(l, v)
}
