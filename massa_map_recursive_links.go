package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type NodeInfo struct {
	IP          string   `json:"ip"`
	Version     string   `json:"version"`
	Connections []string `json:"connections"`
	Lat         float64  `json:"lat"`
	Lon         float64  `json:"lon"`
}

type UnreachableInfo struct {
	IP      string  `json:"ip"`
	Version string  `json:"version,omitempty"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type NodeStatusResponse struct {
	Result struct {
		Version        string                              `json:"version"`
		ConnectedNodes map[string][2]interface{}           `json:"connected_nodes"`
	} `json:"result"`
}

var client = &http.Client{Timeout: 3 * time.Second}
const massaURL = "https://mainnet.massa.net/api/v2"
const defaultPort = "33035"

var visited sync.Map
var graph sync.Map
var unreachable sync.Map
var wg sync.WaitGroup

func main() {
	fmt.Println("Fetching initial node status...")

	// Timeout global (ex. : 3 minutes)
	go func() {
		time.Sleep(1 * time.Minute)
		fmt.Println("⏱️ Timeout global atteint, export forcé…")
		exportJSONGraph()
		os.Exit(0)
	}()

	initialIPs := fetchConnectedIPs(massaURL)

	for _, ip := range initialIPs {
		wg.Add(1)
		go queryNode(ip)
	}

	wg.Wait()
	exportJSONGraph()
}


func fetchConnectedIPs(url string) []string {
	reqBody := []byte(`{"jsonrpc":"2.0","method":"get_status","params":[],"id":1}`)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result NodeStatusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("JSON parse error:", err)
		return nil
	}

	var ips []string
	for _, node := range result.Result.ConnectedNodes {
		rawIP := fmt.Sprintf("%v", node[0])
		cleanIP := strings.TrimPrefix(rawIP, "::ffff:")
		ips = append(ips, cleanIP)
	}
	return ips
}

func queryNode(ip string) {
	defer wg.Done()

	if _, loaded := visited.LoadOrStore(ip, true); loaded {
		return
	}

	url := fmt.Sprintf("http://%s:%s", ip, defaultPort)
	fmt.Printf("➡️  Querying %s\n", url)

	reqBody := []byte(`{"jsonrpc":"2.0","method":"get_status","params":[],"id":1}`)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("❌ %s: %v\n", ip, err)
		lat, lon := getGeo(ip)
		unreachable.Store(ip, &UnreachableInfo{IP: ip, Lat: lat, Lon: lon})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result NodeStatusResponse
	if err := json.Unmarshal(body, &result); err == nil {
		fmt.Printf("✅ %s responded, adding %d new nodes...\n", ip, len(result.Result.ConnectedNodes))
		lat, lon := getGeo(ip)
		conn := []string{}
		for _, node := range result.Result.ConnectedNodes {
			rawIP := fmt.Sprintf("%v", node[0])
			cleanIP := strings.TrimPrefix(rawIP, "::ffff:")
			conn = append(conn, cleanIP)
			wg.Add(1)
			go queryNode(cleanIP)
		}
		graph.Store(ip, &NodeInfo{
			IP:          ip,
			Version:     result.Result.Version,
			Connections: conn,
			Lat:         lat,
			Lon:         lon,
		})
	} else {
		fmt.Printf("⚠️  %s: Failed to parse JSON\n", ip)
		lat, lon := getGeo(ip)
		unreachable.Store(ip, &UnreachableInfo{IP: ip, Lat: lat, Lon: lon})
	}
}

func getGeo(ip string) (float64, float64) {
	resp, err := http.Get("https://ipwho.is/" + ip)
	if err != nil {
		return 0, 0
	}
	defer resp.Body.Close()
	var result struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Success   bool    `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.Success {
		return 0, 0
	}
	return result.Latitude, result.Longitude
}

func exportJSONGraph() {
	type Link struct {
		Source string `json:"source"`
		Target string `json:"target"`
		Out    bool   `json:"out"`
	}
	type Export struct {
		Nodes []NodeInfo        `json:"nodes"`
		Dead  []UnreachableInfo `json:"unreachable"`
		Links []Link            `json:"links"`
	}
	export := Export{}

	graph.Range(func(key, value any) bool {
		node := value.(*NodeInfo)
		export.Nodes = append(export.Nodes, *node)
		for _, conn := range node.Connections {
			export.Links = append(export.Links, Link{
				Source: node.IP,
				Target: conn,
				Out:    true,
			})
			// Add reverse link if the other node is known
			if v, ok := graph.Load(conn); ok {
				other := v.(*NodeInfo)
				for _, back := range other.Connections {
					if back == node.IP {
						export.Links = append(export.Links, Link{
							Source: conn,
							Target: node.IP,
							Out:    false,
						})
						break
					}
				}
			}
		}
		return true
	})

	unreachable.Range(func(key, value any) bool {
		node := value.(*UnreachableInfo)
		export.Dead = append(export.Dead, *node)
		return true
	})

	f, _ := os.Create("graph.json")
	defer f.Close()
	json.NewEncoder(f).Encode(export)
	fmt.Println("🌍 JSON graph exported to graph.json (with unreachable nodes and link direction)")
}
