package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

var (
	// Création d'une nouvelle instance du jeu (définie ailleurs)
	game = NewGame()

	// Chargement du template HTML principal depuis /templates/index.html
	// template.Must arrête le programme si le fichier n'existe pas ou contient des erreurs
	tmpl = template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))
)

// Structure envoyée au template HTML pour afficher l'état du jeu
type PageData struct {
	Board         [rows][cols]Cell // Plateau de jeu
	CurrentPlayer Cell             // Joueur en cours
	State         State            // État du jeu (victoire, en cours, égalité…)
	Message       string           // Message d'erreur ou d'information
}

func main() {

	// Déclare un serveur de fichiers statiques (CSS, JS, images…)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Définition des routes
	http.HandleFunc("/", handleHome)       // Affichage principal
	http.HandleFunc("/play", handlePlay)   // Jouer un coup
	http.HandleFunc("/reset", handleReset) // Réinitialiser le jeu

	// Lancement du serveur
	log.Println("Serveur sur http://localhost:8080 ...")
	http.ListenAndServe("localhost:8080", nil)
}

// Page d'accueil : affiche le plateau et l'état actuel du jeu
func handleHome(w http.ResponseWriter, r *http.Request) {

	// Prépare les données envoyées au template HTML
	data := PageData{
		Board:         game.Board,
		CurrentPlayer: game.CurrentPlayer,
		State:         game.State,
		Message:       game.Message,
	}

	// Exécute le template et affiche la page HTML
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Route pour jouer un coup (méthode POST)
func handlePlay(w http.ResponseWriter, r *http.Request) {

	// Empêche l'accès direct en GET → redirige vers /
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Récupère la colonne envoyée par le formulaire.
	// FormValue renvoie TOUJOURS une string → même si c'est un nombre.
	colStr := r.FormValue("column")

	// Conversion string → int
	// Atoi est nécessaire parce que les valeurs HTML sont des chaînes.
	col, err := strconv.Atoi(colStr)
	if err != nil {
		// Si la conversion échoue : message d'erreur
		game.Message = "Veuillez choisir une colonne valide."
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Joue le coup en utilisant la logique du jeu
	game.Play(col)

	// Recharge la page après l'action
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Réinitialise totalement la partie
func handleReset(w http.ResponseWriter, r *http.Request) {
	game.Reset()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
