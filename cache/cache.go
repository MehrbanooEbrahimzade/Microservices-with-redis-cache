package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/MehrbanooEbrahimzade/golang-redis-example/Data"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type JsonResponse struct {
	Data   []Data.Products `json:"data"`
	Source string          `json:"source"`
}

func GetProducts() (*JsonResponse, error) {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	cachedProducts, err := redisClient.Get("products").Bytes()

	var p1 []Data.Products
	json.Unmarshal(cachedProducts, &p1)
	fmt.Println("from cache : ", p1)

	response := JsonResponse{}

	if err != nil {

		dbProducts, err := fetchFromDb()

		if err != nil {
			return nil, err
		}

		cachedProducts, err = json.Marshal(dbProducts)

		if err != nil {
			return nil, err
		}

		err = redisClient.Set("products", cachedProducts, 20*time.Second).Err()

		if err != nil {
			return nil, err
		}

		response = JsonResponse{Data: dbProducts, Source: "PostgreSQL"}

		return &response, err
	}

	products := []Data.Products{}

	err = json.Unmarshal(cachedProducts, &products)

	if err != nil {
		return nil, err
	}

	response = JsonResponse{Data: products, Source: "Redis Cache"}

	return &response, nil
}

func fetchFromDb() ([]Data.Products, error) {
	fmt.Println("fetch From Db")
	dbUser := "postgres"
	dbName := "sample_company"

	conString := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)

	db, err := sql.Open("postgres", conString)

	if err != nil {
		return nil, err
	}

	queryString := `select
                     product_id,
                     product_name,
                     retail_price
                 from products`

	rows, err := db.Query(queryString)

	if err != nil {
		return nil, err
	}

	var records []Data.Products

	for rows.Next() {

		var p Data.Products

		err = rows.Scan(&p.ProductId, &p.ProductName, &p.RetailPrice)

		records = append(records, p)

		if err != nil {
			return nil, err
		}

	}

	return records, nil
}
