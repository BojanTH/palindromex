/*
Package db is a wrapper for database interactions
*/
package db

import (
	"github.com/jinzhu/gorm"
	// load driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Connection wrapper
type Connection struct {
	Conn             *gorm.DB
	IsOpen           bool
	connectionString string
}

// NewConnection instance
func NewConnection(connectionString string) *Connection {
	return &Connection{connectionString: connectionString}
}

// Open connection
func (connection *Connection) Open() {
	conn, err := gorm.Open("postgres", connection.connectionString)
	if err != nil {
		panic("Failed to connect to database")
	}
	// Configure defaults
	conn.SingularTable(true)

	connection.Conn = conn
	connection.IsOpen = true
}

// Close connection
func (connection *Connection) Close() {
	connection.Conn.Close()
	connection.IsOpen = false
}

func (connection *Connection) Exec(sql string, values ...interface{}) *gorm.DB {
	connection.Open()
	defer connection.Close()
	result := connection.Conn.Exec(sql, values)
	return result;
}