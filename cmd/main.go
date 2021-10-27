package main

import (
	"context"
	"fmt"
	"github.com/igridnet/igrid/api"
	"github.com/igridnet/igrid/env"
	"github.com/igridnet/igrid/mq"
	"github.com/igridnet/igrid/postgres"
	mflog "github.com/igridnet/mproxy/logger"
	"github.com/igridnet/mproxy/pkg/mqtt"
	"github.com/igridnet/mproxy/pkg/session"
	"github.com/igridnet/mproxy/pkg/websocket"
	uapi "github.com/igridnet/users/api"
	"github.com/igridnet/users/factory"
	"github.com/igridnet/users/hasher"
	"github.com/igridnet/users/jwt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	// WS
	defWSHost = "0.0.0.0"
	defWSPath = "/mqtt"
	defWSPort = "8080"

	defWSTargetScheme = "ws"
	defWSTargetHost   = "mosquitto"
	defWSTargetPort   = "8888"
	defWSTargetPath   = "/mqtt"

	envWSHost         = "IGRID_WS_HOST"
	envWSPort         = "IGRID_WS_PORT"
	envWSPath         = "IGRID_WS_PATH"
	envWSTargetScheme = "IGRID_WS_TARGET_SCHEME"
	envWSTargetHost   = "IGRID_WS_TARGET_HOST"
	envWSTargetPort   = "IGRID_WS_TARGET_PORT"
	envWSTargetPath   = "IGRID_WS_TARGET_PATH"

	// MQTT
	defMQTTHost       = "0.0.0.0"
	defMQTTPort       = "1883"
	defMQTTTargetHost = "mosquitto"
	defMQTTTargetPort = "1884"

	envMQTTHost       = "IGRID_MQTT_HOST"
	envMQTTPort       = "IGRID_MQTT_PORT"
	envMQTTSPort      = "IGRID_MQTTS_PORT"
	envMQTTTargetHost = "IGRID_MQTT_TARGET_HOST"
	envMQTTTargetPort = "IGRID_MQTT_TARGET_PORT"

	defLogLevel = "debug"
	envLogLevel = "IGRID_LOG_LEVEL"
)

const (
	envPostgresHost     = "IGRID_POSTGRES_HOST"
	envPostgresPort     = "IGRID_POSTGRES_PORT"
	envPostgresUser     = "IGRID_POSTGRES_USER"
	envPostgresPassword = "IGRID_POSTGRES_PASSWORD"
	envPostgresName     = "IGRID_POSTGRES_DB"
	envPostgresSSLMode  = "IGRID_POSTGRES_SSLMODE"
	defPostgresHost     = "localhost"
	defPostgresPort     = "5432"
	defPostgresUser     = "igridnet"
	defPostgresPassword = "root"
	defPostgresName     = "beanpay"
	defPostgresSSLMode  = "disable"
	envServerPort       = "IGRID_SERVER_PORT"
	defServerPort       = "8080"
	envDebugMode        = "IGRID_DEBUG_MODE"
	defDebugMode        = true
)

func loadDatabaseConf() *postgres.Config {
	var (
		host     = env.String(envPostgresHost, defPostgresHost)
		port     = env.String(envPostgresPort, defPostgresPort)
		user     = env.String(envPostgresUser, defPostgresUser)
		password = env.String(envPostgresPassword, defPostgresPassword)
		name     = env.String(envPostgresName, defPostgresName)
		sslMode  = env.String(envPostgresSSLMode, defPostgresSSLMode)
		//debugMoe = envRead.Bool(envDebugMode,defDebugMode)
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
		//	debugMode     = envRead.Bool(envDebugMode, defDebugMode)
		serverPort = env.String(envServerPort, defServerPort)
	//	secret        = envRead.String(envJWTSigningSecret, defJWTSigningSecret)
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
	us := uapi.NewClient(db, f, tknz, has)

	client := api.NewClient(us)

	handler := client.MakeHandler()

	cfg := loadConfig()

	logger, err := mflog.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	h := mq.New(logger, us)
	errs := make(chan error, 3)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: handler,
	}

	logger.Info(fmt.Sprintf("Starting WebSocket proxy on port %s ", cfg.wsPort))
	go proxyWS(cfg, logger, h, errs)


	logger.Info(fmt.Sprintf("Starting MQTT proxy on port %s ", cfg.mqttPort))
	go proxyMQTT(cfg, logger, h, errs)

	logger.Info(fmt.Sprintf("Starting registry on port %s ", serverPort))
	go runServer(srv,errs)


	go func() {
		quit := make(chan os.Signal, 2)
		signal.Notify(quit, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-quit)
	}()

	<-errs
	log.Println("Shutdown Server ...")
}

type config struct {
	wsHost         string
	wsPort         string
	wsPath         string
	wsTargetScheme string
	wsTargetHost   string
	wsTargetPort   string
	wsTargetPath   string

	mqttHost       string
	mqttPort       string
	mqttsPort      string
	mqttTargetHost string
	mqttTargetPort string

	logLevel string
}

func envRead(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func loadConfig() config {

	return config{
		// WS
		wsHost:         envRead(envWSHost, defWSHost),
		wsPort:         envRead(envWSPort, defWSPort),
		wsPath:         envRead(envWSPath, defWSPath),
		wsTargetScheme: envRead(envWSTargetScheme, defWSTargetScheme),
		wsTargetHost:   envRead(envWSTargetHost, defWSTargetHost),
		wsTargetPort:   envRead(envWSTargetPort, defWSTargetPort),
		wsTargetPath:   envRead(envWSTargetPath, defWSTargetPath),

		// MQTT
		mqttHost:       envRead(envMQTTHost, defMQTTHost),
		mqttPort:       envRead(envMQTTPort, defMQTTPort),
		mqttTargetHost: envRead(envMQTTTargetHost, defMQTTTargetHost),
		mqttTargetPort: envRead(envMQTTTargetPort, defMQTTTargetPort),

		// Log
		logLevel: envRead(envLogLevel, defLogLevel),
	}
}

func runServer(srv *http.Server,errs chan error){
	errs <- srv.ListenAndServe()
}

func proxyWS(cfg config, logger mflog.Logger, handler session.Handler, errs chan error) {
	target := fmt.Sprintf("%s:%s", cfg.wsTargetHost, cfg.wsTargetPort)
	wp := websocket.New(target, cfg.wsTargetPath, cfg.wsTargetScheme, handler, logger)
	http.Handle(cfg.wsPath, wp.Handler())

	errs <- wp.Listen(cfg.wsPort)
}

func proxyMQTT(cfg config, logger mflog.Logger, handler session.Handler, errs chan error) {
	address := fmt.Sprintf("%s:%s", cfg.mqttHost, cfg.mqttPort)
	target := fmt.Sprintf("%s:%s", cfg.mqttTargetHost, cfg.mqttTargetPort)
	mp := mqtt.New(address, target, handler, logger)

	errs <- mp.Listen()
}
