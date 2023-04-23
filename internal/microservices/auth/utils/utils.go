package utils

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
)

func ProtoBySessionModel(s *models.Session) *auth_proto.Session {
	return &auth_proto.Session{
		UID:   s.UID,
		Value: s.SessionID,
	}
}

func SessionModelByProto(protoSession *auth_proto.Session) *models.Session {
	return &models.Session{SessionID: protoSession.Value, UID: protoSession.UID}
}

func LoginFormModelByProto(protoForm *auth_proto.FormLogin) *models.FormLogin {
	return &models.FormLogin{
		Login:    protoForm.Login,
		Password: protoForm.Password,
	}
}

func ProtoByLoginFormModel(form *models.FormLogin) *auth_proto.FormLogin {
	return &auth_proto.FormLogin{
		Login:    form.Login,
		Password: form.Password,
	}
}

func SignupFormModelByProto(protoForm *auth_proto.FormSignup) *models.FormSignUp {
	return &models.FormSignUp{
		Login:     protoForm.Login,
		Password:  protoForm.Password,
		RepeatPw:  protoForm.RepeatPassword,
		FirstName: protoForm.FirstName,
		LastName:  protoForm.LastName,
	}
}

func ProtoBySignupFormModel(form *models.FormSignUp) *auth_proto.FormSignup {
	return &auth_proto.FormSignup{
		Login:          form.Login,
		Password:       form.Password,
		RepeatPassword: form.RepeatPw,
		FirstName:      form.FirstName,
		LastName:       form.LastName,
	}
}

func ProtoByAuthResponseModel(info *models.AuthResponse) *auth_proto.UserInfo {
	return &auth_proto.UserInfo{
		Email:     info.Email,
		FirstName: info.FirstName,
		LastName:  info.LastName,
	}
}

func AuthResponseModelByProto(protoResponse *auth_proto.AuthResponse) *models.AuthResponse {
	return &models.AuthResponse{
		Email:     protoResponse.Info.Email,
		FirstName: protoResponse.Info.FirstName,
		LastName:  protoResponse.Info.LastName,
	}
}

func SessionModelByProtoAuthResponse(protoResponse *auth_proto.AuthResponse) *models.Session {
	return &models.Session{
		UID:       protoResponse.Session.UID,
		SessionID: protoResponse.Session.Value,
	}
}

func ProtoAuthResponseByInfoNSession(info *models.AuthResponse, session *models.Session) *auth_proto.AuthResponse {
	return &auth_proto.AuthResponse{
		Info:    ProtoByAuthResponseModel(info),
		Session: ProtoBySessionModel(session),
	}
}
