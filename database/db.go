package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
	"time"
)

type ProductData struct {
	mu sync.Mutex
	db *sql.DB
}

type Product struct {
	ID     int
	Time   time.Time
	SKU    string
	URL    string
	TITLE  string
	PRICE  float64
	PARAMS string
}

const create string = `
  CREATE TABLE IF NOT EXISTS product_data (
  	  id INTEGER NOT NULL PRIMARY KEY,
	  time DATETIME NOT NULL,
	  sku VARCHAR(255) NOT NULL,
      url TEXT,
      title TEXT,
      price DECIMAL(10,2) NOT NULL,
      params JSONB
  );`

const file string = "productData.db"

// --- TABLE CREATION ---
func CreateTable() (*ProductData, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("createTable error: %v", err)
	}
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	return &ProductData{
		db: db,
	}, nil
}

// --- CREATE ---
func (pd *ProductData) CreateProduct(ProductData *Product) int {
	query := `INSERT INTO product_data (time, sku, url, title, price, params) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := pd.db.Exec(
		query,
		time.Now(),
		ProductData.SKU,
		ProductData.URL,
		ProductData.TITLE,
		ProductData.PRICE,
		ProductData.PARAMS,
	)
	if err != nil {
		log.Fatalf("createProduct error: %v", err)
	}
	id, _ := result.LastInsertId()
	return int(id)
}

// --- READ ---
func (pd *ProductData) GetProduct(sku string) Product {
	var p Product
	query := `SELECT * FROM product_data WHERE sku LIKE ?`
	row := pd.db.QueryRow(query, sku)
	err := row.Scan(
		&p.ID,
		&p.Time,
		&p.SKU,
		&p.URL,
		&p.TITLE,
		&p.PRICE,
		&p.PARAMS)
	if err != nil {
		log.Fatalf("getProduct error: %v", err)
	}
	return p
}

// --- UPDATE ---
func (pd *ProductData) UpdateProduct(ProductData *Product) {
	query := `UPDATE product_data SET sku = ? WHERE id = ?`
	_, err := pd.db.Exec(query, ProductData.SKU, ProductData.ID)
	if err != nil {
		log.Fatalf("updateProduct error: %v", err)
	}
}

// --- DELETE ---
func (pd *ProductData) DeleteProduct(id int) {
	query := `DELETE FROM product_data WHERE id = ?`
	_, err := pd.db.Exec(query, id)
	if err != nil {
		log.Fatalf("deleteProduct error: %v", err)
	}
}
