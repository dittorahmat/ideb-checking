package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// App holds the application dependencies
type App struct {
	DB     *gorm.DB
	Router *mux.Router
	Config *Config
	ReadFileFunc func(name string) ([]byte, error)
}

// Config holds the application configuration
type Config struct {
	InputJSONPath string
	DBPath        string
	ServerPort    string
}

// NewConfig creates a new Config struct from environment variables
func NewConfig() *Config {
	inputJSONPath := os.Getenv("INPUT_JSON_PATH")
	if inputJSONPath == "" {
		inputJSONPath = "memory-bank/input.json" // Default value
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "ideb.db" // Default value
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080" // Default value
	}

	return &Config{
		InputJSONPath: inputJSONPath,
		DBPath:        dbPath,
		ServerPort:    ":" + serverPort,
	}
}

// NewApp creates a new App struct
func NewApp() (*App, error) {
	config := NewConfig()
	db, err := InitDatabase(config.DBPath)
	if err != nil {
		return nil, err
	}

	app := &App{
		DB:     db,
		Router: mux.NewRouter(),
		Config: config,
		ReadFileFunc: os.ReadFile,
	}

	app.RegisterRoutes()
	return app, nil
}

func (a *App) Run() {
	log.Printf("Server starting on port %s...", a.Config.ServerPort)
	if err := http.ListenAndServe(a.Config.ServerPort, a.Router); err != nil {
		log.Fatal(err)
	}
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	app.Run()
}