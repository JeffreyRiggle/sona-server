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
var admin = AddUser{
	EmailAddress: "a@b.c",
	FirstName:    "Admin",
	UserName:     "Administrator",
	Password:     "admin",
}
var adminPermissions = make([]string, 0)

func main() {
	adminPermissions = append(adminPermissions, availablePermissions.master)
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path = ""
	}

	config := processConfig(path)
	initialize(config)

	defer incidentManager.CleanUp()
	startListening(config)
}

func startListening(config Config) {
	router := NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	handler := handlers.CORS(originsOk, headersOk, methodsOk)(router)
	if len(config.Security.Certificate) <= 0 || len(config.Security.Key) <= 0 {
		log.Fatal(http.ListenAndServe(":8080", handler))
	} else {
		log.Fatal(http.ListenAndServeTLS(":8080", config.Security.Certificate, config.Security.Key, handler))
	}
}

func processConfig(file string) Config {
	if len(file) <= 0 {
		return createDefaultConfig()
	}

	log.Printf("Loading config from %v\n", file)
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()

	if err != nil {
		log.Println("Unable to load config file using defaults")
		return createDefaultConfig()
	}

	parser := json.NewDecoder(configFile)
	parser.Decode(&config)

	log.Printf("Loaded config %v\n", config)
	return config
}

func createDefaultConfig() Config {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal("Unable to find current user")
		panic(err)
	}

	return Config{
		Logging: LogConfig{
			Path:    currentUser.HomeDir,
			Enabled: true,
		},
		FileManagerType: 0,
		LocalFileConfig: LocalFileManagerConfig{
			Path: currentUser.HomeDir,
		},
		ManagerType: 0,
		Admin: AdminConfig{
			EmailAddress: "a@b.c",
			Password:     "admin",
		},
	}
}

func initialize(config Config) {
	setupAdmin(config)
	setupLogManager(config)
	setupFileManager(config)
	setupManagers(config)

	hookManager = HookManager{
		config.Hooks.AddedHooks,
		config.Hooks.UpdatedHooks,
		config.Hooks.AttachedHooks,
		config.Hooks.AddedUserHooks,
		config.Hooks.UpdatedUserHooks,
	}
}

func setupAdmin(config Config) {
	if len(config.Admin.EmailAddress) != 0 {
		admin.EmailAddress = config.Admin.EmailAddress
	}

	if len(config.Admin.Password) != 0 {
		admin.Password = config.Admin.Password
	}
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

func setupManagers(config Config) {
	if config.ManagerType == 0 {
		log.Println("Using Runtime managers")
		incidentManager = RuntimeIncidentManager{make(map[int64]*Incident), make(map[int][]Attachment)}
		setupRuntimeUsermanager(config)
		return
	}

	if config.ManagerType == 1 {
		setupDynamoDBManagers(config)
		return
	}

	if config.ManagerType == 2 {
		conn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", config.MYSQL.UserName, config.MYSQL.Password, config.MYSQL.Host, config.MYSQL.Port, config.MYSQL.DBName)
		db, err := sql.Open("mysql", conn)

		if err != nil {
			panic(err)
		}

		setupMySQLIncidentManager(config, db)
		setupSQLUsermanager(config, db)
		return
	}

	if config.ManagerType == 3 {
		setupDataStoreIncidentManager(config)
		return
	}

	panic(fmt.Sprintf("Invalid manager config %v", config.ManagerType))
}

func ensureAdminAccount(startIndex int64) {
	_, found := userManager.GetUser(startIndex)

	if !found {
		log.Println("No administrator found creating admin account")
		_, res := userManager.AddUser(&admin)
		userManager.SetPermissions(res.Id, adminPermissions)
	}
}

func setupDynamoDBManagers(config Config) {
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

	var usr string
	if len(config.DynamoConfig.UserTableOverride) > 0 {
		usr = config.DynamoConfig.UserTableOverride
		log.Printf("Found User table override %v\n", usr)
	} else {
		usr = "Users"
	}

	dbManager := DynamoDBIncidentManager{
		&config.DynamoConfig.Region,
		&config.DynamoConfig.Endpoint,
		&incs, &attach,
	}
	dbManager.Initialize()
	incidentManager = &dbManager

	udbManager := DynamoDBUserManager{
		&config.DynamoConfig.Region,
		&config.DynamoConfig.Endpoint,
		&usr,
		config.User.DefaultPermissions,
	}
	udbManager.Initialize()
	userManager = &udbManager

	ensureAdminAccount(0)
}

func setupMySQLIncidentManager(config Config, db *sql.DB) {
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

func setupRuntimeUsermanager(config Config) {
	userManager = RuntimeUserManager{make(map[int64]*User), make(map[int64]string), make(map[int64][]string), config.User.DefaultPermissions}
	_, res := userManager.AddUser(&admin)
	userManager.SetPermissions(res.Id, adminPermissions)
}

func setupSQLUsermanager(config Config, db *sql.DB) {
	usrMySQLManager := MySQLUserManager{db, config.User.DefaultPermissions}
	usrMySQLManager.Initialize()

	userManager = usrMySQLManager
	ensureAdminAccount(1)
}
