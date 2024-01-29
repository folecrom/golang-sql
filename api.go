package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Structure pour définir la réponse JSON de l'API OpenWeatherMap
type WeatherResponse struct {
	Main struct {
		Temperature float64 `json:"temp"`
	} `json:"main"`
}

func main() {
	// Remplacez "VOTRE_CLE_API" par votre clé API OpenWeatherMap
	apiKey := "68ff09cbdc1317630b3f955444948c1b"
	city := "Nantes" // Remplacez par la ville de votre choix

	// Construire l'URL de l'API OpenWeatherMap
	apiURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)

	// Effectuer la requête HTTP
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return
	}
	defer response.Body.Close()

	// Lire le corps de la réponse
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du corps de la réponse:", err)
		return
	}

	// Décodez la réponse JSON dans la structure WeatherResponse
	var weatherData WeatherResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return
	}

	// Afficher la température
	fmt.Printf("La température actuelle à %s est de %.2f degrés Kelvin.\n", city, weatherData.Main.Temperature)
}
