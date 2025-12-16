package repository

import (
	"bepuas/app/model"
	"database/sql"
)

type StudentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

func (r *StudentRepository) FindAll() ([]model.Student, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.Student
	for rows.Next() {
		var s model.Student
		var advisor sql.NullString

		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&advisor,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}

		if advisor.Valid {
			s.AdvisorID = advisor.String
		}

		data = append(data, s)
	}

	return data, nil
}

func (r *StudentRepository) FindByID(id string) (*model.Student, error) {
	var s model.Student
	var advisor sql.NullString

	err := r.DB.QueryRow(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE id = $1
	`, id).Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&advisor,
		&s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if advisor.Valid {
		s.AdvisorID = advisor.String
	}

	return &s, nil
}

func (r *StudentRepository) FindByAdvisor(lecturerID string) ([]model.Student, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id
		FROM students
		WHERE advisor_id = $1
	`, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.AdvisorID,
		); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

func (r *StudentRepository) UpdateAdvisor(studentID, advisorID string) error {
	_, err := r.DB.Exec(`
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
	`, advisorID, studentID)

	return err
}

func (r *StudentRepository) IsAdvisee(studentID, lecturerID string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM students
			WHERE id = $1 AND advisor_id = $2
		)
	`, studentID, lecturerID).Scan(&exists)

	return exists, err
}
