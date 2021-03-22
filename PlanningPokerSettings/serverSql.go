package PlanningPokerSettings

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

type ServerSql struct {
	DSN     string
	TypeSql string
}

func (serverSql ServerSql) OpenConnection() (*sql.DB, error) {
	return sql.Open(serverSql.TypeSql, serverSql.DSN)
}
