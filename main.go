package main

import (
	"github.com/labstack/echo"
	"github.com/tkanos/gonfig"
	"manager/conf"
	"manager/endpoint"
	"os"
)

var configuration = conf.Configuration{}

func init() {
	err := gonfig.GetConf("./conf/conf.json", &configuration)

	if err != nil {
		panic(err)
	}
}

//main contains all API endpoints
func main() {
	e := echo.New()

	//Create
	e.GET("/api/create", endpoint.Create)

	//Find
	e.GET("/api/find", endpoint.Find)

	//Store
	e.GET("/api/store", endpoint.Store)

	//Transaction
	e.GET("/api/transaction", endpoint.Transaction)

	//Next
	e.GET("/api/next", endpoint.Next)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
