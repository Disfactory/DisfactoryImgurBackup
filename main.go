package main

// Hello World http server
import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

func ConvertToLogLevel(level string) log.Level {
	switch level {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	default:
		return log.DebugLevel
	}
}

type Server struct {
}

func GetCurrentDir() string {
	path, err := os.Executable()
	if err != nil {
		log.Error(err)
	}

	return filepath.Dir(path)
}

func (server *Server) initConfig() {
	// Load configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	defaultConfig := map[string]interface{}{
		"port": 1323,
		"log": map[string]interface{}{
			"rootDir":    filepath.Join(GetCurrentDir(), "logs"),
			"level":      "debug",
			"name":       "server.log",
			"maxSize":    36,
			"maxBackups": 30,
			"maxAge":     180,
		},
		"db": map[string]interface{}{
			"host":     "localhost",
			"port":     5432,
			"user":     "postgres",
			"password": "postgres",
			"dbname":   "postgres",
		},
	}
	viper.SetDefault("config", defaultConfig)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Can't read configuration, %s", err))
	}
}

func (server *Server) SetUp() {
	server.initConfig()

	logDir := viper.GetString("config.log.rootDir")
	logPath := filepath.Join(logDir, viper.GetString("config.log.name"))

	// Set up log
	log.SetFormatter(&log.JSONFormatter{})
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(ConvertToLogLevel(viper.GetString("config.log.level")))
	fileLogWriter := lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    viper.GetInt("config.log.maxSize"), // megabytes
		MaxBackups: viper.GetInt("config.log.maxBackups"),
		MaxAge:     viper.GetInt("config.log.maxAge"), // days
		Compress:   true,                              // disabled by default
	}
	multipleWriter := io.MultiWriter(os.Stdout, &fileLogWriter)
	log.SetOutput(multipleWriter)
}

func (server *Server) Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	log.Infof("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	server := Server{}
	server.SetUp()
	server.Start()
}
