package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"scrapper/schema"
	"scrapper/utils"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// Define the structure of the API response from the external service
type Candidate struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}

type AIResponse struct {
	Candidates []Candidate `json:"candidates"`
}

// ScrapDenyaKolha handles the scraping of product data from multiple websites and querying the AI API
func ScrapDenyaKolha(c echo.Context) error {
	// Extract the search query parameter from the request
	search := c.QueryParam("search")
	if search == "" {
		return c.JSON(400, "search parameter is required")
	}

	// Scrape product data from various websites
	productsMytek, err := utils.ScrapperFromMytek("https://www.mytek.tn/catalogsearch/result/?q=", search)
	if err != nil {
		return c.JSON(500, "Error during scraping Mytek")
	}

	productsSBS, err := utils.ScrapperFromSBS("https://www.sbsinformatique.com/recherche?controller=search&s=", search)
	if err != nil {
		return c.JSON(500, "Error during scraping SBS Informatique")
	}

	productsTunisia, err := utils.ScrapperFromTunisianet("https://www.tunisianet.com.tn/recherche?controller=search&orderby=price&orderway=asc&s=", search)
	if err != nil {
		return c.JSON(500, "Error during scraping Tunisianet")
	}

	// Load environment variables
	loadEnvError := godotenv.Load()
	if loadEnvError != nil {
		fmt.Printf("Error loading the env file: %v", loadEnvError)
		return c.JSON(500, "Error loading environment variables")
	}

	// Retrieve the API key from environment variables
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return c.JSON(500, "API key is missing")
	}

	// Format product details into a string
	formatProductDetails := func(products []schema.ProductDetail) string {
		var formattedProducts []string
		for _, product := range products {
			if product.Name != "" {
				formattedProduct := fmt.Sprintf("%s;%s;%s;%s", product.Name, product.Price, product.Image, product.Link)
				formattedProducts = append(formattedProducts, formattedProduct)
			} else {
				formattedProducts = append(formattedProducts, "N/A") // If no product details found
			}
		}
		return fmt.Sprintf("%s", formattedProducts)
	}

	// Format the product details for each source
	mytekFormatted := formatProductDetails(productsMytek)
	sbsFormatted := formatProductDetails(productsSBS)
	tunisianetFormatted := formatProductDetails(productsTunisia)

	// Prepare the prompt to send to the AI API
	firstPayload := "You will be given a list of products from 3 different websites, Mytek, SBS Informatique, and Tunisianet. Your task to pick the best product from everything, only one product and return the results in the following format: name;price;image;link"

	// API endpoint for generating content from the external AI model
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=" + apiKey

	// Prepare the payload to send to the AI model
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": firstPayload + "\n\nHere's the products that you asked for:\n" + mytekFormatted + "\n" + sbsFormatted + "\n" + tunisianetFormatted},
				},
			},
		},
	}

	// Marshal the payload into JSON format
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return c.JSON(500, "Error preparing request data")
	}

	// Make the API request to the external service
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(500, "Error making request to external API")
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(500, "Error reading API response")
	}

	// Check for any error in the response status code
	if resp.StatusCode != http.StatusOK {
		return c.JSON(500, fmt.Sprintf("Error from external API: %s", body))
	}

	// Unmarshal the API response into the AIResponse struct
	var aiResponse AIResponse
	err = json.Unmarshal(body, &aiResponse)
	if err != nil {
		return c.JSON(500, "Error parsing API response")
	}

	// Extract the product details from the AI response
	if len(aiResponse.Candidates) > 0 && len(aiResponse.Candidates[0].Content.Parts) > 0 {
		extractedText := aiResponse.Candidates[0].Content.Parts[0].Text

		// Return the product data along with the AI response
		return c.JSON(200, map[string]interface{}{
			"mytek":       productsMytek,
			"sbs":         productsSBS,
			"tunisianet":  productsTunisia,
			"AIResponse": extractedText, // This contains the extracted text from the AI response
		})
	}

	// If no candidates were found or the text is missing
	return c.JSON(500, "No product details found in AI response")
}
