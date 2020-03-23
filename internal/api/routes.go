package api

func (s *Server) routes() {
	s.router.Use(s.logRequest)

	s.router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")
	s.router.HandleFunc("/signin", s.handleSignin()).Methods("POST")
	s.router.HandleFunc("/protected", s.isAuth(s.handleProtect())).Methods("GET")
}
