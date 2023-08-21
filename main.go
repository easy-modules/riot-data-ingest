package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Locales struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Locales []string `json:"locales"`
	// Maintenances []string `json:"maintenances"`
	// Incidents    []string `json:"incidents"`
}

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Printf("No such .env file")
		os.Exit(1)
	}
	riot_token := os.Getenv("RIOT_TOKEN")
	fmt.Printf("Riot Token=%s\n", riot_token)

	url := "https://br.api.riotgames.com/val/content/v1/contents?locale=pt-BR&api_key=RGAPI-95faa9dd-791f-4cb3-87df-abb14d44e1ff"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf(resp.Status)
		return
	}

	fmt.Println("Response Body:")
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
}
