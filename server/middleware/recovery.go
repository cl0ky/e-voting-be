package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	logDir    = "log"
	errorFile = "error-panic.log"
)

func writeToErrorFile(s string) {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			log.Printf("Error creating directory: %v", err)
			return
		}
	}

	file, err := os.OpenFile(path.Join(logDir, errorFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening error file: %v", err)
		return
	}
	defer file.Close()

	timenow := time.Now().Format("2006-01-02 15:04:05")
	payload := fmt.Sprintf("==============BEGIN ERROR %s================\n%s\n=============================================================\n", timenow, s)

	_, err = file.WriteString(payload)
	if err != nil {
		log.Printf("Error writing to error file: %v", err)
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				errStack := fmt.Sprintf("%s\n%s", r, debug.Stack())
				writeToErrorFile(errStack)

				log.Printf("\033[31m\n%s\n%s\033[0m", r, debug.Stack())

				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
