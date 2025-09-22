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
	InputJSONPath string `mapstructure:"INPUT_JSON_PATH"`
	DBPath        string `mapstructure:"DB_PATH"`
	ServerPort    string `mapstructure:"SERVER_PORT"`
	DisableAuth   bool   `mapstructure:"DISABLE_AUTH"`
}

// NewConfig creates a new Config struct from environment variables or config file
func NewConfig() (*Config, error) {
	// For prototype, always use defaults
	config := &Config{
		InputJSONPath: "../memory-bank/input.json",
		DBPath:        "ideb.db",
		ServerPort:    ":8080",
		DisableAuth:   true, // Disable authentication by default for prototype
	}

	log.Printf("Using default config: %+v", config)

	return config, nil
}

// NewApp creates a new App struct
func NewApp() (*App, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}

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
