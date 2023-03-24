package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Repo struct {
	Config map[string]interface{}
	conn   *sql.DB
}

func (r *Repo) Conn() (*sql.DB, error) {

	url := fmt.Sprint(r.Config["connection_string"])
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	r.conn = conn

	return conn, nil
}

func (r *Repo) Close() error {
	return r.conn.Close()
}
