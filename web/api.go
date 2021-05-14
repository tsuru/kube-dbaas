package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	echoPrometheus "github.com/globocom/echo-prometheus"
)

var metricsMiddleware = echoPrometheus.MetricsMiddleware()

type api struct {
	// ShutdownTimeout defines the max duration used to wait the web server
	// gracefully shutting down. Defaults to `30 * time.Second`.
	ShutdownTimeout time.Duration

	client.Client

	started  bool
	e        *echo.Echo
	shutdown chan struct{}
}

type API interface {
	Start() error
}

func New(cli client.Client) (API, error) {
	a := &api{
		ShutdownTimeout: 30 * time.Second,
		Client:          cli,
		shutdown:        make(chan struct{}),
	}
	a.e = a.newEcho()
	return a, nil
}

func (a *api) handleSignals() {
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-quit:
		a.Stop()
	case <-a.shutdown:
	}
}

func (a *api) startServer() error {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8888"
	}
	return a.e.StartH2CServer(":"+portStr, &http2.Server{})
}

func (a *api) Start() error {
	a.started = true

	go a.handleSignals()
	if err := a.startServer(); err != http.ErrServerClosed {
		fmt.Printf("problem to start the webserver: %+v", err)
		return err
	}
	fmt.Println("Shutting down the webserver...")
	return nil
}

// Stop shut down the web server.
func (a *api) Stop() error {
	if !a.started {
		return fmt.Errorf("web server is already down")
	}
	if a.shutdown == nil {
		return fmt.Errorf("shutdown channel is not defined")
	}
	close(a.shutdown)
	ctx, cancel := context.WithTimeout(context.Background(), a.ShutdownTimeout)
	defer cancel()
	return a.e.Shutdown(ctx)
}

func (a *api) newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = HTTPErrorHandler

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(metricsMiddleware)
	e.Use(ErrorMiddleware)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/healthcheck", healthcheck)

	group := e.Group("/resources")
	group.GET("/plans", servicePlans)
	group.POST("", a.serviceCreate)
	group.POST("/:instance/bind-app", a.serviceBindApp)

	/*

		group.GET("/:instance/plans", servicePlans)
		group.GET("/:instance", serviceInfo)
		group.PUT("/:instance", serviceUpdate)
		group.GET("/:instance/status", serviceStatus)
		group.DELETE("/:instance", serviceDelete)
		group.GET("/:instance/autoscale", getAutoscale)
		group.POST("/:instance/autoscale", createAutoscale)
		group.PATCH("/:instance/autoscale", updateAutoscale)
		group.DELETE("/:instance/autoscale", removeAutoscale)

		group.DELETE("/:instance/bind-app", serviceUnbindApp)
		group.POST("/:instance/bind", serviceBindUnit)
		group.DELETE("/:instance/bind", serviceUnbindUnit)
		group.GET("/:instance/info", instanceInfo)
	*/

	return e
}
