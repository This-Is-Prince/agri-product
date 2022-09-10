package main

import (
	"fmt"

	"github.com/SaiNageswarS/go-api-boot/odm"
	"github.com/SaiNageswarS/go-api-boot/server"
)

var grpcPort = ":50051"
var webPort = ":8081"
var dbName = "agri-product"

func main() {
	fmt.Println("Hello World")
	server.LoadSecretsIntoEnv(false)

	bootServer := server.NewGoApiBoot()

	client := odm.GetClient()
	db := client.Database(dbName)

	productCollection := db.Collection("product")
	fmt.Println(db.Name())
	fmt.Println(productCollection.Name())

	bootServer.Start(grpcPort, webPort)

}
