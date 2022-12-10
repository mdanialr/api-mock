package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mdanialr/api-mock/api"
	"github.com/spf13/viper"
)

func main() {
	vApp := initAppConfig()
	if err := vApp.ReadInConfig(); err != nil {
		log.Fatalln("failed to init app config:", err)
	}

	responsePath := vApp.GetString("response_dir")
	responseName := vApp.GetString("response_name")
	vResp := initResponseConfig(responsePath, responseName)
	if err := vResp.ReadInConfig(); err != nil {
		log.Fatalln("failed to read response config:", err)
	}
	vResp.WatchConfig() // enable live reload config on changes

	handler := api.NewMainHandler(vResp)
	e := echo.New()
	e.Use(middleware.Logger())

	e.Any("/*", handler.Entry)

	host := vApp.GetString("host")
	port := vApp.GetInt("port")
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", host, port)))
}

// initAppConfig read app config and set default value as necessary.
func initAppConfig() *viper.Viper {
	vApp := viper.New()
	vApp.SetConfigName("app")
	vApp.SetConfigType("yaml")
	vApp.AddConfigPath(".")

	vApp.SetDefault("host", "127.0.0.1")
	vApp.SetDefault("port", 9595)
	vApp.SetDefault("response", ".")

	return vApp
}

// initResponseConfig read response data set config.
func initResponseConfig(dir, name string) *viper.Viper {
	if name == "" {
		name = "response" // default config name is response
	}
	if dir == "" {
		dir = "." // default config path is current directory
	}

	vResp := viper.New()
	vResp.SetConfigName(name)
	vResp.SetConfigType("yaml")
	vResp.AddConfigPath(dir)

	return vResp
}
