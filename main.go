package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dannylesnik/http-inject-context/models"
	"github.com/dannylesnik/http-inject-context/webhandlers"
	"github.com/gorilla/mux"
)

//AddContext -
func AddContext(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, "-", r.RequestURI)

		//Add data to context
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func main() {

	dbase, err := models.InitDB("root:my-said2000@/test")
	if err != nil {
		log.Panic(err)
	}

	ctx := context.WithValue(context.Background(), models.SQLKEY, dbase)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/person/{id}", webhandlers.GetPerson).Methods("GET")
	router.HandleFunc("/update", webhandlers.UpdatePerson).Methods("PUT")
	router.HandleFunc("/add", webhandlers.CreatePerson).Methods("POST")
	router.HandleFunc("/person/{id}", webhandlers.DeletePerson).Methods("DELETE")

	contextedMux := AddContext(ctx, router)

	srv := &http.Server{
		Handler:      contextedMux,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)

}
