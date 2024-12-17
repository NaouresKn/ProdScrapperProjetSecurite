package main

import (
	"fmt"
	"log"
	"net/http"
	"recherche-produit-go/routes"
)

func main() {
	fmt.Println("Démarrage du serveur Go...")

	// Setup routes
	http.HandleFunc("/search", routes.ProductSearchHandler)

	// Start the server
	port := ":8080"
	fmt.Println("Serveur démarré sur http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
