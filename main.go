package main

import (
	"context"
	"encoding/json"
	"example/fxdemo2/db"
	"example/fxdemo2/models"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

const configFile = "appconfig.yaml"

type PlayerService interface {
	//CreatePlayer(player models.Player) error
	GetPlayer(id int) (*models.Player, error)
	GetAllPlayers() (players []models.Player, err error)
	//UpdatePlayer(player Player) error
	//DeletePlayer(id int) error
}

type AppConfig struct {
	DbConfig DbConfig `yaml:"dbConfig"`
	Server   Server   `yaml:"server"`
}

type DbConfig struct {
	PostgresConn string `yaml:"postgresconn"`
}

type Server struct {
	Port int `yaml:"port"`
}

type MyService struct {
	mux           *mux.Router
	logger        *zap.Logger
	playerService PlayerService
	appConfig     *AppConfig
}

func NewMyService(mux *mux.Router, logger *zap.Logger, playersvc PlayerService, appConfig *AppConfig) *MyService {
	return &MyService{mux: mux, logger: logger, playerService: playersvc, appConfig: appConfig}
}

func newPlayerService(db *db.PostGresService) PlayerService {
	return db
}

func newAppConfig() *AppConfig {
	return &AppConfig{}
}

// newMuxRouter creates a new instance of a mux.Router, which is a HTTP request multiplexer.
// This function is used to provide a mux.Router dependency to the MyService struct.
func newMuxRouter(logger *zap.Logger) *mux.Router {
	logger.Info("Starting Router Instance")
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
	s.logger.Info("Index handler called!")
	p, err := s.playerService.GetPlayer(1)
	if err != nil {
		s.logger.Error("Error getting player", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.logger.Sugar().Infof("player %v", p.Name)
	fmt.Fprint(w, "Player: ", p.Name)
}

func (s *MyService) GetAllPlayersHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Get All Players Handler")
	players, err := s.playerService.GetAllPlayers()
	if err != nil {
		s.logger.Error("Error getting players", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(players)
	if err != nil {
		s.logger.Error("Error marshalling players to JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// write players data to response body
	_, err = w.Write(data)
	if err != nil {
		s.logger.Error("Error writing response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// Start starts the HTTP server for the MyService. It listens on port 3333 and serves requests
// using the mux.Router provided to the MyService. If there is an error starting the server,
// it is returned.
func (s *MyService) Start() error {
	port := s.appConfig.parseAppConfig().Server.Port
	PORT := fmt.Sprintf(":%v", port)
	s.logger.Sugar().Infof("Starting FXdemo2 server on port %v", PORT)
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
	s.mux.HandleFunc("/players", s.GetAllPlayersHandler)
}

// runService is a function that sets up the HTTP server for the MyService. It registers the
// IndexHandler function as the handler for the root URL path, and appends lifecycle hooks
// to the Fx application to start and stop the HTTP server.
//
// The OnStart hook starts the HTTP server in a separate goroutine, and the OnStop hook
// stops the HTTP server.
func runService(lifecyle fx.Lifecycle, s *MyService) {

	//port := s.serverConfig.parseServerConfig().Server.Port
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
			newAppConfig,
		),
		fx.Invoke(runService),
	)
	app.Run()
}

func newPostGresClient(logger *zap.Logger, appConfig *AppConfig) *pgxpool.Pool {
	// TODO: implement Postgres client initialization using existing postgres driver
	logger.Info("Starting Postgres Client")
	connString := appConfig.parseAppConfig().DbConfig.PostgresConn
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}
func (appCfg *AppConfig) parseAppConfig() AppConfig {
	var appConfig AppConfig
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		panic(err)
	}
	log.Printf("appConfigs: %v", appConfig.Server)
	return appConfig

}
