package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/olahol/melody"
)

type Application struct {
	Port          uint
	DebuggingMode bool

	Lobbies Lobbies
	m       *melody.Melody
	s       *http.Server

	*Logger
}

type Logger struct {
	logError   *log.Logger
	logWarning *log.Logger
	logInfo    *log.Logger
	logDebug   *log.Logger
}

func (app *Application) Run() {
	app.logInfo.Printf("Starting on http://localhost:%d", app.Port)
	app.logError.Fatal(app.s.ListenAndServe())
}

func (app *Application) ParseFlags() {
	flag.UintVar(&app.Port, "port", DEFAULT_PORT, "Port to run server on")
	flag.BoolVar(&app.DebuggingMode, "debug", false, "Turns on debugging mode")
	flag.Parse()
}

func (app *Application) SetupHandlers() {
	app.logInfo.Println("Setting up handlers")

	// websocket crap
	app.m = melody.New()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/game", app.gameHandler)

	app.m.HandleConnect(app.websocketConnectHandler)
	app.m.HandleMessage(app.websocketMessageHandler)
	app.m.HandleDisconnect(app.websocketDisconnectHandler)

	app.s = &http.Server{
		Addr:           fmt.Sprintf(":%d", app.Port),
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func createDebuggingLog(debuggingEnabled bool, logFlags int) *log.Logger {
	if debuggingEnabled {
		return log.New(os.Stdout, "DEBUG: ", logFlags)
	}

	return log.New(io.Discard, "DEBUG: ", logFlags) // When debugging is disabled it should discard everything
}

func (app *Application) setupLogging() {
	flags := log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix

	app.Logger = &Logger{
		logInfo:    log.New(os.Stdout, "INFO: ", flags),
		logError:   log.New(os.Stdout, "ERROR: ", flags),
		logWarning: log.New(os.Stdout, "WARN: ", flags),
		logDebug:   createDebuggingLog(app.DebuggingMode, flags),
	}

	app.logInfo.Println("Set up logging")
	if app.DebuggingMode {
		app.logInfo.Println("Running in debug mode!")
	}
}
