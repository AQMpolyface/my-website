// packages/database/db.go
package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Connection *sql.DB
}

type User struct {
	Name     string
	Password string
}
type EnvDBConfig struct {
	host     string
	port     string
	username string
	password string
	database string
}

var dataSourceName string = "polyface:@tcp(localhost:3336)/auth"
var config *EnvDBConfig

func ConnectToDB() (*sql.DB, error) {
	fmt.Println("started connecting")
	config := NewEnvDBConfig()
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.GetUsername(), config.GetPassword(), config.GetHost(), config.GetPort(), config.GetDatabase())
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("error connecting to db", err)
		return nil, err
	}
	return db, nil
}
func NewEnvDBConfig() *EnvDBConfig {
	return &EnvDBConfig{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		username: os.Getenv("DB_USERNAME"),
		password: os.Getenv("DB_PASSWORD"),
		database: os.Getenv("DB_DATABASE"),
	}
}

func (c *EnvDBConfig) GetHost() string {
	return c.host
}

func (c *EnvDBConfig) GetPort() string {
	return c.port
}

func (c *EnvDBConfig) GetUsername() string {
	return c.username
}

func (c *EnvDBConfig) GetPassword() string {
	return c.password
}

func (c *EnvDBConfig) GetDatabase() string {
	return c.database
}

func GetUsers(db *sql.DB) (User, error) {
	var user User

	err := db.QueryRow("SELECT USERNAME, PASSWORD FROM authentification").Scan(&user.Name, &user.Password)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
