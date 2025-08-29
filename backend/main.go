package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
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
}

// NewConfig creates a new Config struct from environment variables or config file
func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName(".env") // name of config file (without extension)
	v.SetConfigType("dotenv")    // type of the config file
	v.AddConfigPath(".")       // path to look for the config file in the current directory
	v.AutomaticEnv()           // read in environment variables that match

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error and use defaults/env vars
			log.Println("Config file not found, using environment variables or defaults.")
		} else {
			return nil, fmt.Errorf("Fatal error reading config file: %w", err)
		}
	}

	// Set defaults if not found in config file or environment variables
	v.SetDefault("INPUT_JSON_PATH", "../memory-bank/input.json")
	v.SetDefault("DB_PATH", "ideb.db")
	v.SetDefault("SERVER_PORT", "8080")

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal config: %w", err)
	}

	config.ServerPort = ":" + config.ServerPort // Prepend colon for net/http.ListenAndServe

	return &config, nil
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
