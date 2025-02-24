package handler

import (
	"fmt"
	"log"

	"github.com/victorlui/sma-api/internal/model"
	"github.com/victorlui/sma-api/internal/repository"
)

type StudentHandler struct {
	repo *repository.StudentRepository
}

func NewStudentHandler(repo *repository.StudentRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

func (r *StudentHandler) CreateNewStudent(req model.Student) (*model.Student, error) {

	studentModel := model.Student{
		Name:         req.Name,
		Serie:        req.Serie,
		Dateofbirth:  req.Dateofbirth,
		NameOfMother: req.NameOfMother,
		NameOfDad:    req.NameOfDad,
		Ra:           req.Ra,
		SchoolID:     req.SchoolID,
	}

	student, err := r.repo.CreateNewStudent(studentModel)

	if err != nil {
		fmt.Println("Erro handler", err)
		return &model.Student{}, err
	}

	return &student, nil
}

func (r *StudentHandler) GetAllStudents(page, limit int, name string) ([]model.StudentResponse, int, error) {
	offset := (page - 1) * limit
	students, totalRecords, err := r.repo.GetAllStudents(offset, limit, name)

	if err != nil {

		log.Printf("Erro ao obter estudantes: %v", err)

		return nil, 0, fmt.Errorf("erro ao obter estudantes: %w", err)
	}

	return students, totalRecords, nil
}

func (r *StudentHandler) GetStudentByIDHandler(id int) (model.StudentResponse, error) {
	student, err := r.repo.GetStudentByID(id)

	if err != nil {
		return model.StudentResponse{}, err
	}

	return student, nil
}

func (r *StudentHandler) DeleteHandler(id int) error {

	err := r.repo.DeleteStudentRepository(id)

	if err != nil {
		return err
	}

	return nil
}
