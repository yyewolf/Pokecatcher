package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func StartLogger() {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02d %02d-%02d-%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	Name := "log-" + formatted // Just like 2020-01-02 16-39-05
	Path, _ := filepath.Abs("./crash/" + Name + ".txt")
	logFile, err := os.OpenFile(Path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	//File where everything is recorded
}
