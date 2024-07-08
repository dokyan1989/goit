package main

import (
	"log"

	"github.com/fatih/color"
)

func logInfo(format string, a ...interface{}) {
	log.Println(color.HiGreenString("[INFO] "+format, a...))
}

func logFatal(format string, a ...interface{}) {
	log.Fatal(color.HiRedString("[ERROR] "+format, a...))
}

func logSlice(slice []string, prefix string) {
	for _, s := range slice {
		log.Printf("%s %s \n", prefix, s)
	}
}
