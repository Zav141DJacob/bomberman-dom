package main

import (
	"fmt"
	"net/http"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from api! %s", time.Now())
}

func (app *Application) gameHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.m.HandleRequest(w, r); err != nil {
		app.logError.Println(err)
	}
}
