package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"vmware-exporter/internal/collector"
	"vmware-exporter/internal/helper"
	_ "vmware-exporter/internal/metrics"
)

var (
	listen   = ":9512"
	host     = ""
	username = ""
	password = ""
	features = "host"
	logLevel = "info"
)

func main() {
	flag.StringVar(&listen, "listen", env("ESX_LISTEN", listen), "listen port")
	flag.StringVar(&host, "host", env("ESX_HOST", host), "ESX host addr")
	flag.StringVar(&username, "username", env("ESX_USERNAME", username), "user for ESX")
	flag.StringVar(&password, "password", env("ESX_PASSWORD", password), "password for ESX")
	flag.StringVar(&features, "features", env("ESX_FEATURES", features), "enabled collectors")
	flag.StringVar(&logLevel, "logLevel", env("ESX_LOG", logLevel), "Log level")
	flag.Parse()

	logger := initLogger()

	if host == "" {
		logger.Fatal("Yor must configured system env ESX_HOST or key -host")
	}
	if username == "" {
		logger.Fatal("Yor must configured system env ESX_USERNAME or key -username")
	}
	if password == "" {
		logger.Fatal("Yor must configured system env ESX_PASSWORD or key -password")
	}

	cli, err := helper.NewClient(host, username, password)
	if err != nil {
		logger.Fatal(err)
	}
	defer cli.Logout(context.TODO())

	collect := collector.NewService(cli, logger)
	if err := collect.Start(features); err != nil {
		logger.Fatal(err)
	}
	defer collect.Stop()

	handler := promhttp.Handler()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if err := collect.Scrape(); err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		handler.ServeHTTP(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>VMware Exporter</title></head>
			<body>
			<h1>VMware Exporter</h1>
			<p><a href="` + "/metrics" + `">Metrics</a></p>
			</body>
			</html>`))
	})

	server := &http.Server{Addr: listen}
	defer func() {
		if err := server.Shutdown(context.TODO()); err != nil {
			logger.Error(err)
		}
	}()

	go func() {
		logger.Infof("Exporter start on port %s", listen)
		defer logger.Info("Exporter shutdown")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error(err)
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)

	<-wait
}

func env(key, def string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return def
}

func initLogger() *log.Logger {
	logger := log.New()

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		level = log.InfoLevel
	}

	logger.SetLevel(level)
	logger.Formatter = &log.TextFormatter{DisableTimestamp: false, FullTimestamp: true}

	return logger
}
