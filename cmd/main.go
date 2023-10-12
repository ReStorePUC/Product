package main

import (
	pbProduct "github.com/ReStorePUC/protobucket/product"
	pb "github.com/ReStorePUC/protobucket/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/restore/product/config"
	"github.com/restore/product/controller"
	"github.com/restore/product/handler"
	"github.com/restore/product/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

func main() {
	config.Init()
	dbCfg := config.NewDBConfig()

	db, err := repository.Init(dbCfg)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("user:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserClient(conn)

	uRepo := repository.NewProduct(db)
	uController := controller.NewProduct(uRepo, c)
	uHandler := handler.NewProduct(uController)

	fHandler := handler.NewFile()

	// GRPC
	go func() {
		lis, err := net.Listen("tcp", ":50053")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pbProduct.RegisterProductServer(s, handler.NewProductServer(uController))
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowFiles:       true,
	}))

	router.GET("/search", uHandler.Search)
	router.GET("/product/store/:id", uHandler.ListProduct)
	router.GET("/product/recent", uHandler.ListRecent)
	router.GET("/product/:id", uHandler.GetProduct)

	router.POST("/files", fHandler.UploadFiles)
	router.Static("/view-file/", "./uploads")
	router.DELETE("/file/:file", fHandler.DeleteFile)

	router.POST("/private/product", uHandler.Create)
	router.POST("/private/unavailable/:id", uHandler.Unavailable)

	router.Run(":8080")
}
