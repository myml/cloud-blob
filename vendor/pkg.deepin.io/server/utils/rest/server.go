package rest

// Server is an simple rest api server
type Server struct {
	router *Router
}

// Start Server
func (rs *Server) Start(host, port string) {
	rs.router = StaticRouter()
	rs.router.BootStrap(host + ":" + port)
}
