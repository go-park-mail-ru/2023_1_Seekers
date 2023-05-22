package main

//func main() {
//	size := 1024 * 1024 * 1024
//	mailServiceCon, err := grpc.Dial(
//		"89.208.197.150:8008",
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(size), grpc.MaxCallSendMsgSize(size)),
//	)
//	if err != nil {
//		log.Fatal("failed connect to file microservice", err)
//	}
//
//	mailServiceClient := _mailClient.NewMailClientGRPC(mailServiceCon)
//	mailServiceClient.SendMessage()
//}
