package main

import (
	"fmt"
	"math/rand"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (compatible, MSIE 11, Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.188",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/116.0",
	"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0",
}

var acceptHeaders = []string{
	"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
	"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8,en-US;q=0.5",
	"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8,en;q=0.7",
}

var languages = []string{
	"en-US,en;q=0.9",
	"en-GB,en;q=0.9",
	"en-CA,en;q=0.9",
	"en-AU,en;q=0.9",
	"en-NZ,en;q=0.9",
	"en-ZA,en;q=0.9",
	"en-IE,en;q=0.9",
	"en-IN,en;q=0.9",
	"fr-FR,fr;q=0.9",
	"de-DE,de;q=0.9",
	"es-ES,es;q=0.9",
	"it-IT,it;q=0.9",
	"ru-RU,ru;q=0.9",
	"ja-JP,ja;q=0.9",
	"zh-CN,zh;q=0.9",
	"pt-BR,pt;q=0.9",
	"pl-PL,pl;q=0.9",
	"nl-NL,nl;q=0.9",
}

func getRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func getRandomAccept() string {
	return acceptHeaders[rand.Intn(len(acceptHeaders))]
}

func getRandomLanguage() string {
	return languages[rand.Intn(len(languages))]
}

func generateIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256))
}
