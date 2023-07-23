package db

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"time"
)

type SDatabase struct {
	dns    string
	driver string
}

func NewSDatabase(dns, driver string) *SDatabase {
	return &SDatabase{
		dns:    dns,
		driver: driver,
	}
}

func (database *SDatabase) Conn() (*sql.DB, error) {
	db, err := sql.Open(database.driver, database.dns)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (database *SDatabase) ConnWith(ctx context.Context) (*sql.DB, error) {
	ct, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	resp := make(chan ResponseConn)
	go func() {
		conn, err := database.Conn()
		resp <- ResponseConn{conn, err}
	}()
	for {
		select {
		case <-ct.Done():
			return nil, errors.New("dasda")
		case val := <-resp:
			return val.conn, val.err
		}
	}
}

type ResponseConn struct {
	conn *sql.DB
	err  error
}
