package auth

import (
	"context"
	"database/sql"
	"github.com/mamontmodest/go-rest-api/internal/entity"
	"github.com/mamontmodest/go-rest-api/pkg/db"
	errs "github.com/mamontmodest/go-rest-api/pkg/errors"
	"hash/fnv"
	"time"
)

type Repository interface {
	Login(ctx context.Context, identity entity.Identity) error
	CreateUser(ctx context.Context, identity entity.Identity) error
}

type repository struct {
	db *db.SDatabase
}

// NewRepository creates a new auth repository
func NewRepository(db *db.SDatabase) Repository {
	return repository{db}
}

func (r repository) Login(ctx context.Context, identity entity.Identity) error {
	conn, err := r.db.ConnWith(ctx)
	person := new(entity.Identity)
	person.Login = identity.Login
	person.Password = identity.Password
	if err != nil {
		return err
	}
	defer conn.Close()
	ct, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	resp := make(chan error)
	go func() {
		h := hash(identity.Password)
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		defer tx.Rollback()
		query := "select login from  users where login = $1 and  password = $2"
		err = tx.QueryRow(query, identity.Login, h).Scan(&identity.Login)
		if err != nil {
			resp <- err
		}
		resp <- nil

	}()
	for {
		select {
		case <-ct.Done():
			return errs.CtxError{}
		case err := <-resp:
			return err
		}
	}
}

func (r repository) CreateUser(ctx context.Context, identity entity.Identity) error {
	conn, err := r.db.ConnWith(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	ct, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	resp := make(chan error)
	go func() {
		h := hash(identity.Password)
		query := "insert into users (login, password) values ($1, $2)"
		_, err := conn.Query(query, identity.Login, h)
		if err != nil {
			resp <- err
			return
		}
		resp <- err
		return
	}()
	for {
		select {
		case <-ct.Done():
			return errs.CtxError{}
		case err := <-resp:
			return err
		}
	}
}

type chanIdentity struct {
	Identity entity.Identity
	err      error
}

func hash(s string) uint32 {
	h := fnv.New32()
	h.Write([]byte(s))
	return h.Sum32()
}
