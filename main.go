package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MyService struct {
	mux    *mux.Router
	logger *zap.Logger
}

func NewMyService(mux *mux.Router, logger *zap.Logger) *MyService {
	return &MyService{mux: mux, logger: logger}
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
	logger, _ := zap.NewDevelopment()
	return logger
}

func (s *MyService) IndexHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Index handler called")
	w.Write([]byte("Hello from FxDemo2 Server!"))
}

// Start starts the HTTP server for the MyService. It listens on port 3333 and serves requests
// using the mux.Router provided to the MyService. If there is an error starting the server,
// it is returned.
func (s *MyService) Start() error {
	PORT := ":3333"
	s.logger.Sugar().Infof("Starting FXdemo2 server on port %s", PORT)
	//graceful.Run(fmt.Sprintf(":%d", 3333), 10*time.Second, s.mux)
	err := http.ListenAndServe(PORT, s.mux)
	if err != nil {
		return err
	}
	return err
}

func (s *MyService) Stop() {
	s.logger.Info("Stopping fxdemo2 HTTP server")

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
				s.mux.HandleFunc("/", s.IndexHandler)
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
			zap.NewExample,
			newMuxRouter,
			//newZapLogger,
		),
		fx.Invoke(runService),
	)
	app.Run()
}
