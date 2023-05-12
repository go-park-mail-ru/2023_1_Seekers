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

func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels(in *jlexer.Lexer, out *User2Folder) {
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
		case "UserID":
			out.UserID = uint64(in.Uint64())
		case "FolderID":
			out.FolderID = uint64(in.Uint64())
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels(out *jwriter.Writer, in User2Folder) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"UserID\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.UserID))
	}
	{
		const prefix string = ",\"FolderID\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.FolderID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User2Folder) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User2Folder) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User2Folder) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User2Folder) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels1(in *jlexer.Lexer, out *Recipients) {
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
		case "users":
			if in.IsNull() {
				in.Skip()
				out.Users = nil
			} else {
				in.Delim('[')
				if out.Users == nil {
					if !in.IsDelim(']') {
						out.Users = make([]UserInfo, 0, 1)
					} else {
						out.Users = []UserInfo{}
					}
				} else {
					out.Users = (out.Users)[:0]
				}
				for !in.IsDelim(']') {
					var v1 UserInfo
					(v1).UnmarshalEasyJSON(in)
					out.Users = append(out.Users, v1)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels1(out *jwriter.Writer, in Recipients) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"users\":"
		out.RawString(prefix[1:])
		if in.Users == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Users {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Recipients) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Recipients) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Recipients) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Recipients) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels1(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels2(in *jlexer.Lexer, out *MessagesResponse) {
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
		case "messages":
			if in.IsNull() {
				in.Skip()
				out.Messages = nil
			} else {
				in.Delim('[')
				if out.Messages == nil {
					if !in.IsDelim(']') {
						out.Messages = make([]MessageInfo, 0, 0)
					} else {
						out.Messages = []MessageInfo{}
					}
				} else {
					out.Messages = (out.Messages)[:0]
				}
				for !in.IsDelim(']') {
					var v4 MessageInfo
					(v4).UnmarshalEasyJSON(in)
					out.Messages = append(out.Messages, v4)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels2(out *jwriter.Writer, in MessagesResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"messages\":"
		out.RawString(prefix[1:])
		if in.Messages == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Messages {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessagesResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessagesResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessagesResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessagesResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels2(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels3(in *jlexer.Lexer, out *MessageResponse) {
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
		case "message":
			(out.Message).UnmarshalEasyJSON(in)
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels3(out *jwriter.Writer, in MessageResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix[1:])
		(in.Message).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessageResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessageResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessageResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessageResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels3(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels4(in *jlexer.Lexer, out *MessageInfo) {
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
		case "message_id":
			out.MessageID = uint64(in.Uint64())
		case "from_user_id":
			(out.FromUser).UnmarshalEasyJSON(in)
		case "recipients":
			if in.IsNull() {
				in.Skip()
				out.Recipients = nil
			} else {
				in.Delim('[')
				if out.Recipients == nil {
					if !in.IsDelim(']') {
						out.Recipients = make([]UserInfo, 0, 1)
					} else {
						out.Recipients = []UserInfo{}
					}
				} else {
					out.Recipients = (out.Recipients)[:0]
				}
				for !in.IsDelim(']') {
					var v7 UserInfo
					(v7).UnmarshalEasyJSON(in)
					out.Recipients = append(out.Recipients, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "attachments":
			if in.IsNull() {
				in.Skip()
				out.Attachments = nil
			} else {
				in.Delim('[')
				if out.Attachments == nil {
					if !in.IsDelim(']') {
						out.Attachments = make([]AttachmentInfo, 0, 0)
					} else {
						out.Attachments = []AttachmentInfo{}
					}
				} else {
					out.Attachments = (out.Attachments)[:0]
				}
				for !in.IsDelim(']') {
					var v8 AttachmentInfo
					(v8).UnmarshalEasyJSON(in)
					out.Attachments = append(out.Attachments, v8)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "title":
			out.Title = string(in.String())
		case "created_at":
			out.CreatedAt = string(in.String())
		case "text":
			out.Text = string(in.String())
		case "reply_to":
			if in.IsNull() {
				in.Skip()
				out.ReplyTo = nil
			} else {
				if out.ReplyTo == nil {
					out.ReplyTo = new(MessageInfo)
				}
				(*out.ReplyTo).UnmarshalEasyJSON(in)
			}
		case "seen":
			out.Seen = bool(in.Bool())
		case "favorite":
			out.Favorite = bool(in.Bool())
		case "is_draft":
			out.IsDraft = bool(in.Bool())
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels4(out *jwriter.Writer, in MessageInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"message_id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.MessageID))
	}
	{
		const prefix string = ",\"from_user_id\":"
		out.RawString(prefix)
		(in.FromUser).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"recipients\":"
		out.RawString(prefix)
		if in.Recipients == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v9, v10 := range in.Recipients {
				if v9 > 0 {
					out.RawByte(',')
				}
				(v10).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"attachments\":"
		out.RawString(prefix)
		if in.Attachments == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v11, v12 := range in.Attachments {
				if v11 > 0 {
					out.RawByte(',')
				}
				(v12).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.String(string(in.CreatedAt))
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"reply_to\":"
		out.RawString(prefix)
		if in.ReplyTo == nil {
			out.RawString("null")
		} else {
			(*in.ReplyTo).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"seen\":"
		out.RawString(prefix)
		out.Bool(bool(in.Seen))
	}
	{
		const prefix string = ",\"favorite\":"
		out.RawString(prefix)
		out.Bool(bool(in.Favorite))
	}
	{
		const prefix string = ",\"is_draft\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsDraft))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessageInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessageInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessageInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessageInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels4(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels5(in *jlexer.Lexer, out *FormSearchMessages) {
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
		case "fromUser":
			out.FromUser = string(in.String())
		case "toUser":
			out.ToUser = string(in.String())
		case "folder":
			out.Folder = string(in.String())
		case "filter":
			out.Filter = string(in.String())
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels5(out *jwriter.Writer, in FormSearchMessages) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"fromUser\":"
		out.RawString(prefix[1:])
		out.String(string(in.FromUser))
	}
	{
		const prefix string = ",\"toUser\":"
		out.RawString(prefix)
		out.String(string(in.ToUser))
	}
	{
		const prefix string = ",\"folder\":"
		out.RawString(prefix)
		out.String(string(in.Folder))
	}
	{
		const prefix string = ",\"filter\":"
		out.RawString(prefix)
		out.String(string(in.Filter))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FormSearchMessages) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FormSearchMessages) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FormSearchMessages) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FormSearchMessages) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels5(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels6(in *jlexer.Lexer, out *FormMessage) {
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
		case "recipients":
			if in.IsNull() {
				in.Skip()
				out.Recipients = nil
			} else {
				in.Delim('[')
				if out.Recipients == nil {
					if !in.IsDelim(']') {
						out.Recipients = make([]string, 0, 4)
					} else {
						out.Recipients = []string{}
					}
				} else {
					out.Recipients = (out.Recipients)[:0]
				}
				for !in.IsDelim(']') {
					var v13 string
					v13 = string(in.String())
					out.Recipients = append(out.Recipients, v13)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "title":
			out.Title = string(in.String())
		case "text":
			out.Text = string(in.String())
		case "reply_to":
			if in.IsNull() {
				in.Skip()
				out.ReplyToMessageID = nil
			} else {
				if out.ReplyToMessageID == nil {
					out.ReplyToMessageID = new(uint64)
				}
				*out.ReplyToMessageID = uint64(in.Uint64())
			}
		case "attachments":
			if in.IsNull() {
				in.Skip()
				out.Attachments = nil
			} else {
				in.Delim('[')
				if out.Attachments == nil {
					if !in.IsDelim(']') {
						out.Attachments = make([]Attachment, 0, 2)
					} else {
						out.Attachments = []Attachment{}
					}
				} else {
					out.Attachments = (out.Attachments)[:0]
				}
				for !in.IsDelim(']') {
					var v14 Attachment
					(v14).UnmarshalEasyJSON(in)
					out.Attachments = append(out.Attachments, v14)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels6(out *jwriter.Writer, in FormMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"recipients\":"
		out.RawString(prefix[1:])
		if in.Recipients == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v15, v16 := range in.Recipients {
				if v15 > 0 {
					out.RawByte(',')
				}
				out.String(string(v16))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"reply_to\":"
		out.RawString(prefix)
		if in.ReplyToMessageID == nil {
			out.RawString("null")
		} else {
			out.Uint64(uint64(*in.ReplyToMessageID))
		}
	}
	{
		const prefix string = ",\"attachments\":"
		out.RawString(prefix)
		if in.Attachments == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v17, v18 := range in.Attachments {
				if v17 > 0 {
					out.RawByte(',')
				}
				(v18).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FormMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FormMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FormMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FormMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels6(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels7(in *jlexer.Lexer, out *FormFolder) {
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
		case "name":
			out.Name = string(in.String())
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels7(out *jwriter.Writer, in FormFolder) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FormFolder) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FormFolder) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FormFolder) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FormFolder) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels7(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels8(in *jlexer.Lexer, out *FoldersResponse) {
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
		case "folders":
			if in.IsNull() {
				in.Skip()
				out.Folders = nil
			} else {
				in.Delim('[')
				if out.Folders == nil {
					if !in.IsDelim(']') {
						out.Folders = make([]Folder, 0, 1)
					} else {
						out.Folders = []Folder{}
					}
				} else {
					out.Folders = (out.Folders)[:0]
				}
				for !in.IsDelim(']') {
					var v19 Folder
					(v19).UnmarshalEasyJSON(in)
					out.Folders = append(out.Folders, v19)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "count":
			out.Count = int(in.Int())
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels8(out *jwriter.Writer, in FoldersResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"folders\":"
		out.RawString(prefix[1:])
		if in.Folders == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v20, v21 := range in.Folders {
				if v20 > 0 {
					out.RawByte(',')
				}
				(v21).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"count\":"
		out.RawString(prefix)
		out.Int(int(in.Count))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FoldersResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FoldersResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FoldersResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FoldersResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels8(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels9(in *jlexer.Lexer, out *FolderResponse) {
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
		case "folder":
			(out.Folder).UnmarshalEasyJSON(in)
		case "messages":
			if in.IsNull() {
				in.Skip()
				out.Messages = nil
			} else {
				in.Delim('[')
				if out.Messages == nil {
					if !in.IsDelim(']') {
						out.Messages = make([]MessageInfo, 0, 0)
					} else {
						out.Messages = []MessageInfo{}
					}
				} else {
					out.Messages = (out.Messages)[:0]
				}
				for !in.IsDelim(']') {
					var v22 MessageInfo
					(v22).UnmarshalEasyJSON(in)
					out.Messages = append(out.Messages, v22)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels9(out *jwriter.Writer, in FolderResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"folder\":"
		out.RawString(prefix[1:])
		(in.Folder).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"messages\":"
		out.RawString(prefix)
		if in.Messages == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v23, v24 := range in.Messages {
				if v23 > 0 {
					out.RawByte(',')
				}
				(v24).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FolderResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FolderResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FolderResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FolderResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels9(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels10(in *jlexer.Lexer, out *Folder) {
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
		case "folder_id":
			out.FolderID = uint64(in.Uint64())
		case "folder_slug":
			out.LocalName = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "messages_unseen":
			out.MessagesUnseen = int(in.Int())
		case "messages_count":
			out.MessagesCount = int(in.Int())
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels10(out *jwriter.Writer, in Folder) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"folder_id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.FolderID))
	}
	{
		const prefix string = ",\"folder_slug\":"
		out.RawString(prefix)
		out.String(string(in.LocalName))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"messages_unseen\":"
		out.RawString(prefix)
		out.Int(int(in.MessagesUnseen))
	}
	{
		const prefix string = ",\"messages_count\":"
		out.RawString(prefix)
		out.Int(int(in.MessagesCount))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Folder) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Folder) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Folder) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Folder) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels10(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels11(in *jlexer.Lexer, out *AttachmentInfo) {
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
		case "attachID":
			out.AttachID = uint64(in.Uint64())
		case "fileName":
			out.FileName = string(in.String())
		case "type":
			out.Type = string(in.String())
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels11(out *jwriter.Writer, in AttachmentInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"attachID\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.AttachID))
	}
	{
		const prefix string = ",\"fileName\":"
		out.RawString(prefix)
		out.String(string(in.FileName))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AttachmentInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels11(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AttachmentInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels11(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AttachmentInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels11(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AttachmentInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels11(l, v)
}
func easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels12(in *jlexer.Lexer, out *Attachment) {
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
		case "fileName":
			out.FileName = string(in.String())
		case "fileData":
			out.FileData = string(in.String())
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
func easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels12(out *jwriter.Writer, in Attachment) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"fileName\":"
		out.RawString(prefix[1:])
		out.String(string(in.FileName))
	}
	{
		const prefix string = ",\"fileData\":"
		out.RawString(prefix)
		out.String(string(in.FileData))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Attachment) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels12(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Attachment) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7feca409EncodeGithubComGoParkMailRu20231SeekersInternalModels12(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Attachment) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels12(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Attachment) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7feca409DecodeGithubComGoParkMailRu20231SeekersInternalModels12(l, v)
}
