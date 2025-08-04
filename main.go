// @title Subscription API
// @version 1.0
// @description API для управления подписками

// @contact.name API Support
// @contact.email amet.kemal0032@gmail.com

// @host localhost:8080
// @BasePath /
package main

import (
	"Effective_Mobile_service/internal"
	"io"
	"net/http"
	"text/template"

	_ "Effective_Mobile_service/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data) //заполняем указанный html документ предоставленными данными
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("./web/templates/*.html")), //Говорим где искать файлы с расширением .html
	}
	router := echo.New()

	router.GET("/swagger/*", echoSwagger.WrapHandler)

	router.Renderer = t
	router.GET("/", internal.List)
	router.POST("/create-subscription", internal.CreateSubscription)
	router.PUT("/update-subscription/:id", internal.UpdateSubscription)
	router.DELETE("/delete-subscription/:id", internal.DeleteSubscription)
	router.GET("/subscription/:id", internal.ReadSubscription)

	router.GET("/subscriptions/calculator", func(c echo.Context) error {
		return c.Render(http.StatusOK, "subscription_sum_form", nil)
	})
	router.GET("/subscriptions/total", internal.CalculateSubscriptionsSum)

	router.Logger.Fatal(router.Start(":8080"))
}
