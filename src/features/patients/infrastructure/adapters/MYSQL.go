package adapters

import (
    "apiconsumer/src/core"
	"database/sql"
)

type MYSQL struct {
	conn *sql.DB
}

func NewMysql() *MYSQL {
	conn := core.NewMysql()
	return &MYSQL{
		conn: conn,
	}

}

func (m *MYSQL) FindID(id int) (error) {
	result := m.conn.QueryRow("SELECT id FROM patients WHERE id = ?", id)

	if err := result.Scan(&id); err != nil {
		return err
	}
	return nil
}

