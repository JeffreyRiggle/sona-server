package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
)

var logManager LogManager

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path = ""
	}
	initialize(path)
	defer incidentManager.CleanUp()

	router := NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func initialize(file string) {
	if len(file) <= 0 {
		log.Println("No config provided using default settings")
		useDefaultConfig()
		return
	}

	log.Printf("Loading config from %v\n", file)
	loadConfig(file)
}

func loadConfig(file string) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()

	if err != nil {
		log.Println("Unable to load config file using defaults")
		useDefaultConfig()
		return
	}

	parser := json.NewDecoder(configFile)
	parser.Decode(&config)

	log.Printf("Loaded config %v\n", config)
	setupLogManager(config)
	setupFileManager(config)
	setupIncidentManager(config)
	setupUsermanager(config)

	hookManager = HookManager{
		config.Hooks.AddedHooks,
		config.Hooks.UpdatedHooks,
		config.Hooks.AttachedHooks,
		config.Hooks.AddedUserHooks,
		config.Hooks.UpdatedUserHooks,
	}
}

func useDefaultConfig() {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal("Unable to find current user")
		panic(err)
	}

	logManager = LogManager{currentUser.HomeDir, true}
	logManager.Initialize()
	fileManager = LocalFileManager{currentUser.HomeDir}
	incidentManager = RuntimeIncidentManager{make(map[int64]*Incident), make(map[int][]Attachment)}
	userManager = RuntimeUserManager{make(map[int]*User), make(map[int]string), make(map[int][]string), make([]string, 1)}
}

func setupLogManager(config Config) {
	logManager = LogManager{config.Logging.Path, config.Logging.Enabled}
	logManager.Initialize()
}

func setupFileManager(config Config) {
	if config.FileManagerType > 0 {
		log.Printf("Setting file manager to s3 with region %v and bucket %v\n", config.S3Config.Region, config.S3Config.Bucket)
		s3Manager := S3FileManager{config.S3Config.Region, config.S3Config.Bucket}
		s3Manager.Initialize()
		fileManager = &s3Manager
		return
	}

	if len(config.LocalFileConfig.Path) > 0 {
		log.Printf("Setting local file manager path to %v\n", config.LocalFileConfig.Path)
		fileManager = LocalFileManager{config.LocalFileConfig.Path}
		return
	}

	currentUser, err := user.Current()
	if err != nil {
		log.Fatal("Unable to find current user")
		panic(err)
	}

	fileManager = LocalFileManager{currentUser.HomeDir}
}

func setupIncidentManager(config Config) {
	if config.IncidentManagerType == 0 {
		log.Println("Using Runtime incident manager")
		incidentManager = RuntimeIncidentManager{make(map[int64]*Incident), make(map[int][]Attachment)}
		return
	}

	if config.IncidentManagerType == 1 {
		setupDynamoDBIncidentManager(config)
		return
	}

	if config.IncidentManagerType == 2 {
		setupMySQLIncidentManager(config)
		return
	}

	if config.IncidentManagerType == 3 {
		setupDataStoreIncidentManager(config)
		return
	}

	panic(fmt.Sprintf("Invalid incident manager config %v", config.IncidentManagerType))
}

func setupDynamoDBIncidentManager(config Config) {
	if len(config.DynamoConfig.Region) <= 0 {
		panic("No configured region")
	}

	log.Printf("Using dynamodb incident manager with region %v\n", config.DynamoConfig.Region)

	var incs string
	if len(config.DynamoConfig.IncidentTableOverride) > 0 {
		incs = config.DynamoConfig.IncidentTableOverride
		log.Printf("Found incident table override %v\n", incs)
	} else {
		incs = "Incidents"
	}

	var attach string
	if len(config.DynamoConfig.AttachmentTableOverride) > 0 {
		attach = config.DynamoConfig.AttachmentTableOverride
		log.Printf("Found attachment table override %v\n", attach)
	} else {
		attach = "IncidentAttachments"
	}

	dbManager := DynamoDBIncidentManager{&config.DynamoConfig.Region, &incs, &attach}
	dbManager.Initialize()
	incidentManager = &dbManager
}

func setupMySQLIncidentManager(config Config) {
	conn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", config.MYSQL.UserName, config.MYSQL.Password, config.MYSQL.Host, config.MYSQL.Port, config.MYSQL.DBName)
	db, err := sql.Open("mysql", conn)

	if err != nil {
		panic(err)
	}

	mysqlManager := MySQLManager{db}
	mysqlManager.Initialize()
	incidentManager = &mysqlManager
}

func setupDataStoreIncidentManager(config Config) {
	context, client := CreateDataStoreClient(config.DataStore.ProjectName, config.DataStore.AuthFile)
	dataStoreManager := DataStoreIncidentManager{
		context,
		client,
		0,
	}
	incidentManager = &dataStoreManager
}

func setupUsermanager(config Config) {
	userManager = RuntimeUserManager{make(map[int]*User), make(map[int]string), make(map[int][]string), config.User.DefaultPermissions}
}
