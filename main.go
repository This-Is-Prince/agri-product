package main

import (
	"fmt"

	"github.com/SaiNageswarS/go-api-boot/server"
	"github.com/This-Is-Prince/agri-product/pb"
)

var grpcPort = ":50051"
var webPort = ":8081"

func main() {
	fmt.Println("Hello World")
	server.LoadSecretsIntoEnv(false)
	inject := NewInject()

	bootServer := server.NewGoApiBoot()

	pb.RegisterSearchServiceServer(bootServer.GrpcServer, inject.SearchService)
	pb.RegisterListProductServiceServer(bootServer.GrpcServer, inject.ListProductService)
	bootServer.Start(grpcPort, webPort)

}
