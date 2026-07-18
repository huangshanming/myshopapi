package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

func main() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJteW1hbGwtdXNlci1rZXkiLCJ1c2VyX2lkIjo4LCJyb2xlIjoidXNlciIsImlzcyI6Im15bWFsbCIsInN1YiI6IjgiLCJleHAiOjE3ODQ0NjQzMzUsImlhdCI6MTc4NDM3NzkzNX0.1iQDCZDy-VLAVbF9SkR4YZdhLlrtMH91je1d3mQ4M5Q"
	reqBody := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"product_id": 12,
				"quantity":   1,
				"sku_id":     28,
			},
		},
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(100)
	// 本地无 APISIX 时直连 order-service；有网关可改回 http://localhost:9080/api/v1/orders
	url := "http://127.0.0.1:8883/api/v1/orders"

	for i := 0; i < 100; i++ {
		go func() {
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("X-User-Id", "8")
			req.Header.Set("X-User-Role", "user")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err) // 不要在 err!=nil 时再碰 resp
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("status:", resp.StatusCode)
			fmt.Println(string(body))
			wg.Done()
			defer resp.Body.Close()
		}()
	}
	wg.Wait()

}
