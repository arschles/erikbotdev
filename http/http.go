package http

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/erikstmartin/erikbotdev/bot"
)

var hub *Hub

func Start(addr string, webPath string) error {
	hub = newHub()
	go hub.run()

	// Kubernetes readiness and liveness probes
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(bot.WebPath(), "public")))))
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(bot.MediaPath()))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(bot.WebPath(), "public", "index.html"))
	})

	return http.ListenAndServe(addr, nil)
}

func BroadcastMessage(msg Message) error {
	log.Printf("Broadcasting message %+v", msg)
	return hub.BroadcastMessage(msg)
}

func BroadcastChatMessage(user *bot.User, msg string) error {
	m := &ChatMessage{
		User: user,
		Text: msg,
	}

	return hub.BroadcastMessage(m)
}
