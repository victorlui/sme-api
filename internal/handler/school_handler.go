package handler

import (
	"fmt"
	"log"

	"github.com/victorlui/sma-api/internal/model"
	"github.com/victorlui/sma-api/internal/repository"
)

type SchoolHandler struct {
	repo *repository.SchoolRepository
}

func NewSchoolHandler(repo *repository.SchoolRepository) *SchoolHandler {
	return &SchoolHandler{repo: repo}
}

func (u *SchoolHandler) CreateNewSchool(req model.School) (*model.School, error) {

	school := model.School{
		Name: req.Name,
	}

	createdUser, err := u.repo.CreateNewSchool(school)

	// retorna apenas os campos informado na resposta
	reponse := model.School{
		ID:   createdUser.ID,
		Name: createdUser.Name,
	}

	if err != nil {
		return nil, err
	}

	return &reponse, nil
}

func (r *SchoolHandler) UpdateSchoolHandler(id int, req model.School) (*model.School, error) {

	updateSchool, err := r.repo.UpdateSchool(id, req)

	if err != nil {
		return &model.School{}, err
	}

	return &updateSchool, nil
}

func (r *SchoolHandler) GetSchoolsAll(page, limit int, name string) ([]model.School, int, error) {
	offset := (page - 1) * limit
	students, totalRecords, err := r.repo.GetAllSchools(offset, limit, name)

	if err != nil {

		log.Printf("Erro ao obter escolas: %v", err)

		return nil, 0, fmt.Errorf("erro ao obter escolas: %w", err)
	}

	return students, totalRecords, nil
}

func (r *SchoolHandler) DeleteSchool(id int) error {
	err := r.repo.DeleteSchool(id)

	if err != nil {
		return err
	}

	return nil
}
