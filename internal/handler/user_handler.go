package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/victorlui/sma-api/internal/model"
	"github.com/victorlui/sma-api/internal/repository"
	"github.com/victorlui/sma-api/internal/utils"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (u *UserHandler) GetAllUsers() ([]model.CreateUserResponse, error) {
	ctx := context.Background()

	return u.repo.GetAllUsers(ctx)
}

func (u *UserHandler) CreateNewUser(req model.CreateUserRequest) (*model.CreateUserResponse, error) {

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	createdUser, err := u.repo.CreateNewUser(user)

	// retorna apenas os campos informado na resposta
	reponse := model.CreateUserResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
		Phone: createdUser.Phone,
	}

	if err != nil {
		return nil, err
	}

	return &reponse, nil
}

func (u *UserHandler) GetUserByID(id int) (*model.User, error) {
	user, err := u.repo.GetUserByID(id)

	if err != nil {
		return &model.User{}, err
	}

	return &user, nil
}

func (u *UserHandler) Login(ctx context.Context, email, password string) (model.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)

	if err != nil {
		fmt.Println("Usuario não encontrado")
		return model.User{}, errors.New("usuário ou senha inválidos")
	}

	//log.Printf("Usuário encontrado: %+v\n", user)

	if !utils.CheckPassword(user.Password, password) {
		fmt.Println("CheckPassword", user.Password)
		return model.User{}, errors.New("usuário ou senha inválidos")
	}

	// Gerar JWT
	token, err := utils.GenerateJWT(user.ID, user.Email, time.Hour*24) // Expira em 24h
	if err != nil {
		return model.User{}, errors.New("error ao gerar jwt")
	}

	reponse := model.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Token: token,
	}
	return reponse, nil
}
