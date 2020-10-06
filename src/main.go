package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	project := os.Getenv("PUB_PROJECT")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", func(c echo.Context) error {
		ctx := c.Request().Context()
		client, err := firestore.NewClient(ctx, project)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer func() {
			if client != nil {
				if err := client.Close(); err != nil {
					log.Print(err)
				}
			}
		}()

		m, err := unmarshal(c.Request().Body)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		data := string(m.Message.Data)

		ids := strings.Split(data, ":")
		log.Printf("[ids] %#+v", ids)
		log.Printf("[len(ids)] %d", len(ids))

		if len(ids) == 3 {
			// region
			ss, err := client.Collection("region").Doc(ids[1]).Get(ctx)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			log.Printf("[region] %#+v", ss.Data())
		} else if len(ids) == 4 {
			// school
			ss, err := client.Collection("region").Doc(ids[1]).Collection("school").Doc(ids[2]).Collection("operation").Doc(ids[3]).Get(ctx)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			log.Printf("[school] %#+v", ss.Data())
		} else {
			ss, err := client.Collection("region").Doc(ids[1]).Collection("school").Doc(ids[2]).Collection("operation").Doc(ids[4]).Get(ctx)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			log.Printf("[other] %#+v", ss.Data())
		}

		delay(ids[0])

		log.Printf("Message[ID:%s][Data:%s]", m.Message.ID, data)
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

func unmarshal(r io.ReadCloser) (*PubSubMessage, error) {
	var m PubSubMessage
	body, err := ioutil.ReadAll(r)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		return nil, err
	}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		return nil, err
	}
	return &m, nil
}

func delay(data string) {
	if strings.Contains(data, "region") {
		time.Sleep(5 * time.Second)
	}
	if strings.Contains(data, "school") {
		time.Sleep(5 * time.Second)
	}
	if strings.Contains(data, "grade") {
		time.Sleep(10 * time.Second)
	}
	if strings.Contains(data, "class") {
		time.Sleep(5 * time.Second)
	}
	if strings.Contains(data, "teacher") {
		time.Sleep(3 * time.Second)
	}
	if strings.Contains(data, "student") {
		time.Sleep(10 * time.Second)
	}
}
