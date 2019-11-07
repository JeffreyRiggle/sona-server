package main

// Config defines the configuration that can be used on this web service.
// The ManagerType controls what manager to use (0 = runtime, 1 = dynamodb, 2 = mysql, 3 = datastore)
// The FileManagerType controls what file manager to use (0 = local, 1 = S3)
type Config struct {
	ManagerType     int                    `json:"managertype"`
	FileManagerType int                    `json:"filemanagertype"`
	DynamoConfig    DynamoDBConfig         `json:"dynamodb"`
	MYSQL           MySQLConfig            `json:"mysql"`
	DataStore       DataStoreConfig        `json:"datastore"`
	LocalFileConfig LocalFileManagerConfig `json:"fileconfig"`
	S3Config        S3FileManagerConfig    `json:"s3config"`
	Hooks           WebHooks               `json:"webhooks"`
	Logging         LogConfig              `json:"logging"`
	User            UserConfig             `json:"userconfig"`
}

// WebHook defines an endpoint to call.
// The method represents what method will be called on the http endpoint (example GET, PUT, POST)
// The URL is the endpoint to call (example https://mysite.com/api)
// The Body is the body to send to that endpoint.
type WebHook struct {
	Method string  `json:"method"`
	URL    string  `json:"url"`
	Body   WebBody `json:"body"`
}

// WebBody defines the json body to send to an endpoint.
type WebBody struct {
	Items []WebBodyItem `json:"items"`
}

// WebBodyItem represents a single item in the json body to send.
// The key is the json key value.
// The value is the json value.
// The subsititute flag indicates if you would like sona to attempt to swap out the value
// with known information.
//
// Example of a substitution config {"key", "incident", "value": "id", "substitute": true}
// If this was called when incident 1 was created the following json would be sent to an endpoint
// {"incident", "1"}
//
// Example of a complex subsitution {"key", "message", "value": "{{id}} was created", "substitute": true}
// If this was called when incident 1 was created the following json would be sent to an endpoint
// {"message", "1 was created"}
type WebBodyItem struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Substitute bool   `json:"substitute"`
}

// WebHooks defines the hooks to call for different events.
// The AddedHooks are web hooks to call when an incident has been created.
// The UpdatedHooks are web hooks to call when an incident has been updated.
// The AttachedHooks are web hooks to call when an attachment has been added to an incident.
// The UpdatedUserHooks are web hooks to call when a user is updated.
type WebHooks struct {
	AddedHooks       []WebHook `json:"addedhooks"`
	UpdatedHooks     []WebHook `json:"updatedhooks"`
	AttachedHooks    []WebHook `json:"attachedhooks"`
	AddedUserHooks   []WebHook `json:"addedUserHooks"`
	UpdatedUserHooks []WebHook `json:"updatedUserHooks"`
}

// DynamoDBConfig is the configuration to use if the dynamodb mananger is in use.
// The Region controls what AWS region your db will be created/maintained in.
// The IncidentTableOverride will override the default incident table name and use that instead.
// The AttachmentTableOverride will override the default attachment table name and use that instead.
type DynamoDBConfig struct {
	Region                  string `json:"region"`
	IncidentTableOverride   string `json:"incidenttableoverride"`
	AttachmentTableOverride string `json:"attachmenttableoverride"`
}

// LocalFileManagerConfig controls the configuration of the local file manager if it is in use.
// The Path represents what physical path to store incident attachments in.
type LocalFileManagerConfig struct {
	Path string `json:"path"`
}

// S3FileManagerConfig controls the configuration of the s3 file manager if it is in use.
// The Region controls what region your bucket will be stored/maintained in.
// The Bucket is the bucket to use of incident attachments.
type S3FileManagerConfig struct {
	Region string `json:"region"`
	Bucket string `json:"bucket"`
}

// LogConfig controls how logging is handled.
// Enabled controls whether logging will happen at all.
// Path controls the root folder in which log files will be stored.
type LogConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

// MySQLConfig controls the configuration of a mysql database if it is in use.
// The UserName controls what user to log in as.
// The Password controls what password to use for the user.
// The Host controls what host you will attempt to login to.
// The Port controls what port will will login to on the host.
// The DBName controls what DB you will login to.
type MySQLConfig struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"dbname"`
}

// DataStoreConfig controls the configuration of a google cloud datastore if it is in use.
// The ProjectName controls what google cloud project will be used.
// The AuthFile controls what json file to use for authentication.
type DataStoreConfig struct {
	ProjectName string `json:"projectname"`
	AuthFile    string `json:"authfile"`
}

// UserConfig to use when creating new users
// The DefaultPermissions are the permissions that should be granted to all new users.
type UserConfig struct {
	DefaultPermissions []string `json:"defaultpermissions"`
}
