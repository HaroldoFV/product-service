package main

import (
	"database/sql"
	"fmt"
	"github.com/HaroldoFV/product-service/configs"
	_ "github.com/HaroldoFV/product-service/docs"
	"github.com/HaroldoFV/product-service/internal/infra/database"
	"github.com/HaroldoFV/product-service/internal/infra/web"
	"github.com/HaroldoFV/product-service/internal/infra/web/webserver"
	"github.com/HaroldoFV/product-service/internal/usecase"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"path/filepath"
)

// @title Product Service API
// @version 1.0
// @description This is a product microservice API.
// @host localhost:8000
// @BasePath /api/v1
func main() {
	dir, _ := os.Getwd()
	fmt.Println("Diretório atual:", dir)

	config, err := configs.LoadConfig(dir)
	if err != nil {
		rootDir := filepath.Join(dir, "..", "..")
		config, err = configs.LoadConfig(rootDir)
		if err != nil {
			fmt.Println("Erro ao carregar configurações:", err)
			panic(err)
		}
	}
	fmt.Printf("Configurações carregadas: %+v\n", config)

	db, err := sql.Open(config.DBDriver, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	webServer := webserver.NewWebServer(":" + config.WebServerPort)

	productRepository := database.NewProductRepository(db)
	createProductUseCase := usecase.NewCreateProductUseCase(productRepository)
	webProductHandler := web.NewWebProductHandler(createProductUseCase, productRepository)

	webServer.AddHandler(http.MethodPost, "/products", webProductHandler.Create)
	webServer.AddHandler(http.MethodGet, "/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+config.WebServerPort+"/swagger/doc.json"),
	))

	fmt.Println("Starting web server on port", config.WebServerPort)
	go func() {
		err = webServer.Start()
		if err != nil {
			panic(err)
		}
	}()
	select {}
}
