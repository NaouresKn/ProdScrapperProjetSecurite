package models

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type CNNResponse struct {
	PredictedClass string  `json:"predicted_class"`
	Confidence     float64 `json:"confidence"`
}

func AnalyzeWithCNN(product string) string {
	client := resty.New()

	// Replace with your TensorFlow Serving URL
	url := "http://localhost:8501/v1/models/cnn:predict"

	// Dummy payload
	payload := map[string]interface{}{
		"instances": []map[string]string{{"input_text": product}},
	}

	// Send POST request
	resp, err := client.R().
		SetBody(payload).
		Post(url)

	if err != nil {
		fmt.Println("Erreur CNN:", err)
		return "Analyse échouée"
	}

	return string(resp.Body())
}
