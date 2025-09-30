package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func ReqLog() gin.HandlerFunc {
	return gin.LoggerWithFormatter(reqlog)
}

func reqlog(param gin.LogFormatterParams) string {
	logmap := make(map[string]any, 0)
	logmap["Host"] = param.ClientIP
	logmap["Time"] = param.TimeStamp.Format(time.RFC1123)
	logmap["Method"] = param.Method
	logmap["Path"] = param.Path
	logmap["HTTP Version"] = param.Request.Proto
	logmap["Status Code"] = param.StatusCode
	logmap["Response Time"] = param.Latency
	logmap["User Agent"] = param.Request.UserAgent()
	logmap["Error"] = param.ErrorMessage

	logbyte, _ := json.MarshalIndent(logmap, "", "  ")
	return fmt.Sprintf("\n================================================\n%s \n================================================\n", logbyte)
}
