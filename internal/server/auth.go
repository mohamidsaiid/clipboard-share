package server

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server)ValidateUser(r *http.Request) bool{
	providedSecretKey := s.getSecretKey(r)	
	hashedSecretKey := s.user.Get().SecretKey

	err := bcrypt.CompareHashAndPassword([]byte(hashedSecretKey), providedSecretKey)

	return err == nil
}

func (s *Server) getSecretKey(r *http.Request) []byte {
	secretKey := r.Header.Get("Authorization")
	secretKeyString, _ := strings.CutPrefix(secretKey, "Bearer ")
	return []byte(secretKeyString)
}