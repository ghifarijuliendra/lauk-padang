package models

// Kategori lauk padang
type Category struct {
	ID    int    `json:"id"`
	Nama  string `json:"name"`
	Harga int    `json:"price"`
	Stok  int    `json:"stock"`
}
