package repository

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/victorlui/sma-api/internal/model"
)

type StudentRepository struct {
	db         *pgxpool.Pool
	schoolRepo *SchoolRepository
}

func NewStudentRepository(db *pgxpool.Pool, schoolRepo *SchoolRepository) *StudentRepository {
	return &StudentRepository{db: db, schoolRepo: schoolRepo}
}

func (conn *StudentRepository) CreateNewStudent(student model.Student) (model.Student, error) {

	idSchool := int(student.SchoolID)
	_, err := conn.schoolRepo.GetSchoolById(idSchool)

	if err != nil {
		return model.Student{}, err
	}

	query := "INSERT INTO student (name, serie, date_of_birth, name_of_mother, name_of_dad,ra, school_id) VALUES ($1, $2, $3 ,$4, $5, $6, $7) RETURNING id"

	row := conn.db.QueryRow(context.Background(), query, student.Name, student.Serie, student.Dateofbirth, student.NameOfMother, student.NameOfDad, student.Ra, student.SchoolID)

	var id int

	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)
			if pgErr.Code == "23505" {
				return model.Student{}, errors.New("ALUNO JÁ REGISTRADO")
			}
		}
		return model.Student{}, err
	}

	student.ID = id

	return student, nil

}

func (conn *StudentRepository) GetAllStudents(offset, limit int, name string) ([]model.StudentResponse, int, error) {
	query := `SELECT s.id, s.name, s.serie, s.date_of_birth, s.name_of_mother, s.name_of_dad, s.ra, 
                     sc.id AS school_id, sc.name AS school_name 
              FROM student s 
              INNER JOIN schools sc ON s.school_id = sc.id`

	queryCount := "SELECT COUNT(*) FROM student WHERE name ILIKE $1"

	var totalRecords int
	var params []interface{}
	paramIndex := 1

	// Decodifica o nome (caso venha com aspas ou caracteres especiais)
	decodedName, err := url.QueryUnescape(name)
	if err != nil {
		return nil, 0, err
	}
	decodedName = strings.Trim(decodedName, `"'`)

	if decodedName != "" {
		query += fmt.Sprintf(" WHERE s.name ILIKE $%d", paramIndex)
		params = append(params, "%"+decodedName+"%")
		paramIndex++
	}

	err = conn.db.QueryRow(context.Background(), queryCount, "%"+name+"%").Scan(&totalRecords)

	if err != nil {
		return []model.StudentResponse{}, 0, err
	}

	// Adicionando paginação
	query += fmt.Sprintf(" ORDER BY s.id LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	params = append(params, limit, offset)

	rows, err := conn.db.Query(context.Background(), query, params...)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var students []model.StudentResponse

	for rows.Next() {
		var student model.StudentResponse

		if err := rows.Scan(
			&student.ID,
			&student.Name,
			&student.Serie,
			&student.Dateofbirth,
			&student.NameOfMother,
			&student.NameOfDad,
			&student.Ra,
			&student.School.ID,
			&student.School.Name,
		); err != nil {
			return nil, 0, err
		}

		students = append(students, student)
	}

	if len(students) == 0 {
		return []model.StudentResponse{}, 0, nil
	}

	return students, totalRecords, nil
}

func (conn *StudentRepository) GetStudentByID(id int) (model.StudentResponse, error) {
	query := "SELECT s.id, s.name, s.serie, s.date_of_birth, s.name_of_mother, s.name_of_dad, s.ra,sc.id AS school_id, sc.name AS school_name FROM student s INNER JOIN schools sc ON s.school_id = sc.id where s.id = $1"

	var student model.StudentResponse

	err := conn.db.QueryRow(context.Background(), query, id).Scan(
		&student.ID,
		&student.Name,
		&student.Serie,
		&student.Dateofbirth,
		&student.NameOfMother,
		&student.NameOfDad,
		&student.Ra,
		&student.School.ID,
		&student.School.Name,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return model.StudentResponse{}, errors.New("NENHUM ALUNO ENCONTRADO")
		}
		return model.StudentResponse{}, err
	}

	return student, nil

}

func (conn *StudentRepository) DeleteStudentRepository(id int) error {

	var isExist bool

	queryExist := "SELECT EXISTS(SELECT 1 FROM student WHERE id = $1)"

	err := conn.db.QueryRow(context.Background(), queryExist, id).Scan(&isExist)

	if err != nil {
		return fmt.Errorf("erro ao verificar alunos %v", err)
	}

	if !isExist {
		return fmt.Errorf("nenhum aluno encontrado com esse id")
	}

	// Iniciar transação
	tx, err := conn.db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background()) // Faz rollback se algo der errado
		}
	}()

	// Deletar os serviços relacionados ao aluno
	queryDeleteServices := "DELETE FROM service_students WHERE student_id = $1"
	_, err = tx.Exec(context.Background(), queryDeleteServices, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar atendimento do aluno: %w", err)
	}

	queryDeleteStudent := "DELETE FROM student WHERE id = $1"
	_, err = tx.Exec(context.Background(), queryDeleteStudent, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar aluno: %w", err)
	}

	// Confirmar a transação
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil

}
