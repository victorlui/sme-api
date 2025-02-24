package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/victorlui/sma-api/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]model.CreateUserResponse, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, email, phone FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.CreateUserResponse
	for rows.Next() {
		var user model.CreateUserResponse
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) CreateNewUser(user model.User) (model.User, error) {
	query := `INSERT INTO users (name,email,phone,password) VALUES ($1,$2,$3,$4) RETURNING id`

	row := r.db.QueryRow(context.Background(), query, user.Name, user.Email, user.Phone, user.Password)

	var id int

	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
			if pgErr.Code == "23505" {
				return model.User{}, errors.New("USUÁRIO JÁ REGISTRADO COM ESSE E-MAIL")
			}
		}
		return model.User{}, err
	}

	user.ID = id

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	query := "SELECT id, name, email, phone, password FROM users WHERE email = $1"

	var user model.User

	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password)

	if err != nil {
		return model.User{}, nil
	}

	return user, nil
}

func (conn *UserRepository) GetUserByID(id int) (model.User, error) {
	query := "SELECT id, name, email, phone FROM users where id = $1"

	var user model.User

	err := conn.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return model.User{}, errors.New("NENHUM USUÁRIO ENCONTRADO")
		}
		return model.User{}, err
	}

	return user, nil

}
