package controllers

func (server *Server) InitializeRoutes() {

	//Signup API routes
	server.Router.POST("/api/v1/signup", server.signUp)

}
