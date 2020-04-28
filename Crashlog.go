package main

import (
	"fmt"
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
	_, err := os.Stat("crash")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("crash", 0755)
		if errDir != nil {
			//
		}
	}
	Path, _ := filepath.Abs("./crash/" + Name + ".txt")
	logFile, err := os.OpenFile(Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
	}
	OStdout = os.Stdout
	os.Stderr = logFile
	os.Stdout = logFile
	//File where everything is recorded
}
