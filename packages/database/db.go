// packages/database/db.go
package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	Connection *sql.DB
}

type User struct {
	Name      string
	Password  string
	SessionId string
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
	//fmt.Println("started connecting")
	config := NewEnvDBConfig()
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.GetUsername(), config.GetPassword(), config.GetHost(), config.GetPort(), config.GetDatabase())
	fmt.Println(connectionString)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("error connecting to db", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Database connection is not established:", err)
		return nil, err
	}

	_, err = db.Exec("SELECT 1;")
	if err != nil {
		fmt.Println("Failed to execute test query:", err)
	} else {
		fmt.Println("Test query executed successfully")
	}
	//fmt.Println("done connecting")
	return db, nil
}

func AddUser(db *sql.DB, username, password string) error {

	hashedPassword, err := hashPassword(password)
	if err != nil {
		fmt.Println("error hashing passwd:", err)
		return err
	}
	// Generate a new UUID
	userUUID, err := MakeUuid(db)

	// Insert the new user into the database
	query := "INSERT INTO authentification (username, password, uuid) VALUES (?, ?, ?)"
	_, err = db.Exec(query, username, hashedPassword, userUUID)
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}

	fmt.Println("User added successfully")
	return nil
}

func CheckUsername(db *sql.DB, username string) (bool, error) {
	var username1 string

	//fmt.Println("username:", username)

	username = strings.TrimSpace(username)

	err := db.QueryRow("SELECT username FROM authorized_usernames WHERE username = ?", username).Scan(&username1)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows returned for username:", username)
			return false, nil
		}
		fmt.Println("Error executing query:", err)
		return false, err
	}

	fmt.Println("username1:", username1)
	if username1 == "" {
		// User does not exist
		return false, nil
	}
	return true, nil

}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}
	return string(hashedPassword), nil
}

func CheckUserCredentials(db *sql.DB, username, password string) (bool, error) {
	var storedHash string

	var query string
	query = fmt.Sprintf("SELECT password FROM authentification WHERE username = '%s'; ", username)
	err := db.QueryRow(query).Scan(&storedHash)

	if err == sql.ErrNoRows {
		// User does not exist
		return false, nil
	} else if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		// Password does not match
		return false, nil
	}

	// User exists and password matches
	return true, nil
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

func GetUsers(db *sql.DB) (User, error) {
	var user User

	err := db.QueryRow("SELECT USERNAME, PASSWORD FROM authentification").Scan(&user.Name, &user.Password)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func MakeUuid(db *sql.DB) (string, error) {
	var compUuid string
	var newUuid uuid.UUID
	for {
		newUuid = uuid.NewV4()
		err := db.QueryRow("SELECT UUID FROM authentication WHERE UUID = ?", newUuid).Scan(&compUuid)

		if err == sql.ErrNoRows {
			// UUID does not exist
			// continuing the loop until there is a new uuid
		} else if err != nil {
			return "", err
		} else {
			//break if no one has the same uuid
			break
		}

	}

	return newUuid.String(), nil
}

func CheckUuid(db *sql.DB, uuid string) (bool, error) {
	var count int
	// Query to count how many times the UUID exists
	err := db.QueryRow("SELECT COUNT(*) FROM authentication WHERE UUID = ?", uuid).Scan(&count)
	if err != nil {
		// Handle any potential error
		return false, err
	}

	// If count is greater than 0, the UUID exists so true is retuwurned
	return count > 0, nil
}

/*func CreateNewUser(db *sql.DB, username string, password string, uuid string, id int) error {

}*/

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
