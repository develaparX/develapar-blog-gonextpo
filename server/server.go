package main

import (
	"database/sql"
	"develapar-server/config"
	"develapar-server/controller"
	"develapar-server/repository"
	"develapar-server/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	uS      service.UserService
	cS      service.CategoryService
	aS      service.ArticleService
	bS      service.BookmarkService
	engine  *gin.Engine
	portApp string
}

func (s *Server) initiateRoute() {
	routerGroup := s.engine.Group("/api/v1")
	controller.NewUserController(s.uS, routerGroup).Route()
	controller.NewCategoryController(s.cS, routerGroup).Route()
	controller.NewArticleController(s.aS, routerGroup).Route()
	controller.NewBookmarkController(s.bS, routerGroup).Route()
}

func (s *Server) Start() {
	s.initiateRoute()
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

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)

	userService := service.NewUserservice(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	articleService := service.NewArticleService(articleRepo)
	bookmarkService := service.NewBookmarkService(bookmarkRepo)

	return &Server{
		cS:      categoryService,
		uS:      userService,
		aS:      articleService,
		bS:      bookmarkService,
		portApp: portApp,
		engine:  gin.Default(),
	}
}
