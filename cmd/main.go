package main

import (
	pb "github.com/ReStorePUC/protobucket/generated"
	"github.com/gin-gonic/gin"
	"github.com/restore/product/config"
	"github.com/restore/product/controller"
	"github.com/restore/product/handler"
	"github.com/restore/product/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	config.Init()
	dbCfg := config.NewDBConfig()

	db, err := repository.Init(dbCfg)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserClient(conn)

	uRepo := repository.NewProduct(db)
	uController := controller.NewProduct(uRepo, c)
	uHandler := handler.NewProduct(uController)

	fHandler := handler.NewFile()

	router := gin.Default()
	router.GET("/search", uHandler.Search)
	router.GET("/product/store/:id", uHandler.ListProduct)
	router.GET("/product/recent", uHandler.ListRecent)
	router.GET("/product/:id", uHandler.GetProduct)

	router.POST("/files", fHandler.UploadFiles)
	router.GET("/file/:file", fHandler.GetFile)
	router.DELETE("/file/:file", fHandler.DeleteFile)

	router.POST("/private/product", uHandler.Create)
	router.POST("/private/unavailable/:id", uHandler.Unavailable)

	router.Run(":8080")
}
