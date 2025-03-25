package users

import (
	"fmt"
	"http_server/config"
	"http_server/repository"
	"http_server/repository/models"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
    PqUniqueViolation = "23505"
)

type UsersRepo struct {
    db *sqlx.DB
    cfg config.PostgreSQLConfig
}

func NewUsersRepo(cfg config.PostgreSQLConfig) (*UsersRepo, error) {
    connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
        cfg.Host, cfg.Port, cfg.DB, cfg.User, cfg.Password,
    )
    db, err := sqlx.Connect("postgres", connStr)

    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return &UsersRepo{
        db: db,
    }, nil
}

func (r *UsersRepo) GetUser(key uuid.UUID) (*models.User, error) {
    row := r.db.QueryRowx("SELECT * FROM users WHERE id = $1", key)
    var user models.User

    if err := row.StructScan(&user); err != nil {
        log.Printf("getting user: %s", err.Error())
        return nil, repository.ErrorUserNotFound
    }

    return &user, nil
}

func (r *UsersRepo) GetUserByCred(login string, password string) (*models.User, error) {
    row := r.db.QueryRowx("SELECT * FROM users WHERE login = $1 AND password = $2", login, password)
    var user models.User

    if err := row.StructScan(&user); err != nil {
        log.Printf("getting user by cred: %s", err.Error())
        return nil, repository.ErrorWrongUserCreds
    }

    return &user, nil
}

func (r *UsersRepo) PostUser(key uuid.UUID, login string, password string) error {
    res, err := r.db.Exec(`INSERT INTO users (id, login, password)
        VALUES ($1, $2, $3)`, key, login, password)

    if err != nil {
        log.Printf("posting user: %s", err.Error())

        if err.(*pq.Error).Code == PqUniqueViolation {
            return repository.ErrorUserAlreadyExists
        } else {
            return repository.ErrorInternalQueryError
        }
    }

    if n, err := res.RowsAffected(); err != nil || n != 1 {
        return repository.ErrorUserKeyAlreadyUsed
    }

    return nil
}
