package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gogrpccrud/internal/repositories"
	"gogrpccrud/internal/storage"
	"gogrpccrud/pkg/api"
)

func main() {
	server := grpc.NewServer()

	db, err := storage.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userService := repositories.NewUserServise(db)

	api.RegisterUserServer(server, userService)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	reflection.Register(server)

	if err := server.Serve(listener); err != io.EOF {
		log.Fatal(err)
	}
}
