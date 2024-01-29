package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"io/ioutil"
	"net/http"

	"github.com/redis/go-redis/v9"
)

// Structure pour définir la réponse JSON de l'API OpenWeatherMap
type WeatherResponse struct {
	Main struct {
		Temperature float64 `json:"temp"`
	} `json:"main"`
}

// Fonction pour initialiser le client Redis
func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Adresse du serveur Redis
		Password: "",               // Pas de mot de passe
		DB:       0,                // Base de données par défaut
	})
	return client
}

// Fonction pour obtenir les données météorologiques depuis l'API avec mise en cache dans Redis
func getCachedWeatherData(client *redis.Client, city string) (float64, error) {
	// Tentative d'obtenir les données depuis le cache Redis
	cachedData, err := client.Get(context.TODO(), city).Result()
	if err == redis.Nil {
		// Données non trouvées dans Redis, besoin de faire une requête à l'API
		data, err := fetchWeatherData(city)
		if err != nil {
			return 0, err
		}

		// Analyse des données JSON pour extraire la température
		var weatherData WeatherResponse
		if err := json.Unmarshal([]byte(data), &weatherData); err != nil {
			return 0, err
		}

		// Mettez en cache les données avec un TTL de 300 secondes (5 minutes)
		client.Set(context.TODO(), city, weatherData.Main.Temperature, 300*time.Second)

		return weatherData.Main.Temperature, nil
	} else if err != nil {
		return 0, err
	}

	// Les données ont été trouvées dans le cache, retournez-les
	var temperature float64
	if _, err := fmt.Sscanf(cachedData, "%f", &temperature); err != nil {
		return 0, err
	}
	return temperature, nil
}

// Fonction pour faire la requête à l'API OpenWeatherMap
func fetchWeatherData(city string) (string, error) {
	// Remplacez "VOTRE_CLE_API" par votre clé API OpenWeatherMap
	apiKey := "68ff09cbdc1317630b3f955444948c1b"
	apiURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)

	// Effectuer la requête HTTP
	response, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Lire le corps de la réponse
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// Retourner les données de la réponse sous forme de chaîne de caractères
	return string(body), nil
}

func main() {
	// Initialisez le client Redis
	redisClient := newRedisClient()

	// Remplacez "VOTRE_CLE_API" par votre clé API OpenWeatherMap
	city := "Paris" // Remplacez par la ville de votre choix

	// Obtenez les données météorologiques avec mise en cache
	temperature, err := getCachedWeatherData(redisClient, city)
	if err != nil {
		log.Fatal(err)
	}

	// Affichez la température
	fmt.Printf("La température actuelle à %s est de %.2f degrés Celsius.\n", city, temperature)
}
