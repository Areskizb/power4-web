package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

var (
	game = NewGame()
	tmpl = template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))
)

type PageData struct {
	Board         [rows][cols]Cell
	CurrentPlayer Cell
	State         State
	Message       string
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/play", handlePlay)
	http.HandleFunc("/reset", handleReset)

	log.Println("Serveur sur http://localhost:8080 ...")
	http.ListenAndServe("localhost:8080", nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Board:         game.Board,
		CurrentPlayer: game.CurrentPlayer,
		State:         game.State,
		Message:       game.Message,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	colStr := r.FormValue("column")
	col, err := strconv.Atoi(colStr)
	if err != nil {
		game.Message = "Veuillez choisir une colonne valide."
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	game.Play(col)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleReset(w http.ResponseWriter, r *http.Request) {
	game.Reset()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
