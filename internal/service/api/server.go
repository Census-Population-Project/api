package api

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/database"
	"github.com/Census-Population-Project/API/internal/service/api/middleware"

	"github.com/Census-Population-Project/API/internal/service/api/handlers/system"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type State string

const (
	Exit State = "exit"
)

type ServerInterface interface {
	InitAPI()
	Start()
}

type ServerEngine struct {
	Core    *http.Server
	Channel chan interface{}
}

type Server struct {
	WaitGroup *sync.WaitGroup
	Engine    *ServerEngine
	Config    *config.Config
	Logger    *logrus.Logger
	Database  *database.DataBase
	Redis     *redis.Client
}

func (s *Server) InitAPI() {
	r := chi.NewRouter()

	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.HttpLoggerMiddleware(s.Logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: s.Config.Server.AllowOrigins,
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Authorization"},
	}))

	systemHandlers := system.NewSystemHandler(s.Config)

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/system", systemHandlers.Router)
	})

	s.Engine.Core.Handler = r
}

func (s *Server) Start() {
	s.Logger.Info("Starting API server...")
	s.WaitGroup.Add(1)
	go func() {
		defer s.WaitGroup.Done()

		if err := s.Engine.Core.ListenAndServe(); err != nil {
			s.Logger.Fatalf("Failed to start the server: %v", err)
			s.Engine.Channel <- Exit
		}
	}()

	s.Logger.Info("API server started!")
}

func NewServerHttp(
	log *logrus.Logger, cfg *config.Config,
	db *database.DataBase, rdb *redis.Client, wg *sync.WaitGroup,
) *Server {
	return &Server{
		WaitGroup: wg,
		Engine: &ServerEngine{
			Core: &http.Server{
				Addr:    cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port),
				Handler: nil,
			},
			Channel: make(chan interface{}),
		},
		Config:   cfg,
		Logger:   log,
		Database: db,
		Redis:    rdb,
	}
}
