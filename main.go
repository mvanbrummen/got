package main

import (
	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"

	"github.com/mvanbrummen/got-std/handler"
	"github.com/mvanbrummen/got-std/util"
)

func init() {
	// init config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// create the app directory
	util.InitGotDir(viper.GetString("got.dir"))
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	handler := &handler.Handler{}

	r.Static("/static", "./static")
	r.GET("/repository/:repoName", handler.RepositoryHandler)
	r.GET("/repository/:repoName/blob/*rest", handler.FileHandler)
	r.POST("/repository/:repoName", handler.RepositoryHandlerPost)
	r.GET("/repository/:repoName/commits", handler.CommitsHandler)
	r.GET("/", handler.IndexHandler)

	r.Run(viper.GetString("application.port"))
}
