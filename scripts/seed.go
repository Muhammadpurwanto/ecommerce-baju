package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var baseURL = "http://localhost:8080/api/v1"

func main() {
	log.Println("=== Starting Seed & Test ===")

	// 1. Register User
	email := fmt.Sprintf("tester_%d@example.com", time.Now().Unix())
	password := "rahasia123"
	registerPayload := map[string]string{
		"name":     "Tester Admin",
		"email":    email,
		"password": password,
	}
	sendRequest("POST", "/auth/register", registerPayload, "")

	// Tunggu sebentar agar RabbitMQ sempat memproses user.registered
	time.Sleep(2 * time.Second)

	// 2. Ubah role user menjadi 'admin' via koneksi database langsung
	log.Println("Updating user role to admin via database...")
	dsn := "root:secretpassword@tcp(127.0.0.1:3306)/user_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to user_db: %v", err)
	}
	result := db.Exec("UPDATE users SET role = 'admin' WHERE email = ?", email)
	if result.Error != nil {
		log.Fatalf("Failed to update role: %v", result.Error)
	}
	log.Println("User role updated successfully!")

	// 3. Login untuk mendapatkan Token Admin
	loginPayload := map[string]string{
		"email":    email,
		"password": password,
	}
	respMap := sendRequest("POST", "/auth/login", loginPayload, "")
	token := respMap["data"].(map[string]interface{})["access_token"].(string)
	log.Println("Got Admin Token:", token[:20]+"...")

	// 4. Test Profile & Address
	sendRequest("GET", "/users/profile", nil, token)

	addressPayload := map[string]interface{}{
		"label":       "Kantor",
		"recipient":   "Tester Admin",
		"phone":       "08123456789",
		"province":    "DKI Jakarta",
		"city":        "Jakarta Selatan",
		"district":    "Setiabudi",
		"postal_code": "12920",
		"detail":      "Gedung A, Lt 5",
		"is_default":  true,
	}
	sendRequest("POST", "/users/addresses", addressPayload, token)

	// 5. Create Category (Admin Only)
	catPayload := map[string]interface{}{
		"name":        fmt.Sprintf("Kategori %d", time.Now().Unix()),
		"description": "Koleksi pakaian pria terbaru",
	}
	catResp := sendRequest("POST", "/categories/", catPayload, token)
	catID := uint(catResp["data"].(map[string]interface{})["id"].(float64))

	// 6. Create Product (Admin Only)
	prodPayload := map[string]interface{}{
		"category_id": catID,
		"name":        fmt.Sprintf("Kemeja Flanel %d", time.Now().Unix()),
		"description": "Kemeja flanel bahan premium.",
		"brand":       "LocalPride",
		"gender":      "men",
		"base_price":  150000,
		"weight":      300,
		"is_active":   true,
	}
	prodResp := sendRequest("POST", "/products/", prodPayload, token)
	prodID := uint(prodResp["data"].(map[string]interface{})["id"].(float64))

	// 7. Create Variant
	varPayload := map[string]interface{}{
		"product_id":       prodID,
		"sku":              fmt.Sprintf("FLN-M-%d", time.Now().Unix()),
		"size":             "M",
		"color":            "Merah",
		"stock":            100,
		"price_adjustment": 0,
	}
	varRoute := fmt.Sprintf("/products/%d/variants", prodID)
	varResp := sendRequest("POST", varRoute, varPayload, token)
	varID := uint(varResp["data"].(map[string]interface{})["id"].(float64))

	// 8. Add to Cart
	cartPayload := map[string]interface{}{
		"product_id": prodID,
		"variant_id": varID,
		"quantity":   2,
	}
	sendRequest("POST", "/carts/items", cartPayload, token)
	sendRequest("GET", "/carts/", nil, token)

	// 9. Checkout / Create Order
	orderPayload := map[string]interface{}{
		"shipping_cost":    20000,
		"shipping_address": "Gedung A, Lt 5, Setiabudi, Jakarta Selatan",
		"courier":          "JNE",
		"notes":            "Kirim pagi",
		"items": []map[string]interface{}{
			{
				"product_id": prodID,
				"variant_id": varID,
				"quantity":   2,
				"price":      150000,
			},
		},
	}
	orderResp := sendRequest("POST", "/orders/", orderPayload, token)
	orderID := uint(orderResp["data"].(map[string]interface{})["id"].(float64))

	// 10. Pay Order
	payPayload := map[string]interface{}{
		"order_id": orderID,
		"amount":   320000, // 2 * 150000 + 20000
	}
	sendRequest("POST", "/payments/", payPayload, token)

	log.Println("=== All Tests & Seeds Completed Successfully ===")
}

func sendRequest(method, path string, payload interface{}, token string) map[string]interface{} {
	url := baseURL + path
	var reqBody io.Reader

	if payload != nil {
		jsonBytes, _ := json.Marshal(payload)
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("[%s] %s failed: %v", method, path, err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(bodyBytes, &result)

	success, ok := result["success"].(bool)
	if !ok || !success {
		log.Fatalf("[%s] %s returned error: %v", method, path, string(bodyBytes))
	}

	log.Printf("[%s] %s -> OK", method, path)
	return result
}
