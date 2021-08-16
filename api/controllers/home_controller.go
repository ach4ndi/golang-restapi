package controllers

import (
	"net/http"

	"../responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Selamat Datang di API test ini")

}
