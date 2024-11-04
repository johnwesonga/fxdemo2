package main

import (
	"context"
	"example/fxdemo2/db"
	"example/fxdemo2/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type PlayerService interface {
	//CreatePlayer(player models.Player) error
	GetPlayer(id int) (*models.Player, error)
	//GetAllPlayers() (players []models.Player, err error)
	//UpdatePlayer(player Player) error
	//DeletePlayer(id int) error
}

type MyService struct {
	mux           *mux.Router
	logger        *zap.Logger
	playerService PlayerService
}

func NewMyService(mux *mux.Router, logger *zap.Logger, playersvc PlayerService) *MyService {
	return &MyService{mux: mux, logger: logger, playerService: playersvc}
}

func newPlayerService(db *db.PostGresService) PlayerService {
	return db
}

// newMuxRouter creates a new instance of a mux.Router, which is a HTTP request multiplexer.
// This function is used to provide a mux.Router dependency to the MyService struct.
func newMuxRouter() *mux.Router {
	router := mux.NewRouter()
	return router
}

// newZapLogger creates a new instance of a zap.Logger in development mode.
// This function is used to provide a zap.Logger dependency to the MyService struct.
func newZapLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)
	var err error
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	//logger, _ := zap.NewDevelopment()
	return logger
}

func (s *MyService) IndexHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Index handler called")
	p, err := s.playerService.GetPlayer(1)
	if err != nil {
		s.logger.Error("Error getting player", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.logger.Sugar().Infof("player %v", p.Name)
	w.Write([]byte("Hello from FxDemo2 Server!"))
}

// Start starts the HTTP server for the MyService. It listens on port 3333 and serves requests
// using the mux.Router provided to the MyService. If there is an error starting the server,
// it is returned.
func (s *MyService) Start() error {
	PORT := ":3333"
	s.logger.Sugar().Infof("Starting FXdemo2 server on port %s", PORT)
	err := http.ListenAndServe(PORT, s.mux)
	if err != nil {
		return err
	}
	return err
}

func (s *MyService) Stop() {
	s.logger.Info("Stopping fxdemo2 HTTP server")
}

func (s *MyService) LoadHandlers() {
	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	s.mux.HandleFunc("/", s.IndexHandler)
}

// runService is a function that sets up the HTTP server for the MyService. It registers the
// IndexHandler function as the handler for the root URL path, and appends lifecycle hooks
// to the Fx application to start and stop the HTTP server.
//
// The OnStart hook starts the HTTP server in a separate goroutine, and the OnStop hook
// stops the HTTP server.
func runService(lifecyle fx.Lifecycle, s *MyService) {

	lifecyle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				s.LoadHandlers()
				go s.Start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				s.Stop()
				return nil
			},
		},
	)
}

// main is the entry point for the FxDemo2 application. It creates an Fx application that
// provides a MyService, a zap.Logger, and a mux.Router, and invokes the runService function
// to start the service. The application is then run.
func main() {
	// app is an Fx application that provides a MyService, a zap.Logger, and a mux.Router,
	// and invokes the runService function to start the service.
	app := fx.New(

		fx.Provide(
			NewMyService,
			//zap.NewExample,
			newMuxRouter,
			newZapLogger,
			//newRedisClient,
			newPostGresClient,
			db.NewPostGresService,
			newPlayerService,
		),
		fx.Invoke(runService),
	)
	app.Run()
}

func newPostGresClient() *pgxpool.Pool {
	// TODO: implement Postgres client initialization using existing postgres driver
	connString := "postgres://postgres@localhost:5432/players"
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}
