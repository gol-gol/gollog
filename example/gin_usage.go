package main

import (
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gol-gol/gollog"
)

var HTTPAt = flag.String("bind", ":8080", "Address to bind at.")

func main() {
	flag.Parse()
	ginUp(*HTTPAt)
}

func ginUp(listenAt string) {
	router := gin.Default()
	gollog.SetGinLog(router)
	router.GET("/", rootHandler)
	router.MaxMultipartMemory = 1
	router.Run(listenAt)
}

func rootHandler(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		map[string]string{"status": "OK"},
	)
}
