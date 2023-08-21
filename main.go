package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"riot-data-ingest/api"
	"riot-data-ingest/client"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Locales struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Locales []string `json:"locales"`
	// Maintenances []string `json:"maintenances"`
	// Incidents    []string `json:"incidents"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	c := client.NewClient(api.RegionBrasil, http.DefaultClient, logrus.New())
	r, err := c.DoRequest("GET", "/val/content/v1/contents", nil)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	fmt.Println(r.Body)
	//url := "https://br.api.riotgames.com/val/content/v1/contents?locale=pt-BR&api_key=RGAPI-95faa9dd-791f-4cb3-87df-abb14d44e1ff"

	_, err = io.Copy(os.Stdout, r.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
}
