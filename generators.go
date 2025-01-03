package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0",
	}

	acceptHeaders = []string{
		"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"application/json, text/javascript, */*;q=0.01",
	}

	languages = []string{
		"en-US,en;q=0.9", "fr-FR,fr;q=0.9", "de-DE,de;q=0.9",
		"es-ES,es;q=0.9", "ru-RU,ru;q=0.9", "ja-JP,ja;q=0.9",
	}
)

func getRandomItem(list []string) string {
	return list[rand.Intn(len(list))]
}

func getRandomUserAgent() string {
	return getRandomItem(userAgents)
}

func getRandomAccept() string {
	return getRandomItem(acceptHeaders)
}

func getRandomLanguage() string {
	return getRandomItem(languages)
}

func generateIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func init() {
	rand.Seed(time.Now().UnixNano()) // Seed random
}
