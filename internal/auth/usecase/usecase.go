package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
)

type useCase struct {
	authRepo auth.RepoI
	userUC   _user.UseCaseI
	mailUC   mail.UseCaseI
	fileUC   file_storage.UseCaseI
}

func New(ar auth.RepoI, uc _user.UseCaseI, mUC mail.UseCaseI, fUC file_storage.UseCaseI) auth.UseCaseI {
	return &useCase{
		authRepo: ar,
		userUC:   uc,
		mailUC:   mUC,
		fileUC:   fUC,
	}
}

func (u *useCase) SignIn(form models.FormLogin) (*models.AuthResponse, *models.Session, error) {
	email, err := pkg.ValidateLogin(form.Login)
	if err != nil {
		return nil, nil, auth.ErrInvalidLogin
	}
	user, err := u.userUC.GetUserByEmail(email)
	if err != nil {
		return nil, nil, auth.ErrWrongPw
	}

	if user.Password != form.Password {
		return nil, nil, auth.ErrWrongPw
	}

	// когда логинимся, то обновляем куку, если ранее была, то удалится и пересоздастся
	err = u.DeleteSessionByUID(user.ID)
	session, err := u.CreateSession(user.ID)
	if err != nil {
		return nil, nil, auth.ErrFailedCreateSession
	}

	f, err := u.fileUC.Get(config.S3AvatarBucket, user.Avatar)
	if err != nil {
		return nil, nil, err
	}

	return &models.AuthResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Image: models.Image{
			Name: f.Name,
			Data: f.Data,
		},
	}, session, nil
}

func (u *useCase) SignUp(form models.FormSignUp) (*models.AuthResponse, *models.Session, error) {
	if form.RepeatPw != form.Password {
		return nil, nil, auth.ErrPwDontMatch
	}

	email, err := pkg.ValidateLogin(form.Login)
	if err != nil || len(form.Login) > 30 || len(form.Login) < 3 {
		return nil, nil, auth.ErrInvalidLogin
	}

	user, err := u.userUC.CreateUser(models.User{
		Email:     email,
		Password:  form.Password,
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Avatar:    config.DefaultAvatar,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("cant create user: %w", err)
	}

	err = u.mailUC.CreateHelloMessage(user.ID)
	if err != nil {
		return nil, nil, auth.ErrInternalHelloMsg
	}

	session, err := u.CreateSession(user.ID)
	if err != nil {
		return nil, nil, auth.ErrFailedCreateSession
	}

	f, err := u.fileUC.Get(config.S3AvatarBucket, user.Avatar)
	if err != nil {
		return nil, nil, err
	}

	return &models.AuthResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Image: models.Image{
			Name: f.Name,
			Data: f.Data,
		},
	}, session, nil
}

func (u *useCase) CreateSession(uID uint64) (*models.Session, error) {
	value, err := pkg.String(config.CookieLen)
	if err != nil {
		return nil, fmt.Errorf("cant create session: %w", err)
	}
	newSession := models.Session{
		UID:       uID,
		SessionID: value,
	}

	err = u.authRepo.CreateSession(newSession)
	if err != nil {
		return nil, fmt.Errorf("cant create session: %w", err)
	}

	return &newSession, nil
}

func (u *useCase) DeleteSession(sessionID string) error {
	err := u.authRepo.DeleteSession(sessionID)
	if err != nil {
		return fmt.Errorf("cant delete session: %w", err)
	}

	return nil
}

func (u *useCase) DeleteSessionByUID(uID uint64) error {
	err := u.authRepo.DeleteSessionByUID(uID)
	if err != nil {
		return fmt.Errorf("cant delete session by id: %w", err)
	}

	return nil
}

func (u *useCase) GetSession(sessionID string) (*models.Session, error) {
	s, err := u.authRepo.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("cant get session: %w", err)
	}

	return s, nil
}

func (u *useCase) GetSessionByUID(uID uint64) (*models.Session, error) {
	s, err := u.authRepo.GetSessionByUID(uID)
	if err != nil {
		return nil, fmt.Errorf("cant get session by user %w", err)
	}

	return s, nil
}

func (u *useCase) GetSessionByEmail(email string) (*models.Session, error) {
	user, err := u.userUC.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("cant get user by email")
	}
	return u.GetSessionByUID(user.ID)
}
