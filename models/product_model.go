package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"mvc/config"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

const cacheExpiration = 1 * time.Microsecond

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	Description string  `json:"description"`
}
type ProductModel struct {
	db    *sql.DB
	cache *redis.Client
}

func NewProductModel(db config.MySQLConnection, cache config.RedisConnection) *ProductModel {
	dbConn, err := db.GetConnection()
	if err != nil {
		log.Fatalf("error conection to database: %s", err.Error())
	}
	cacheConn, err := cache.GetConnection()
	if err != nil {
		log.Fatalf("error conection to cache: %s", err.Error())
	}
	productModel := &ProductModel{db: dbConn, cache: cacheConn}
	return productModel
}

func (m *ProductModel) GetByID(id int64) (*Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)
	cacheResult, err := m.cache.Get(cacheKey).Result()
	if err == nil {
		var product Product
		err = json.Unmarshal([]byte(cacheResult), &product)
		if err == nil {
			return &product, nil
		}
	}

	query := "SELECT id, name, price, quantity, description FROM `my-mvc`.product_list WHERE id = ?"
	row := m.db.QueryRow(query, id)
	product := &Product{}
	err = row.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Description)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	} else if err != nil {
		return nil, fmt.Errorf("error scanning product row: %s", err.Error())
	}

	// Set the product in cache
	productJSON, err := json.Marshal(product)
	if err != nil {
		return nil, fmt.Errorf("error marshalling product to JSON: %s", err.Error())
	}
	err = m.cache.Set(cacheKey, string(productJSON), cacheExpiration).Err()
	if err != nil {
		return nil, fmt.Errorf("error setting product in cache: %s", err.Error())
	}

	return product, nil
}
func (m *ProductModel) GetAll() ([]*Product, error) {
	query := "SELECT id, name, price, quantity, description FROM `my-mvc`.product_list"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var products []*Product
	for rows.Next() {
		product := &Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (m *ProductModel) Create(product *Product) error {
	query := "INSERT INTO `my-mvc`.product_list (name, price, quantity, description) VALUES(?, ?, ?, ?)"
	result, err := m.db.Exec(query, product.Name, product.Price, product.Quantity, product.Description)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = int64(int(id))

	return nil
}
func (m *ProductModel) Update(product *Product) error {
	query := "UPDATE `my-mvc`.product_list SET name=?, price=?, quantity=?, description=? WHERE id=?"
	_, err := m.db.Exec(query, product.Name, product.Price, product.Quantity, product.Description, product.ID)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("product:%d", product.ID)
	jsonResult, err := json.Marshal(product)
	if err == nil {
		m.cache.Set(cacheKey, jsonResult, cacheExpiration*time.Minute)
	}

	return nil
}
func (m *ProductModel) Delete(id int) error {
	query := "DELETE FROM `my-mvc`.product_list WHERE id = ?"
	_, err := m.db.Exec(query, id)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("product:%d", id)
	m.cache.Del(cacheKey)

	return nil
}
