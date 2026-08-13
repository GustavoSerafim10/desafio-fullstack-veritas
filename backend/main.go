package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewTaskStore("tasks.json")
	taskHandler := NewTaskHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/tasks", taskHandler)
	mux.Handle("/tasks/", taskHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	handler := withLogging(withCORS(mux))

	const addr = ":8080"
	log.Printf("Mini Kanban API rodando em http://localhost%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
