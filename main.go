package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Kategori lauk padang
type Category struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

// In-memory storage
var categories = []Category{
	{ID: 1, Nama: "Telur Dadar", Harga: 10000, Stok: 15},
	{ID: 2, Nama: "Tunjang", Harga: 20000, Stok: 20},
	{ID: 3, Nama: "Dendeng Merah", Harga: 18000, Stok: 10},
}

func main() {
	// Handler 1: Cek Status
	http.HandleFunc("/laukpadang", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// Handler 2: CRUD Lauk (GET, POST, PUT, DELETE)
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

		// get && post
		if idStr == "" {
			// GET localhost:9090/api/categories
			if r.Method == "GET" {
				json.NewEncoder(w).Encode(categories)
				return
			}

			// POST localhost:9090/api/categories
			if r.Method == "POST" {
				var kategoriBaru Category
				if err := json.NewDecoder(r.Body).Decode(&kategoriBaru); err != nil {
					http.Error(w, "Invalid request", http.StatusBadRequest)
					return
				}

				kategoriBaru.ID = len(categories) + 1
				categories = append(categories, kategoriBaru)

				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(kategoriBaru)
				return
			}
		}

		// get by id && put && delete
		if idStr != "" {
			// convert id ke int
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "ID harus angka", http.StatusBadRequest)
				return
			}

			// GET localhost:9090/api/categories/{id}
			if r.Method == "GET" {
				for _, c := range categories {
					if c.ID == id {
						json.NewEncoder(w).Encode(c)
						return
					}
				}
				http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
				return
			}

			// PUT localhost:9090/api/categories/{id}
			if r.Method == "PUT" {
				var updateData Category
				if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
					http.Error(w, "Invalid request body", http.StatusBadRequest)
					return
				}

				for i := range categories {
					if categories[i].ID == id {
						updateData.ID = id
						categories[i] = updateData
						json.NewEncoder(w).Encode(updateData)
						return
					}
				}
				http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
				return
			}

			// DELETE localhost:9090/api/categories/{id}
			if r.Method == "DELETE" {
				for i, c := range categories {
					if c.ID == id {
						// Teknik menghapus element dari Slice:
						// Gabungkan data sebelum index 'i' dengan data sesudah index 'i'
						categories = append(categories[:i], categories[i+1:]...)

						json.NewEncoder(w).Encode(map[string]string{
							"message": "Sukses delete data",
						})
						return
					}
				}
				http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
				return
			}
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Cek apakah ada environment variable PORT (dari Railway)
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090" // Default port kalau jalan di laptop (lokal)
	}

	fmt.Println("Server running di port " + port)

	// Ganti ":9090" jadi variabel address
	addr := ":" + port
	http.ListenAndServe(addr, nil)
}
