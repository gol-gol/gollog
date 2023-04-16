package gollog

/*
Helper GinLog to show request details
*/

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

type ginLog struct {
	LogLevel          string
	HttpStatusCode    int    `json:"status,omitempty"`
	HttpMethod        string `json:",omitempty"`
	HttpPath          string `json:"path,omitempty"`
	AfterMicroseconds int    `json:"usec,omitempty"`
	DataLength        int    `json:"len,omitempty"`

	ClientIp    string `json:"src,omitempty"`
	UserAgent   string `json:"user-agent,omitempty"`
	HttpReferer string `json:"referer,omitempty"`

	RequestedAt string `json:"at,omitempty"`
}

var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	GinTimeFormat = time.RFC3339
	GinEnableUTC  = false
)

func ginLogger(ctx *gin.Context) {
	startCtx := time.Now()
	path := ctx.Request.URL.Path

	ctx.Next()

	stopCtx := time.Since(startCtx)
	msUsed := int(math.Ceil(float64(stopCtx.Nanoseconds()) / 1000.0))

	statusCode := ctx.Writer.Status()
	clientIP := ctx.ClientIP()
	clientUserAgent := ctx.Request.UserAgent()
	referer := ctx.Request.Referer()
	dataLength := ctx.Writer.Size()
	httpMethod := ctx.Request.Method
	if dataLength < 0 {
		dataLength = 0
	}

	if GinEnableUTC {
		startCtx = startCtx.UTC()
	}
	startTime := startCtx.Format(GinTimeFormat)

	entry := ginLog{
		HttpStatusCode:    statusCode,
		HttpMethod:        httpMethod,
		HttpPath:          path,
		AfterMicroseconds: msUsed,
		DataLength:        dataLength,
		ClientIp:          clientIP,
		UserAgent:         clientUserAgent,
		HttpReferer:       referer,
		RequestedAt:       startTime,
	}

	if len(ctx.Errors) > 0 {
		entry.LogLevel = "PANIC"
	} else {
		if statusCode > 499 {
			entry.LogLevel = "ERROR"
		} else if statusCode > 399 {
			entry.LogLevel = "WARN"
		} else {
			entry.LogLevel = "INFO"
		}
	}

	logBytes, err := json.Marshal(&entry)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(logBytes))
}

func SetGinLog(router *gin.Engine) {
	router.Use(ginLogger, gin.Recovery())
}
