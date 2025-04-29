package server

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server)ValidateUser(r *http.Request) (valid bool, exists bool){
	providedSecretKey := s.getSecretKey(r)	
	user, ok := s.user.Get()
	if !ok {
		exists = false
		return
	}
	hashedSecretKey := user.SecretKey
	err := bcrypt.CompareHashAndPassword([]byte(hashedSecretKey), providedSecretKey)

	return err == nil, true
}

func (s *Server) getSecretKey(r *http.Request) []byte {
	secretKey := r.Header.Get("Authorization")
	secretKeyString, _ := strings.CutPrefix(secretKey, "Bearer ")
	return []byte(secretKeyString)
}