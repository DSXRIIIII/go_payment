package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunHttpServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("http-addr")
	RunHttpServerOnAddr(addr, wrapper)

}

func RunHttpServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	apiRouter := gin.New()
	wrapper(apiRouter)
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}
