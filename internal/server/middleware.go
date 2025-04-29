package server

import "net/http"

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, exists := s.ValidateUser(r)
		if !exists {
			http.Redirect(w, r, "/secretkey", http.StatusSeeOther)
			return
		}

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
