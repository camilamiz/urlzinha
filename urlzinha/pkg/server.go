package urlzinha

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func StartServer() {
	r := mux.NewRouter()
	postUrlHanlder := &PostUrlHandler{}
	r.HandleFunc("/", postUrlHanlder.Handle).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("server was shutdown unexpectedly: %s", err)
	}
	fmt.Println("server was shutdown gracefully.")
}
