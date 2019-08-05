package db

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
//	_ "github.com/lib/pq"
//	_ "github.com/mattn/go-sqlite3"
	. "pkg.deepin.io/server/utils/logger"
)

var (
	ErrRecordExist = errors.New("Record Has Exist")
)

type Model struct {
	ID        uint64    `gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

var Maria *gorm.DB

func InitMaria(username, password, host, port, database string) {
	Maria = ConnectMaria(username, password, host, port, database)
}

func ConnectMaria(username, password, host, port, database string) *gorm.DB {
	dataSourceFormat := "%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dataSource := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	Logger.Debug("Connect to database %v ", fmt.Sprintf(dataSourceFormat, username, "******", host, port, database))
	db, err := gorm.Open("mysql", dataSource)
	if nil != err {
		panic("init gorm failed: " + err.Error())
	}
	//db, err := gorm.Open("postgres", "user=gorm dbname=gorm sslmode=disable")
	// db, err := gorm.Open("sqlite3", "/tmp/gorm.db")

	// You can also use an existing database connection handle
	// dbSql, _ := sql.Open("postgres", "user=gorm dbname=gorm sslmode=disable")
	// db := gorm.Open("postgres", dbSql)

	// Get database connection handle [*sql.DB](http://golang.org/pkg/database/sql/#DB)
	db.DB()
	// Then you could invoke `*sql.DB`'s functions with it
	db.DB().Ping()
	db.DB().SetMaxIdleConns(2048)
	db.DB().SetMaxOpenConns(8096)

	// db.LogMode(true)

	// Disable table name's pluralization
	db.SingularTable(true)
	return db
}

func GetMaria() *gorm.DB {
	return Maria
}
