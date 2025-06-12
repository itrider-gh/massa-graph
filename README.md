# Massa Graph

Massa Graph is a visualization tool for displaying the Massa blockchain network topology using OpenStreetMap and Leaflet.js.

It consists of:
- A Go script to recursively retrieve the list of Massa nodes and their connections.
- A frontend (index.html) that visually displays these nodes and connections on a world map.

## 🧠 Features

- Recursively maps Massa nodes and their connections using `massa_map_recursive_links.go`
- Visual representation of the network with Leaflet.js and OpenStreetMap
- Lightweight, easy to deploy

---

## ⚙️ Backend – Node Mapping

### File: `massa_map_recursive_links.go`

This Go script fetches a list of Massa nodes by recursively crawling the network. It builds a map of active nodes and their connection relationships, outputting the data in a format compatible with the frontend.

### How to Run

Make sure you have Go installed (`go version`).

1. Clone the repository:
   ```bash
   https://github.com/itrider-gh/massa-graph.git
   cd massa-graph
    ````

2. Build and run the Go script:

   ```bash
   go run massa_map_recursive_links.go
   ```

This will generate or serve the data that the frontend can use to display the graph.

---

## 🌍 Frontend – Visualization

### File: `index.html`

This file uses **Leaflet.js** and **OpenStreetMap** to render Massa nodes on an interactive map. Each node is plotted based on geolocation (if available), and connections between nodes are drawn as lines.

### How to View

You can open the frontend directly in your browser:

```bash
xdg-open index.html   # Linux
open index.html       # macOS
start index.html      # Windows
```

Or serve it using any static HTTP server (like Python’s built-in one):

```bash
python3 -m http.server
```

---

## 📦 Dependencies

* [Go](https://golang.org/)
* [Leaflet.js](https://leafletjs.com/)
* [OpenStreetMap](https://www.openstreetmap.org/)

---

## ✨ Author

Made by [itrider](https://github.com/itrider-gh) – contributions welcome!

