package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/smtp/client"
)

func main() {
	err := client.SendMail(&models.User{
		Email:      "test@mailbx.ru",
		Password:   "12345",
		FirstName:  "Max",
		LastName:   "Vlasov",
		Avatar:     "default_avatar",
		IsExternal: false,
	}, "geraldy12319@yandex.ru", "test", "body", "mailbx.ru")
	fmt.Println(err)
}
