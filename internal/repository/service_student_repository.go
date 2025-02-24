package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/victorlui/sma-api/internal/model"
)

type ServiceStudentRepository struct {
	db          *pgxpool.Pool
	studentRepo *StudentRepository
	userRepo    *UserRepository
}

func NewServiceStudentRepository(db *pgxpool.Pool, studentRepo *StudentRepository, userRepo *UserRepository) *ServiceStudentRepository {
	return &ServiceStudentRepository{db: db, studentRepo: studentRepo, userRepo: userRepo}
}

func (conn *ServiceStudentRepository) CreateServiceStudent(student_service model.ServiceStudent) (model.ServiceStudent, error) {
	_, err := conn.studentRepo.GetStudentByID(student_service.StudentID)

	if err != nil {
		fmt.Println("Erro ao encontrar aluno", err)
		return model.ServiceStudent{}, errors.New("NENHUM ALUNO ENCONTRADO")
	}

	_, err = conn.userRepo.GetUserByID(student_service.UserID)

	if err != nil {
		fmt.Println("Erro ao encontrar usuario", err)
		return model.ServiceStudent{}, errors.New("NENHUM USUARIO ENCONTRADO")
	}

	query := "INSERT INTO service_students (student_id, user_id, date_service, reason, file) VALUES ($1, $2, $3, $4, $5) RETURNING id;"

	row := conn.db.QueryRow(context.Background(), query,
		student_service.StudentID,
		student_service.UserID,
		student_service.DateService,
		student_service.Reason,
		student_service.File,
	)

	var id int

	if err := row.Scan(&id); err != nil {
		fmt.Println("Error service", err)
		return model.ServiceStudent{}, err
	}

	student_service.ID = id

	return student_service, nil
}

func (conn *ServiceStudentRepository) UpdateServiceStudent(id int, student_service model.ServiceStudent) (model.ServiceStudent, error) {
	var exists bool
	queryCheck := "SELECT EXISTS(SELECT 1 FROM service_students WHERE id = $1)"
	err := conn.db.QueryRow(context.Background(), queryCheck, id).Scan(&exists)
	if err != nil {
		fmt.Println("Erro ao verificar atendimento", err)
		return model.ServiceStudent{}, errors.New("erro ao verificar atendimento")
	}
	if !exists {
		return model.ServiceStudent{}, errors.New("nenhum atendimento encontrado para atualizar")
	}

	_, err = conn.studentRepo.GetStudentByID(student_service.StudentID)

	if err != nil {
		fmt.Println("Erro ao encontrar aluno", err)
		return model.ServiceStudent{}, errors.New("NENHUM ALUNO ENCONTRADO")
	}

	_, err = conn.userRepo.GetUserByID(student_service.UserID)

	if err != nil {
		fmt.Println("Erro ao encontrar usuario", err)
		return model.ServiceStudent{}, errors.New("NENHUM USUARIO ENCONTRADO")
	}

	query := "UPDATE service_students SET updated_at = NOW()"
	params := []interface{}{}
	setClauses := []string{}

	if student_service.Reason != "" {
		setClauses = append(setClauses, fmt.Sprintf("reason = $%d", len(params)+1))
		params = append(params, student_service.Reason)
	}
	if student_service.File != "" {
		setClauses = append(setClauses, fmt.Sprintf("file = $%d", len(params)+1))
		params = append(params, student_service.File)
	}

	// Verificar se ao menos um campo foi informado
	if len(setClauses) == 0 {

		return model.ServiceStudent{}, fmt.Errorf("nenhum campo para atualizar")
	}

	// Concatenar a cláusula SET
	query += ", " + strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d RETURNING id", len(params)+1)

	params = append(params, id)

	var updatedID int

	if err := conn.db.QueryRow(context.Background(), query, params...).Scan(&updatedID); err != nil {
		return model.ServiceStudent{}, fmt.Errorf("erro ao atualizar atendimento: %w", err)
	}

	student_service.ID = updatedID
	return student_service, nil
}

func (conn *ServiceStudentRepository) GetAllServicesStudents(offset, limit int, date, idStudent string) ([]model.ServiceStudentResponse, int, error) {
	var totalRecords int
	var query string
	var args []interface{}
	var conditions []string

	// Construir a consulta SQL dinamicamente com LIMIT e OFFSET
	query = `
		SELECT 
			ss.id,
			ss.reason,
			ss.file,
			ss.date_service,
			s.id,
			s.name,
			s.serie,
			s.date_of_birth,
			s.name_of_mother,
			s.name_of_dad,
			s.ra,
			u.id,
			u.name,
			u.email
		FROM 
			service_students ss
		JOIN 
			student s ON ss.student_id = s.id
		JOIN 
			users u ON ss.user_id = u.id
	`

	queryCount := "SELECT COUNT(*) FROM service_students ss JOIN student s ON ss.student_id = s.id JOIN users u ON ss.user_id = u.id"

	if date != "" {
		startDate := strings.Trim(date, `"`) + " 23:59:59"
		//endDate := strings.Trim(date_end, `"`) + " 23:59:59"

		conditions = append(conditions, "DATE(ss.date_service) = $"+strconv.Itoa(len(args)+1))
		args = append(args, startDate)
	}

	if idStudent != "" {
		conditions = append(conditions, "s.id = $"+strconv.Itoa(len(args)+1))
		args = append(args, idStudent)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := conn.db.QueryRow(context.Background(), queryCount, args...).Scan(&totalRecords)

	if err != nil {
		return []model.ServiceStudentResponse{}, 0, err
	}

	query += fmt.Sprintf(`
		ORDER BY ss.date_service DESC
		LIMIT $%d OFFSET $%d;
	`, len(args)+1, len(args)+2)

	args = append(args, limit, offset)

	rows, err := conn.db.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return []model.ServiceStudentResponse{}, 0, err
	}
	defer rows.Close()

	var service_students []model.ServiceStudentResponse

	for rows.Next() {
		var ss model.ServiceStudentResponse
		if err := rows.Scan(
			&ss.ID,
			&ss.Reason,
			&ss.File,
			&ss.DateService,
			&ss.Student.ID,
			&ss.Student.Name,
			&ss.Student.Serie,
			&ss.Student.Dateofbirth,
			&ss.Student.NameOfMother,
			&ss.Student.NameOfDad,
			&ss.Student.Ra,
			&ss.User.ID,
			&ss.User.Name,
			&ss.User.Email,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			return []model.ServiceStudentResponse{}, 0, err
		}

		service_students = append(service_students, ss)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return []model.ServiceStudentResponse{}, 0, err
	}

	if len(service_students) == 0 {
		return []model.ServiceStudentResponse{}, 0, nil
	}

	return service_students, totalRecords, nil
}

func (conn *ServiceStudentRepository) DeleteServicesStudents(id int) error {
	query := "SELECT EXISTS(SELECT 1 FROM service_students WHERE id = $1)"

	var isExist bool

	err := conn.db.QueryRow(context.Background(), query, id).Scan(&isExist)

	if err != nil {
		fmt.Println("Erro ao encontrar atendimento", err)
		return errors.New("erro ao encontrar atendimento")
	}

	if !isExist {
		return errors.New("nenhum atendimento a ser deletado")
	}

	_, err = conn.db.Exec(context.Background(), "DELETE FROM service_students WHERE id = $1", id)

	if err != nil {
		fmt.Println("Erro ao deletar atendimento", err)
		return errors.New("erro ao deletar atendimento")
	}

	return nil
}
