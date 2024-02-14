package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/hc", func(c echo.Context) error {
		containerName := os.Getenv("CONTAINER_NAME")
		return c.JSON(http.StatusOK, map[string]string{"message": containerName})
	})

	e.GET("/connect", func(c echo.Context) error {
		res, err := GetMessage()
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("ERROR: %s\n", err))
		}
		return c.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("connected === %s", *res),
		})
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

func GetMessage() (*string, error) {
	u, err := url.ParseRequestURI(fmt.Sprintf("http://%s:%s/hc", os.Getenv("SERVER_CONTAINER_HOST"), os.Getenv("SERVER_CONTAINER_PORT")))
	if err != nil {
		log.Errorf("Error: Parse URI === %s", err)
		return nil, err
	}

	var tmp struct {
		Message string `json:"message"`
	}
	resp, err := http.Get(fmt.Sprintf("%v", u))

	if err != nil {
		log.Errorf("Error: Http Request === %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &tmp); err != nil {
		log.Errorf("Error: Unmarshal === %s", err)
		return nil, err
	}

	return &tmp.Message, nil
}
