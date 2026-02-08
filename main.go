package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"store-api-go/internal/database"
	"store-api-go/internal/handlers"
	"store-api-go/internal/models"
	"store-api-go/internal/repositories"
	"store-api-go/internal/services"
	"strings"

	"github.com/spf13/viper"
)

// config
type Config struct {
	Port          string `mapstructure:"PORT"`
	DBConn        string `mapstructure:"DB_CONN"`
	DBMaxOpenConn string `mapstructure:"DB_MAX_OPEN_CONNECTION"`
	BaseURL       string `mapstructure:"BASE_URL"`
}

// data
var categories = []models.Category{
	{ID: 1, Name: "Animal", Description: " A living thing that moves around to find food and eats plants or other animals for energy."},
	{ID: 2, Name: "Plant", Description: "A living thing that has leaves and roots that usually grow in the ground."},
	{ID: 3, Name: "Bacteria", Description: "Tiny, single-celled organisms with no nucleus."},
}

// main func
func main() {
	// Load the env
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:          viper.GetString("PORT"),
		DBConn:        viper.GetString("DB_CONN"),
		DBMaxOpenConn: viper.GetString("DB_MAX_OPEN_CONNECTION"),
	}

	// DB setup
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Define the layers
	categoryRepo := repositories.NewCategoryRepo(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productRepo := repositories.NewProductRepo(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	transactionRepo := repositories.NewTransactionRepo(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Setup routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Server is running",
		})
	})

	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)

	http.HandleFunc("/api/report/hari-ini", transactionHandler.HandleReportToday)

	// Serve the api
	address := config.BaseURL + ":" + config.Port
	fmt.Println("Server running on ", address)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Print(err)
	}
}
