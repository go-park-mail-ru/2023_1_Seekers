package mail_proto

//go:generate protoc --go_out=./ --go-grpc_out=./ --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional --go_opt=paths=source_relative mail.proto
