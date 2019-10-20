package mysqlc

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"constants"
)

var (
	MySqlDB  *sql.DB
	mySqlErr error
)

func Init() {
	/*
	Connecting to mysql
	Raising panic, if connection is not made properly
	 */
	DBConnecting()
	MySqlDB.SetMaxOpenConns(constants.MySQLMaxOpenConnection)
	MySqlDB.SetMaxIdleConns(constants.MySQLMaxIdleConnection)
	if mysqlPingErr := MySqlDB.Ping(); mysqlPingErr != nil {
		panic(mysqlPingErr.Error())
	}
}

func DBConnecting() {
	/*
	Opening a connection to mysql
	 */
	MySqlDB, mySqlErr = sql.Open("mysql",
		constants.MySQLUserName+":"+constants.MySQLPassword+"@/"+constants.MySQLDBName)
	if mySqlErr != nil {
		panic(mySqlErr.Error())
	}
}

func DBClosing() {
	/*
	Closing mysql connection
	 */
	MySqlDB.Close()
}

