package settings

type MongoDbOptions struct {
	DatabaseName     string `json:"databaseName"`
	ConnectionString string `json:"connectionString"`
}
