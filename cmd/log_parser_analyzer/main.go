// Package main est l'entrée du programme log_parser_analyzer.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"

	"log_parser_analyzer/pkg/analyseur"
	"log_parser_analyzer/pkg/log/parser"
	"log_parser_analyzer/pkg/storage"
)

var (
	// Flag pour spécifier le fichier de logs
	logFile = flag.String("log-file", "", "Fichier de logs")
	// Flag pour spécifier la durée d'analyse
	analyzeDuration = flag.Duration("analyze-duration", 0, "Durée d'analyse en secondes")
)

func main() {
	// Initialisation des flags
	flag.Parse()

	// Vérification de la présence du fichier de logs
	if *logFile == "" {
		log.Fatal("Fichier de logs non spécifié")
	}

	// Création de l'analyseur
	a := analyseur.NewAnalyzer()

	// Lancement de l'analyse du fichier de logs
	if err := a.Analyze(context.Background(), *logFile, *analyzeDuration); err != nil {
		log.Fatal(err)
	}

	// Affichage des résultats de l'analyse
	if err := printResults(a.GetResults()); err != nil {
		log.Fatal(err)
	}
}

// printResults affiche les résultats de l'analyse sous forme de tableau
func printResults(results []*analyseur.Result) error {
	// Création d'un tableau pour afficher les résultats
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Heure", "Type", "Description"})

	// Affichage des résultats
	for _, r := range results {
		table.Append([]string{r.Date, r.Hour, r.Type, r.Description})
	}
	table.Render()

	return nil
}

// SIGINT et SIGTERM pour fermer proprement le programme
func shutdown(ctx context.Context) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func init() {
	// Définition de la fonction de fermeture pour les terminaux
	term.SetRawMode(true)
	defer term.Restore()
}