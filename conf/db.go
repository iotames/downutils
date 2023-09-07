package conf

import (
	"fmt"
)

const DRIVER_SQLITE3 = "sqlite3"
const DRIVER_MYSQL = "mysql"
const DRIVER_POSTGRES = "postgres"
const SQLITE_FILENAME = "sqlite3.db"

type Database struct {
	Driver, Host, Username, Password, Name string
	Port, NodeID                           int
}

func GetDatabase() Database {
	return Database{Driver: dbDriver, Host: dbHost, Username: dbUsername, Password: dbPassword, Name: dbName, Port: dbPort, NodeID: dbNodeId}
}

func (d Database) GetAddr() string {
	return fmt.Sprintf("%s:%d", d.Host, d.Port)
}

func (d Database) GetDSN() string {
	dsnMap := map[string]string{
		DRIVER_MYSQL:    fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", d.Username, d.Password, d.Host, d.Port, d.Name),
		DRIVER_POSTGRES: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", d.Host, d.Username, d.Password, d.Name, d.Port),
	}
	dsn, ok := dsnMap[d.Driver]
	if !ok {
		dsnLen := len(dsnMap)
		ds := make([]string, dsnLen)
		for k := range dsnMap {
			ds = append(ds, k)
		}
		errMsg := fmt.Sprintf("ENV error: DB_DRIVER only Support: %v", ds)
		panic(errMsg)
	}
	return dsn
}
