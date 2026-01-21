package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)


func main(){
	// иницилизация AuthService
	AuthSvc
	MassageSvc
	// иницилизация ручек (отправка сообщений пользователю и всем пользователям)
	authH := handlers.NewAuthHandler(authSvc)
	// роутер
	router := gin.Default()
	router.POST("/auth/login",authH.Login)
	router.POST("/send/:id",)
	router.POST("/send/broadcast",)

	// Protected endpoints
	// protected := router.Group("/")

	// start http server
	port := os.Getenv("PORT")
	if port == "" {
        port = "8080"
    }
	addr := fmt.Sprintf(":%s", port)
    log.Printf("starting server on %s", addr)
    if err := http.ListenAndServe(addr, router); err != nil {
        log.Fatal("server failed:", err)
    }
// 	  router := gin.Default()
//   router.GET("/ping", func(c *gin.Context) {
//     c.JSON(200, gin.H{
//       "message": "pong",
//     })
//   })
//   router.Run() // listens on 0.0.0.0:8080 by default
}