package util

import (
	"log"
	"os"
)

var Log log.Logger = *log.New(os.Stdout, "", log.Lshortfile|log.LUTC)
