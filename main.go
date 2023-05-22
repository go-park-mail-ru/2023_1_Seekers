package main

import (
	"fmt"
	_mailClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/client"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	size := 1024 * 1024 * 1024
	mailServiceCon, err := grpc.Dial(
		"127.0.0.1:8008",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(size), grpc.MaxCallSendMsgSize(size)),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	mailServiceClient := _mailClient.NewMailClientGRPC(mailServiceCon)
	fmt.Println(mailServiceClient.SendMessage(8, models.FormMessage{
		FromUser:         "geraldy12319@mail.ru",
		Recipients:       []string{"test@mailbx.ru"},
		Title:            "тест",
		Text:             "",
		ReplyToMessageID: nil,
		Attachments:      nil,
	}))
}
