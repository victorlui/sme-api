package handler

import (
	"github.com/victorlui/sma-api/internal/model"
	"github.com/victorlui/sma-api/internal/repository"
)

type ServiceStudentHandler struct {
	repo *repository.ServiceStudentRepository
}

func NewServiceStudentHandler(repo *repository.ServiceStudentRepository) *ServiceStudentHandler {
	return &ServiceStudentHandler{repo: repo}
}

func (r *ServiceStudentHandler) CreateNewServiceStudentHandler(req model.ServiceStudent) (*model.ServiceStudent, error) {

	student_service := model.ServiceStudent{
		StudentID:   req.StudentID,
		UserID:      req.UserID,
		DateService: req.DateService,
		Reason:      req.Reason,
		File:        req.File,
	}

	createdStudentService, err := r.repo.CreateServiceStudent(student_service)

	if err != nil {
		return &model.ServiceStudent{}, err
	}

	return &createdStudentService, nil
}

func (r *ServiceStudentHandler) UpdateNewServiceStudentHandler(id int, req model.ServiceStudent) (*model.ServiceStudent, error) {

	student_service := model.ServiceStudent{
		StudentID:   req.StudentID,
		UserID:      req.UserID,
		DateService: req.DateService,
		Reason:      req.Reason,
		File:        req.File,
	}

	createdStudentService, err := r.repo.UpdateServiceStudent(id, student_service)

	if err != nil {
		return &model.ServiceStudent{}, err
	}

	return &createdStudentService, nil
}

func (r *ServiceStudentHandler) GetAllServiceStudents(page, limit int, date, idStudent string) (*[]model.ServiceStudentResponse, error) {
	offset := (page - 1) * limit

	services_students, err := r.repo.GetAllServicesStudents(offset, limit, date, idStudent)

	if err != nil {
		return &[]model.ServiceStudentResponse{}, err
	}

	return &services_students, nil
}

func (r *ServiceStudentHandler) DeleteServiceStudentHandler(id int) error {
	err := r.repo.DeleteServicesStudents(id)

	if err != nil {
		return err
	}

	return nil
}
