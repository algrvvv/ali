package utils

import (
	"fmt"
	"os"

	"github.com/algrvvv/ali/v2/logger"
)

func CheckError(v any) {
	if v != nil {
		logger.SaveDebugf("error: %v", v)
		fmt.Println("error occured; see ~/.ali/ali.log for more information")
		os.Exit(1)
	}
}
