package main

import (
	"context"
	"fmt"
	"github.com/igridnet/igrid/api"
	"github.com/igridnet/igrid/env"
	"github.com/igridnet/igrid/postgres"
	uapi "github.com/igridnet/users/api"
	"github.com/igridnet/users/factory"
	"github.com/igridnet/users/hasher"
	"github.com/igridnet/users/jwt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envPostgresHost                      = "IGRID_POSTGRES_HOST"
	envPostgresPort                      = "IGRID_POSTGRES_PORT"
	envPostgresUser                      = "IGRID_POSTGRES_USER"
	envPostgresPassword                  = "IGRID_POSTGRES_PASSWORD"
	envPostgresName                      = "IGRID_POSTGRES_DB"
	envPostgresSSLMode                   = "IGRID_POSTGRES_SSLMODE"
	defPostgresHost                      = "localhost"
	defPostgresPort                      = "5432"
	defPostgresUser                      = "igridnet"
	defPostgresPassword                  = "root"
	defPostgresName                      = "beanpay"
	defPostgresSSLMode                   = "disable"
	envServerPort                        = "IGRID_SERVER_PORT"
	defServerPort                        = "8080"
	envDebugMode                         = "IGRID_DEBUG_MODE"
	defDebugMode                         = true
)

func loadDatabaseConf() *postgres.Config {
	var (
		host     = env.String(envPostgresHost, defPostgresHost)
		port     = env.String(envPostgresPort, defPostgresPort)
		user     = env.String(envPostgresUser, defPostgresUser)
		password = env.String(envPostgresPassword, defPostgresPassword)
		name     = env.String(envPostgresName, defPostgresName)
		sslMode  = env.String(envPostgresSSLMode, defPostgresSSLMode)
		//debugMoe = env.Bool(envDebugMode,defDebugMode)
	)

	return &postgres.Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
		SSLMode:  sslMode,
	}
}


func main() {

	var (
	//	debugMode     = env.Bool(envDebugMode, defDebugMode)
		serverPort    = env.String(envServerPort, defServerPort)
	//	secret        = env.String(envJWTSigningSecret, defJWTSigningSecret)
	)

	err := postgres.Initialize(context.TODO(), loadDatabaseConf())
	if err != nil {
		log.Fatalf("could not initialize database %v\n", err)
	}

	db, err := postgres.ConnectWithConfig(context.TODO(), loadDatabaseConf())
	if err != nil {
		log.Fatalf("could not connect to database %v\n", err)
	}
	tknz := jwt.NewTokenizer("secret")
	has := hasher.New()
	f := factory.NewFactory()
	us := uapi.NewClient(db,f,tknz,has)

	client := api.NewClient(us)

	handler := client.MakeHandler()

	//debugMode := env.Bool(envDebugMode, defDebugMode)
	//fmt.Printf("debug mode is %t\n", debugMode)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s",serverPort),
		Handler: handler,
	}

	go func() {
		fmt.Printf("starting server on port %s\n",serverPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")


}
