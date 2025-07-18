package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gogo/protobuf/version"
	"go.uber.org/zap"
	"inkgo/authentication"
	"inkgo/authentication/oauth"
	"inkgo/common"
	"inkgo/config"
	"inkgo/controller"
	"inkgo/database"
	_ "inkgo/docs"
	"inkgo/middleware"
	"inkgo/repository"
	"inkgo/service"
	"inkgo/utils/set"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const Language = "zh"

type Server struct {
	engine      *gin.Engine
	config      *config.Config
	redis       *redis.Client
	logger      *zap.Logger
	repository  repository.Repository
	controllers []controller.Controller
}

func New(conf *config.Config, logger *zap.Logger) (*Server, error) {

	// rateLimit
	rateLimitMiddleware := middleware.RateLimitMiddleware(&conf.Server.RateLimitsConfigs)

	// mysql
	db, err := database.NewMysql(&conf.DB)
	if err != nil {
		return nil, err
	}

	// initable redis
	rdb, err := database.NewRedis(&conf.Redis)
	if err != nil {
		return nil, err
	}

	//new initable repository
	repository := repository.NewRepository(db, rdb)
	if conf.DB.Migrate {
		if err = repository.Migrate(); err != nil {
			return nil, err
		}
	}

	// user
	userService := service.NewUserService(repository.User())
	userController := controller.NewUserController(userService)
	// authentication
	tokenService := service.NewTokenService(repository.Token())
	authService := service.NewAuthService(repository.Auth())
	jwtService := authentication.NewJWT(&conf.JWTConfig, tokenService)
	oauthManager := oauth.NewOAuthManager(conf.OAuthConfigs)
	authContoller := controller.NewAuthController(userService, jwtService, oauthManager, authService)
	// Post
	PostService := service.NewPostService(repository.Post())
	PostController := controller.NewPostController(PostService)

	//category
	categoryService := service.NewCategoryService(repository.Category())
	categoryController := controller.NewCategoryController(categoryService)

	//tag
	tagService := service.NewTagService(repository.Tag())
	tagController := controller.NewTagController(tagService)

	// favorite
	favoriteService := service.NewFavoriteService(repository.Favorite())
	favoriteController := controller.NewFavoriteController(favoriteService)

	//comment
	commentService := service.NewCommentService(repository.Comment())
	commentController := controller.NewCommentController(commentService)

	// like
	likeService := service.NewLikeService(repository.Like())
	likeController := controller.NewLikeController(likeService)

	controllers := []controller.Controller{userController, PostController, categoryController, tagController, authContoller, commentController, likeController, favoriteController}

	//logger
	logs := service.NewLoggerService(&conf.Logger)
	if err := logs.WriteLog(); err != nil {
		return nil, err
	}

	gin.SetMode(conf.Server.ENV)
	e := gin.New()

	e.Use(
		rateLimitMiddleware,
		middleware.CORSMiddleware(),
		middleware.LoggerMiddleWare(),
		middleware.Recovery(),
	)

	//e.LoadHTMLFiles("")

	return &Server{
		engine:      e,
		config:      conf,
		logger:      logger,
		repository:  repository,
		controllers: controllers,
	}, nil
}

func (s *Server) Routers() {
	root := s.engine
	// register non-resource routers
	root.GET("/", common.WrapFunc(s.getRoutes))
	root.GET("/index", controller.Index)
	root.GET("/healthz", common.WrapFunc(s.Ping))
	root.GET("/version", common.WrapFunc(version.Get))
	//root.GET("/metrics", gin.WrapH(promhttp.Handler()))
	root.Any("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))

	// swagger doc
	if gin.Mode() != gin.ReleaseMode {
		root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	public := root.Group("/api/v1")

	api := root.Group("/api/v1")
	api.Use(middleware.AuthenticationMiddleware(authentication.NewJWT(&s.config.JWTConfig, s.repository.Token()), s.repository.User()))
	controllers := make([]string, 0, len(s.controllers))
	for _, router := range s.controllers {
		switch router.Name() {
		case "auth":
			//case "email", "auth", "users":
			// 这些路由不需要登录验证
			router.RegisterRoute(public)
		default:
			// 其他路由需要登录验证
			router.RegisterRoute(api)
		}
		controllers = append(controllers, router.Name())
	}
	zap.S().Infof("server enabled controllers: %v", controllers)
}

func (s *Server) Run() error {
	defer s.Close()
	s.Routers()

	addr := fmt.Sprintf("%s:%d", s.config.Server.Address, s.config.Server.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("failed to start server, %v", err)
		}
	}()

	// 平滑关闭进程
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Server.GracefulShutdownPeriod)*time.Second)
	defer cancel()

	ch := <-sig
	zap.S().Infof("Receive signal: %s", ch)

	return server.Shutdown(ctx)
}

func (s *Server) Close() {
	if err := s.repository.Close(); err != nil {
		zap.S().Warnf("failed to close repository, %v", err)
	}
}

func (s *Server) getRoutes() []string {
	paths := set.NewString()
	for _, r := range s.engine.Routes() {
		if r.Path != "" {
			paths.Insert(r.Path)
		}
	}
	return paths.Slice()
}

type ServerStatus struct {
	Ping         bool `json:"ping"`
	DBRepository bool `json:"dbRepository"`
}

func (s *Server) Ping() *ServerStatus {
	status := &ServerStatus{Ping: true}

	ctx, cannel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cannel()

	if err := s.repository.Ping(ctx); err == nil {
		status.DBRepository = true
	}

	return status
}

// @Summary Healthz
// @Produce json
// @Tags healthz
// @Success 200 {string}  string    "ok"
// @Router /healthz [get]
func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
