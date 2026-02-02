package routes

import (
	"go_chat/src/controllers"
	"go_chat/src/middlewares"
	"go_chat/src/ws"
	"net/http"
)

func RegisterRoutes(r *http.ServeMux, uc *controllers.UserController, cc *controllers.ChatController, mc *controllers.MessageController, hub *ws.Hub) {
	r.HandleFunc("GET /user/", middlewares.Protect(uc.SearchUsers))
	r.HandleFunc("POST /user/", uc.Register)
	r.HandleFunc("POST /user/login", uc.AuthUser)
	r.HandleFunc("GET /user/info", uc.Info)

	r.HandleFunc("POST /chat/", middlewares.Protect(cc.AccessChat))
	r.HandleFunc("GET /chat/", middlewares.Protect(cc.FetchChats))

	r.HandleFunc("GET /messages/", middlewares.Protect(mc.Fetch))

	r.HandleFunc("GET /ws/", middlewares.Protect(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	}))
}
