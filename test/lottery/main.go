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
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJteW1hbGwtdXNlci1rZXkiLCJ1c2VyX2lkIjoxLCJyb2xlIjoidXNlciIsImlzcyI6Im15bWFsbCIsInN1YiI6IjEiLCJleHAiOjE3ODQ1MTAyNTAsImlhdCI6MTc4NDQyMzg1MH0.rHQurDZPj73daantcOhCmN7yzDTYa47SRYkw9pM-0RE"
	reqBody := map[string]interface{}{}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}
	num := 1000
	wg := sync.WaitGroup{}
	wg.Add(num)
	// 本地无 APISIX 时直连 order-service；有网关可改回 http://localhost:9080/api/v1/orders
	url := "http://127.0.0.1:5175/api/v1/lottery/draw"

	for i := 0; i < num; i++ {
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
