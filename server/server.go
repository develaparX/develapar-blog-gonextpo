package main

import (
	"database/sql"
	"develapar-server/config"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine  *gin.Engine
	portApp string
}

func (s *Server) Start() {
	s.engine.Run(s.portApp)
}

func NewServer() *Server {

	co, _ := config.NewConfig()

	urlConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", co.Host, co.Port, co.User, co.Password, co.Name)

	db, err := sql.Open(co.Driver, urlConnection)
	if err != nil {

		log.Fatal(err)
	}

	portApp := co.AppPort

	return &Server{
		portApp: portApp,
		engine:  gin.Default(),
	}
}
