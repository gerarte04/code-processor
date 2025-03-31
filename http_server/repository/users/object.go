package users

import (
	"http_server/pkg/database"
	"http_server/repository"
	"http_server/repository/models"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UsersRepo struct {
    db *sqlx.DB
    ep database.DBErrorProcessor
}

func NewUsersRepo(db *sqlx.DB, ep database.DBErrorProcessor) *UsersRepo {
    return &UsersRepo{
        db: db,
        ep: ep,
    }
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

        if r.ep.ProcessError(err) == database.ErrorUniqueViolation {
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
