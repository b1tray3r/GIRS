package main

import (
	"embed"
	"log/slog"
	"net/http"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/b1tray3r/go-openapi3/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	redmine "github.com/nixys/nxs-go-redmine/v5"
)

//go:embed swagger/*
var StaticFiles embed.FS

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Warn("No configuration file found")
		}
	}
}

func main() {
	initConfig()

	auth := gitea.SetBasicAuth(
		viper.GetString("girs.gitea.user"),
		viper.GetString("girs.gitea.password"),
	)
	client, err := gitea.NewClient(
		viper.GetString("girs.gitea.url"),
		auth)
	if err != nil {
		panic(err.Error())
	}

	srv := server.NewEchoServer(
		"bu-contentserv-devops",
		map[string]int64{
			server.RedmineTrackerID:      viper.GetInt64("girs.redmine.tracker_id"),
			server.RedmineClosedStatusID: viper.GetInt64("girs.redmine.closed_status_id"),
		},
		redmine.Init(redmine.Settings{
			Endpoint: viper.GetString("girs.redmine.url"),
			APIKey:   viper.GetString("girs.redmine.api_key"),
		}),
		client,
	)

	srv.Echo.GET(
		"/swagger/*",
		echo.WrapHandler(
			http.StripPrefix(
				"/",
				http.FileServer(
					http.FS(
						StaticFiles,
					),
				),
			),
		),
	)

	if err := srv.Echo.Start(viper.GetString("girs.server.url")); err != nil {
		panic(err.Error())
	}
}
