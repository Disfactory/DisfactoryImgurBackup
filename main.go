package main

// Hello World http server
import (
	"disfactory/imgur-backup/adapter"
	"disfactory/imgur-backup/domain"
	"disfactory/imgur-backup/service"
	"disfactory/imgur-backup/utils"

	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Server struct {
	backupRepository       domain.BackupRepository
	factoryImageRepository domain.FactoryImageRepository
	imageDownloader        domain.ImageDownloader

	backupImgurService *service.BackupImgurService
	httpService        *service.HttpService
}

func GetCurrentDir() string {
	path, err := os.Executable()
	if err != nil {
		utils.Logger().Error(err)
	}

	return filepath.Dir(path)
}

func generateDefaultConfig() map[string]interface{} {
	return map[string]interface{}{
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
			"port":     "5432",
			"user":     "postgres",
			"password": "postgres",
			"dbname":   "postgres",
		},
		"backup": map[string]interface{}{
			"rootDir": filepath.Join(GetCurrentDir(), "backup"),
			// Run cronjob at 00:00:00 everyday
			"cronjob": "0 0 0 * * *",
		},
		"http": map[string]interface{}{
			"address": "127.0.0.1:8000",
		},
	}
}

func (server *Server) initConfig() {
	// Load configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("config", generateDefaultConfig())

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Can't read configuration, %s", err))
	}
}

func (server *Server) SetUp() {
	server.initConfig()

	logDir := viper.GetString("config.log.rootDir")
	logPath := filepath.Join(logDir, viper.GetString("config.log.name"))

	utils.Init()
	logger := utils.Logger()
	logger.SetLevel(utils.ConvertToLogLevel(viper.GetString("config.log.level")))
	fileLogWriter := lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    viper.GetInt("config.log.maxSize"), // megabytes
		MaxBackups: viper.GetInt("config.log.maxBackups"),
		MaxAge:     viper.GetInt("config.log.maxAge"), // days
		Compress:   true,                              // disabled by default
	}
	multipleWriter := io.MultiWriter(os.Stdout, &fileLogWriter)
	logger.SetOutput(multipleWriter)

	// Init backup repository
	backupDir := viper.GetString("config.backup.rootDir")
	logger.Infof("Backup dir: %s", backupDir)
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err := os.MkdirAll(backupDir, os.ModePerm)
		if err != nil {
			logger.Error(err)
		}
	}
	server.backupRepository = adapter.NewLocalFileBackupRepository(backupDir)

	// Init factory image repository
	pgParams := adapter.PGParameters{
		Host:     viper.GetString("config.db.host"),
		Port:     viper.GetString("config.db.port"),
		Account:  viper.GetString("config.db.user"),
		Password: viper.GetString("config.db.password"),
		DBName:   viper.GetString("config.db.dbname"),
	}
	dbFactoryImageRepository, err := adapter.NewDBFactoryImageRepository(pgParams)
	if err != nil {
		logger.Panic(err)
	}
	server.factoryImageRepository = dbFactoryImageRepository

	// Init image downloader
	server.imageDownloader = adapter.NewHttpImageDownloader()

	// Init backup imgur service
	server.backupImgurService = service.NewBackupImgurService(
		server.factoryImageRepository,
		server.backupRepository,
		server.imageDownloader,
		viper.GetString("config.backup.cronjob"),
	)

	// Init http service
	server.httpService = service.NewHttpService(
		viper.GetString("config.http.address"),
		backupDir,
	)
}

func (server *Server) Test() {
	// Try to get some images from FactoryImageRepository
	factoryImages, err := server.factoryImageRepository.GetImages(100, 0)
	if err != nil {
		utils.Logger().Error(err)
	} else {
		utils.Logger().Infof("Get %d images from FactoryImageRepository", len(factoryImages))
	}
}

func (server *Server) Start() {
	server.backupImgurService.Start()
	server.httpService.Start()
}

func CreateDefaultConfigFile() {
	fmt.Println("Generting config file")
	// Get current directory and create config file
	configFilePath := filepath.Join(GetCurrentDir(), "config.yaml")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		file, err := os.Create(configFilePath)
		if err != nil {
			utils.Logger().Error(err)
		}
		defer file.Close()

		defaultConfig := generateDefaultConfig()
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.Set("config", defaultConfig)
		viper.WriteConfigAs(configFilePath)
	} else {
		utils.Logger().Warnf("Config file %s already exists", configFilePath)
	}
}

func CheckConfig() {
	server := Server{}
	server.SetUp()
	server.Test()
}

func main() {
	// action: start, init
	action := flag.String("action", "", "Action to do")
	flag.Parse()
	if *action == "init" {
		CreateDefaultConfigFile()
		return
	} else if *action == "start" {
		server := Server{}
		server.SetUp()
		server.Start()
	} else if *action == "check" {
		CheckConfig()
	} else {
		fmt.Println("Invalid action")
	}
}
