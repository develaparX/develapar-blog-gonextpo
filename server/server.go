package main

import (
	"database/sql"
	"develapar-server/config"
	"develapar-server/controller"
	"develapar-server/middleware"
	"develapar-server/repository"
	"develapar-server/service"
	"fmt"
	"log"
	"net/http"

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
	jS      service.JwtService
	mD middleware.AuthMiddleware
	engine  *gin.Engine
	portApp string
}

func (s *Server) initiateRoute() {
	routerGroup := s.engine.Group("/api/v1")
	controller.NewUserController(s.uS, routerGroup).Route()
	controller.NewCategoryController(s.cS, routerGroup).Route()
	controller.NewArticleController(s.aS, s.mD, routerGroup).Route()
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

// CORSMiddleware adalah middleware yang akan menangani CORS
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("CORS Middleware hit:", c.Request.Method, c.Request.URL.Path)
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}

func NewServer() *Server {
	co, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	urlConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", co.Host, co.Port, co.User, co.Password, co.Name)

	db, err := sql.Open(co.Driver, urlConnection)
	if err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()
	
	// PENTING: Pastikan middleware CORS dipasang sebelum router
	engine.Use(CORSMiddleware())

	portApp := co.AppPort

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)
	tagRepo := repository.NewTagRepository(db)
	articleTagRepo := repository.NewArticleTagRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)

	jwtService := service.NewJwtService(co.SecurityConfig)
	userService := service.NewUserservice(userRepo, jwtService)
	categoryService := service.NewCategoryService(categoryRepo)
	articleService := service.NewArticleService(articleRepo)
	bookmarkService := service.NewBookmarkService(bookmarkRepo)
	tagService := service.NewTagService(tagRepo)
	articleTagService := service.NewArticleTagService(tagRepo, articleTagRepo)
	commentService := service.NewCommentService(commentRepo)
	likeService := service.NewLikeService(likeRepo)

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	
	return &Server{
		cS:      categoryService,
		uS:      userService,
		aS:      articleService,
		bS:      bookmarkService,
		tS:      tagService,
		jS:      jwtService,
		atS:     articleTagService,
		coS:     commentService,
		lS:      likeService,
		mD: authMiddleware,
		portApp: portApp,
		engine:  engine,
	}
}