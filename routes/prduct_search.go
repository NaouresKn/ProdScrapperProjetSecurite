package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recherche-produit-go/models"
	"recherche-produit-go/scrapers"
)

type SearchRequest struct {
	ProductName string `json:"product"`
}

type SearchResponse struct {
	Message string   `json:"message"`
	Results []string `json:"results"`
}

func ProductSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// Scrape data from websites
	results := scrapers.ScrapeSites(req.ProductName)

	// Analyze with CNN (fake example here)
	analysis := models.AnalyzeWithCNN(req.ProductName)

	// Respond with combined results
	resp := SearchResponse{
		Message: fmt.Sprintf("Recherche pour '%s' terminée. CNN: %s", req.ProductName, analysis),
		Results: results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
