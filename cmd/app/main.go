package main

import app "github.com/Tamplier2911/gorest/internal"

// @title Go REST API example
// @version 2.0
// @description This is a sample rest api realized in go language for education purposes.
//
// @contact.email artyom.nikolaev@syahoo.com
//
// @host localhost:8000
// @BasePath /api/v2
func main() {
	app := app.Monolith{}
	app.Setup()
	app.Start()
}
