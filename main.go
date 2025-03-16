package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fiatjaf/khatru"
	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/afero"
)

var (
	relay  = khatru.NewRelay()
	config = loadConfig()
	fs     afero.Fs
)

func main() {
	nostr.InfoLogger = log.New(io.Discard, "", 0)
	green := "\033[32m"
	reset := "\033[0m"
	fmt.Println(green + art + reset)
	log.Println("🚀 Mattress is booting up")
	fs = afero.NewOsFs()

	initRelay()

	go func() {
		refreshTrustNetwork()

		go backupDatabase()
	}()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates/static"))))
	http.HandleFunc("/", relayHandler)

	addr := fmt.Sprintf("%s:%d", config.RelayBindAddress, config.RelayPort)

	log.Printf("🔗 listening at %s", addr)
	http.ListenAndServe(addr, nil)
}

func relayHandler(w http.ResponseWriter, r *http.Request) {
	var relay *khatru.Relay
	urlPath := r.URL.Path

	if urlPath == "/relay" {
		relay = relay
	}

	relay.ServeHTTP(w, r)
}
