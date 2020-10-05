package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", func(c echo.Context) error {
		var m PubSubMessage
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			log.Printf("ioutil.ReadAll: %v", err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		if err := json.Unmarshal(body, &m); err != nil {
			log.Printf("json.Unmarshal: %v", err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		log.Printf("Message.ID %s", m.Message.ID)
		log.Printf("Message.Data %s", string(m.Message.Data))
		return c.String(http.StatusOK, "OK")
	})

	e.Logger.Fatal(e.Start(":8080"))
}

type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}
