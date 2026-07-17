package analyseur

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/go-logfmt/logfmt"
	"github.com/olekukonko/tablewriter"
)

// Analytics est une structure pour stocker les données d'analyse
type Analytics struct {
	Anomalies    []Anomalie
	Patterns     map[string]int
	Stats        map[string]float64
	TotalAnomalies int
}

// Anomalie est une struct pour stocker une anomalie détectée
type Anomalie struct {
	Timestamp string
	Severity  string
	Message   string
}

// DetectAnomalies détecte les anomalies dans les logs
func (a *Analytics) DetectAnomalies(ctx context.Context, logs []logfmt.Log) error {
	// Créer une copie du tableau de logs pour éviter de modifier l'original
	logsCopy := make([]logfmt.Log, len(logs))
	copy(logsCopy, logs)

	// Initialiser les stats et les patterns
	a.Stats = make(map[string]float64)
	a.Patterns = make(map[string]int)

	// Définir les seuils pour les anomalies
	anomalyThreshold := 5
	patternThreshold := 3

	// Iterer sur les logs pour détecter les anomalies
	for _, log := range logsCopy {
		// Définir le timestamp du log
		timestamp := log.Timestamp

		// Définir la sévérité du log
		severity := log.Severity

		// Définir le message du log
		message := log.Message

		// Vérifier si le log est une anomalie
		if severity == "ERROR" || severity == "CRITICAL" {
			// Créer une nouvelle anomalie
			anomaly := Anomalie{
				Timestamp: timestamp,
				Severity:  severity,
				Message:   message,
			}

			// Ajouter l'anomalie aux anomalies
			a.Anomalies = append(a.Anomalies, anomaly)

			// Incrémente le nombre d'anomalies totales
			a.TotalAnomalies++
		}

		// Vérifier si le log contient un pattern
		pattern := ""
		if matched, err := matchPattern(message, a.Patterns); err == nil && matched {
			// Incrémente le compteur du pattern
			a.Patterns[pattern]++

			// Vérifier si le pattern est une anomalie
			if a.Patterns[pattern] >= patternThreshold {
				// Créer une nouvelle anomalie
				anomaly := Anomalie{
					Timestamp: timestamp,
					Severity:  "INFO",
					Message:   fmt.Sprintf("Pattern '%s' trouvé", pattern),
				}

				// Ajouter l'anomalie aux anomalies
				a.Anomalies = append(a.Anomalies, anomaly)
			}
		}
	}

	return nil
}

// matchPattern vérifie si le message contient un pattern
func matchPattern(message string, patterns map[string]int) (bool, error) {
	// Créer une copie du tableau de patterns pour éviter de modifier l'original
	patternsCopy := make(map[string]int)
	for pattern, count := range patterns {
		patternsCopy[pattern] = count
	}

	// Iterer sur les patterns pour vérifier si le message en contient un
	for pattern, count := range patternsCopy {
		if strings.Contains(message, pattern) {
			return true, nil
		}

		// Vérifier si le nombre de occurrences du pattern est supérieur au seuil
		if occurrences := strings.Count(message, pattern); occurrences > 0 {
			// Incrémente le compteur du pattern
			patternsCopy[pattern]++

			// Vérifier si le pattern est une anomalie
			if patternsCopy[pattern] >= 5 {
				return true, nil
			}
		}
	}

	return false, nil
}

// getStats calcule les stats des anomalies
func (a *Analytics) getStats() {
	// Initialiser les stats
	for key := range a.Stats {
		a.Stats[key] = 0
	}

	// Iterer sur les anomalies pour calculer les stats
	for _, anomaly := range a.Anomalies {
		// Définir la catégorie de l'anomalie
		category := anomaly.Severity

		// Incrémente la stat pour la catégorie
		a.Stats[category]++
	}

	// Calculer la moyenne des stats
	mean := math.NaN()
	for _, stat := range a.Stats {
		mean += stat
	}
	mean /= float64(len(a.Stats))
}

// printStats imprime les stats des anomalies
func (a *Analytics) printStats() {
	// Initialiser la table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Categorie", "Compteur"})

	// Iterer sur les stats pour imprimer la table
	for key, value := range a.Stats {
		table.Append([]string{key, fmt.Sprintf("%d", value)})
	}

	// Imprimer la table
	table.Render()
}

// printAnomalies imprime les anomalies détectées
func (a *Analytics) printAnomalies() {
	// Initialiser la table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Timestamp", "Sévérité", "Message"})

	// Iterer sur les anomalies pour imprimer la table
	for _, anomaly := range a.Anomalies {
		table.Append([]string{anomaly.Timestamp, anomaly.Severity, anomaly.Message})
	}

	// Imprimer la table
	table.Render()
}

// runAnalyseur exécute l'analyseur
func (a *Analytics) runAnalyseur(ctx context.Context) error {
	// Détecter les anomalies
	if err := a.DetectAnomalies(ctx, []logfmt.Log{}); err != nil {
		return err
	}

	// Calculer les stats
	a.getStats()

	// Imprimer les stats
	a.printStats()

	// Imprimer les anomalies
	a.printAnomalies()

	return nil
}