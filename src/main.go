package main

import (
	"go_chat/src/controllers"
	"go_chat/src/database"
	"go_chat/src/repositories"
	"go_chat/src/routes"
	"go_chat/src/services"
	"go_chat/src/ws"
	"log"
	"net/http"
)

func main() {
	db := database.Connect()

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)
	
	chatRepo := repositories.NewChatRepository(db)
	chatService := services.NewChatService(chatRepo)
	chatController := controllers.NewChatController(chatService)
	
	messageRepo := repositories.NewMessageRepository(db)
	messageService := services.NewMessageService(messageRepo)

	hub := ws.NewHub(messageService, chatService)
	go hub.Run()

	router := http.NewServeMux()

	routes.RegisterRoutes(router, userController, chatController, hub)

	http.Handle("/api/", http.StripPrefix("/api", router))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Print("There was an error trying to initialize the server: ", err)
	}
}
