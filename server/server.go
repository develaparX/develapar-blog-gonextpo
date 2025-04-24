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
	tS      service.TagService
	atS     service.ArticleTagService
	coS     service.CommentService
	lS      service.LikeService
	engine  *gin.Engine
	portApp string
}

func (s *Server) initiateRoute() {
	routerGroup := s.engine.Group("/api/v1")
	controller.NewUserController(s.uS, routerGroup).Route()
	controller.NewCategoryController(s.cS, routerGroup).Route()
	controller.NewArticleController(s.aS, routerGroup).Route()
	controller.NewBookmarkController(s.bS, routerGroup).Route()
	controller.NewTagController(s.tS, routerGroup).Route()
	controller.NewArticleTagController(s.atS, routerGroup).Route()
	controller.NewCommentController(s.coS, routerGroup).Route()
	controller.NewLikeController(s.lS, routerGroup).Route()
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
	tagRepo := repository.NewTagRepository(db)
	articleTagRepo := repository.NewArticleTagRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)

	userService := service.NewUserservice(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	articleService := service.NewArticleService(articleRepo)
	bookmarkService := service.NewBookmarkService(bookmarkRepo)
	tagService := service.NewTagService(tagRepo)
	articleTagService := service.NewArticleTagService(articleTagRepo)
	commentService := service.NewCommentService(commentRepo)
	likeService := service.NewLikeService(likeRepo)

	return &Server{
		cS:      categoryService,
		uS:      userService,
		aS:      articleService,
		bS:      bookmarkService,
		tS:      tagService,
		atS:     articleTagService,
		coS:     commentService,
		lS:      likeService,
		portApp: portApp,
		engine:  gin.Default(),
	}
}
