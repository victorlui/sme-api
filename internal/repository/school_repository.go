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

type SchoolRepository struct {
	db *pgxpool.Pool
}

func NewSchoolRepository(db *pgxpool.Pool) *SchoolRepository {
	return &SchoolRepository{db: db}
}

func (s *SchoolRepository) CreateNewSchool(school model.School) (model.School, error) {
	query := `INSERT INTO schools (name) VALUES ($1) RETURNING id`

	row := s.db.QueryRow(context.Background(), query, school.Name)

	var id int

	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// => 42601
			if pgErr.Code == "23505" {
				return model.School{}, errors.New("ESCOLA JÁ REGISTRADO COM ESSE NOME")
			}
		}
		return model.School{}, err
	}

	school.ID = id

	return school, nil

}

func (conn *SchoolRepository) UpdateSchool(id int, school model.School) (model.School, error) {
	var exists bool
	queryCheck := "SELECT EXISTS(SELECT 1 FROM schools WHERE id = $1)"

	err := conn.db.QueryRow(context.Background(), queryCheck, id).Scan(&exists)

	if err != nil {
		return model.School{}, fmt.Errorf("erro ao verificar escola %v", err.Error())
	}

	if !exists {
		return model.School{}, fmt.Errorf("nenhuma escola encontrada para ser atualizada")
	}

	var idSchool int
	query := "UPDATE schools SET name = $1, updated_at = NOW() WHERE id = $2 RETURNING id"

	err = conn.db.QueryRow(context.Background(), query, school.Name, id).Scan(&idSchool)

	if err != nil {
		return model.School{}, fmt.Errorf("erro ao atualizar a escola %v", err)
	}

	school.ID = idSchool

	return school, nil

}

func (s *SchoolRepository) GetSchoolById(id int) (model.School, error) {
	query := "SELECT id, name FROM schools WHERE id = $1 "

	var school model.School

	err := s.db.QueryRow(context.Background(), query, id).Scan(&school.ID, &school.Name)

	if err != nil {
		if err == pgx.ErrNoRows {
			return model.School{}, errors.New("NENHUMA ESCOLA ENCONTRADA")
		}
		return model.School{}, err
	}

	return school, nil
}

func (conn *SchoolRepository) GetAllSchools(offset, limit int, name string) ([]model.School, int, error) {
	query := `SELECT id, name FROM schools`
	queryCount := "SELECT COUNT(*) FROM schools WHERE name ILIKE $1"

	var totalRecords int
	var params []interface{}
	paramIndex := 1

	decodedName, err := url.QueryUnescape(name)
	if err != nil {
		return nil, 0, err
	}
	decodedName = strings.Trim(decodedName, `"'`)

	if decodedName != "" {
		query += fmt.Sprintf(" WHERE name ILIKE $%d", paramIndex)
		params = append(params, "%"+decodedName+"%")
		paramIndex++
	}

	err = conn.db.QueryRow(context.Background(), queryCount, "%"+name+"%").Scan(&totalRecords)
	if err != nil {
		fmt.Println("err:", err)
		return nil, 0, err
	}

	// Adicionando paginação
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	params = append(params, limit, offset)

	rows, err := conn.db.Query(context.Background(), query, params...)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var schools []model.School

	for rows.Next() {
		var school model.School

		if err := rows.Scan(
			&school.ID,
			&school.Name,
		); err != nil {
			return nil, 0, err
		}

		schools = append(schools, school)
	}

	if len(schools) == 0 {
		return []model.School{}, 0, nil
	}

	return schools, totalRecords, nil
}

func (conn *SchoolRepository) DeleteSchool(id int) error {
	var count int
	query := "SELECT COUNT(*) FROM student WHERE school_id = $1"

	err := conn.db.QueryRow(context.Background(), query, id).Scan(&count)

	if err != nil {
		return fmt.Errorf("erro ao verificar alunos: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("não foi possivel deletar escola, pois há alunos nela")
	}

	_, err = conn.db.Exec(context.Background(), "DELETE FROM schools WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("não foi possivel deletar escola: %v", err)
	}

	return nil
}
