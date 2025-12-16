package repository

import (
	"bepuas/app/model"
	"database/sql"
)

type LecturerRepository struct {
	DB *sql.DB
}

func NewLecturerRepository(db *sql.DB) *LecturerRepository {
	return &LecturerRepository{DB: db}
}

func (r *LecturerRepository) FindAll() ([]model.Lecturer, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.Lecturer
	for rows.Next() {
		var l model.Lecturer
		if err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.LecturerID,
			&l.Department,
			&l.CreatedAt,
		); err != nil {
			return nil, err
		}
		data = append(data, l)
	}
	return data, nil
}

func (r *LecturerRepository) FindByID(id string) (*model.Lecturer, error) {
	var l model.Lecturer

	err := r.DB.QueryRow(`
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE id = $1
	`, id).Scan(
		&l.ID,
		&l.UserID,
		&l.LecturerID,
		&l.Department,
		&l.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &l, nil
}
