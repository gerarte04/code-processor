package users

import (
	"cpapp/http_server/repository"
	"cpapp/http_server/repository/models"
	"cpapp/pkg/database"
	"cpapp/pkg/database/postgres"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UsersRepo struct {
    db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
    return &UsersRepo{
        db: db,
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
    row := r.db.QueryRowx("SELECT * FROM users WHERE login = $1", login)
    var user models.User

    if err := row.StructScan(&user); err != nil {
        log.Printf("getting user by cred: %s", err.Error())
        return nil, repository.ErrorWrongUserCreds
    }

    if user.Password != password {
        log.Printf("getting user by cred: wrong password\n")
        return nil, repository.ErrorWrongPassword
    }

    return &user, nil
}

func (r *UsersRepo) PostUser(key uuid.UUID, login string, password string) error {
    res, err := r.db.Exec(`INSERT INTO users (id, login, password)
        VALUES ($1, $2, $3)`, key, login, password)

    if err != nil {
        log.Printf("posting user: %s", err.Error())

        if postgres.ProcessError(err) == database.ErrorUniqueViolation {
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
