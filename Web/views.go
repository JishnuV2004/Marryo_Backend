package web

import (
	"github.com/gofiber/template/html/v2"
)

func InitViews() *html.Engine {
	engine := html.New("./Web/Templates", ".html")
	return engine
}

