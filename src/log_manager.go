package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// LogManager will manage any logging needed by this application.
type LogManager struct {
	LogPath string // The root path to store log files under.
	Enabled bool   // If logging should happen or not.
}

// Initialize setups up the log manager and creates any needed files.
func (manager LogManager) Initialize() {
	if !manager.Enabled {
		return
	}

	filePath := manager.LogPath + "/" + time.Now().Format("2006-01-02")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Path does not exist creating path")
		os.MkdirAll(filePath, 0777)
	}

	var file string
	var iter int
	for {
		file = filePath + "/" + "sonalog" + strconv.Itoa(iter)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			break
		}
		iter++
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Unable to create log file")
	}

	log.SetOutput(f)
}

// LogFatal logs a fatal entry but only if logging is enabled.
func (manager LogManager) LogFatal(v ...interface{}) {
	if !manager.Enabled {
		return
	}

	log.Fatal(v)
}

// LogFatalf logs a fatal formated entry but only if logging is enabled.
func (manager LogManager) LogFatalf(format string, v ...interface{}) {
	if !manager.Enabled {
		return
	}

	log.Fatalf(format, v)
}

// LogFatalln logs a fatal entry but only if logging is enabled.
func (manager LogManager) LogFatalln(v ...interface{}) {
	if !manager.Enabled {
		return
	}

	log.Fatalln(v)
}

// LogPanic logs a panic entry but only if logging is enabled.
func (manager LogManager) LogPanic(v ...interface{}) {
	if !manager.Enabled {
		return
	}

	log.Panic(v)
}

// LogPanicf logs a formatted panic entry but only if logging is enabled.
func (manager LogManager) LogPanicf(format string, v ...interface{}) {
	if !manager.Enabled {
		return
	}

	log.Panicf(format, v)
}

// LogPanicln logs a panic entry but only if logging is enabled.
func (manager LogManager) LogPanicln(v ...interface{}) {
	if !manager.Enabled {
		return
	}

	log.Panicln(v)
}

// LogPrint logs an entry but only if logging is enabled.
func (manager LogManager) LogPrint(v ...interface{}) {
	if !manager.Enabled {
		return
	}

	log.Print(v)
}

// LogPrintf logs a fomratted entry but only if logging is enabled.
func (manager LogManager) LogPrintf(format string, v ...interface{}) {
	if !manager.Enabled {
		return
	}

	log.Printf(format, v)
}

// LogPrintln logs an entry but only if logging is enabled.
func (manager LogManager) LogPrintln(v ...interface{}) {
	if !manager.Enabled {
		return
	}

	log.Println(v)
}
