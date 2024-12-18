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

func ScrapDenyaKolha(c echo.Context) error {

	search := c.QueryParam("search")
	if search == "" {
		return c.JSON(400, "search parameter is required")
	}


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


	loadEnvError := godotenv.Load()
	if loadEnvError != nil {
		fmt.Printf("Error loading the env file: %v", loadEnvError)
		return c.JSON(500, "Error loading environment variables")
	}


	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return c.JSON(500, "API key is missing")
	}


	formatProductDetails := func(products []schema.ProductDetail) string {
		var formattedProducts []string
		for _, product := range products {
			if product.Name != "" {
				formattedProduct := fmt.Sprintf("%s;%s;%s;%s", product.Name, product.Price, product.Image, product.Link)
				formattedProducts = append(formattedProducts, formattedProduct)
			} else {
				formattedProducts = append(formattedProducts, "N/A") // if no product details found
			}
		}
		return fmt.Sprintf("%s", formattedProducts)
	}


	mytekFormatted := formatProductDetails(productsMytek)
	sbsFormatted := formatProductDetails(productsSBS)
	tunisianetFormatted := formatProductDetails(productsTunisia)


	firstPayload := "You will be given a list of products from 3 different websites, Mytek, SBS Informatique, and Tunisianet. Your task to pick the best product from each website and return the results in the following format: name;price;image;link#name;price;image;link#name;price;image;link where the first product is from Mytek, the second from SBS Informatique, and the third from Tunisianet. If you can't find a product from a website, you should return 'N/A' instead of the product details. The price should be in TND. The image should be a direct link to the image. The link should be a direct link to the product page. The products should be separated by a '#' character. The product details should be separated by a ';' character. The product details should be in the order name;price;image;link. The products should be in the order Mytek#SBS Informatique#Tunisianet. If you can't find a product from a website, you should return 'N/A' instead of the product details. and do nothing less and nothing more!"

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=" + apiKey


	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": firstPayload + "\n\nHere's the products that you asked for:\n" + mytekFormatted + "\n" + sbsFormatted + "\n" + tunisianetFormatted},
				},
			},
		},
	}


	jsonData, err := json.Marshal(payload)
	if err != nil {
		return c.JSON(500, "Error preparing request data")
	}


	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(500, "Error making request to external API")
	}
	defer resp.Body.Close()


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(500, "Error reading API response")
	}


	if resp.StatusCode != http.StatusOK {
		return c.JSON(500, fmt.Sprintf("Error from external API: %s", body))
	}


	return c.JSON(200, map[string]interface{}{
		"mytek":       productsMytek,
		"sbs":         productsSBS,
		"tunisianet":  productsTunisia,
		"AIResponse": string(body),
	})
}
