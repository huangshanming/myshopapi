package main

import "fmt"
import "mymall/routes"

func main() {
	router := routes.InitRouter()
	router.Run(":8087")
	fmt.Println("Hello World")
}
