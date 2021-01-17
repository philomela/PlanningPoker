package ServerPlanningPoker

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

type ServerSql struct {
	DSN     string
	TypeSql string
}

func (serverSql ServerSql) ConfigureSqlServer(DSN string, TypeSql string) {
	serverSql.DSN = DSN
	serverSql.TypeSql = TypeSql
}

func (serverSql ServerSql) OpenConnection() (*sql.DB, error) {
	return sql.Open(serverSql.TypeSql, serverSql.DSN) //Переписать реализацию и не открывать экземпляр один для всех обработчиков
}
