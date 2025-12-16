package repository

import (
	"database/sql"

	"bepuas/app/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindAllUser() ([]model.UserResponse, error) {
	rows, err := r.DB.Query(`
		SELECT u.id, u.username, u.email, u.full_name, u.role_id, r.name
		FROM users u
		LEFT JOIN roles r ON r.id = u.role_id
		WHERE u.is_active = true
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserResponse

	for rows.Next() {
		var u model.UserResponse
		var roleID sql.NullString
		var roleName sql.NullString

		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.FullName,
			&roleID,
			&roleName,
		)
		if err != nil {
			return nil, err
		}

		if roleID.Valid {
			u.RoleID = &roleID.String
		}

		if roleName.Valid {
			u.RoleName = &roleName.String
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) FindUserByID(id string) (*model.User, error) {
	var u model.User
	err := r.DB.QueryRow(`
		SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at
		FROM users WHERE id=$1
	`, id).Scan(
		&u.ID, &u.Username, &u.Email,
		&u.FullName, &u.RoleID, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) CreateUser(u *model.User, passwordHash string) error {
	_, err := r.DB.Exec(`
		INSERT INTO users (id, username, email, password_hash, full_name, is_active)
		VALUES ($1,$2,$3,$4,$5, true)
	`,
		u.ID, u.Username, u.Email, passwordHash, u.FullName,
	)
	return err
}

func (r *UserRepository) GetUserRole(userID string) (roleID string, roleName string, err error) {
	var roleIDNull sql.NullString
	var roleNameNull sql.NullString

	err = r.DB.QueryRow(`
		SELECT r.id, r.name
		FROM users u
		LEFT JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`, userID).Scan(&roleIDNull, &roleNameNull)

	if err != nil {
		return "", "", err
	}

	if roleIDNull.Valid {
		roleID = roleIDNull.String
	}
	if roleNameNull.Valid {
		roleName = roleNameNull.String
	}

	return roleID, roleName, nil
}

func (r *UserRepository) DeleteUser(id string) error {
	_, err := r.DB.Exec(`
		UPDATE users SET is_active=false WHERE id=$1
	`, id)
	return err
}

func (r *UserRepository) UpdateUserRole(userID, roleID string) error {
	_, err := r.DB.Exec(`
		UPDATE users SET role_id=$1, updated_at=NOW()
		WHERE id=$2
	`, roleID, userID)
	return err
}

func (r *UserRepository) UpsertStudent(s *model.Student) error {
	_, err := r.DB.Exec(`
		INSERT INTO students (id, user_id, student_id, program_study, academic_year)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (user_id) DO UPDATE SET
			student_id=EXCLUDED.student_id,
			program_study=EXCLUDED.program_study,
			academic_year=EXCLUDED.academic_year
	`,
		s.ID, s.UserID, s.StudentID, s.ProgramStudy, s.AcademicYear,
	)
	return err
}

func (r *UserRepository) UpsertLecturer(l *model.Lecturer) error {
	_, err := r.DB.Exec(`
		INSERT INTO lecturers (id, user_id, lecturer_id, department)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (user_id) DO UPDATE SET
			lecturer_id=EXCLUDED.lecturer_id,
			department=EXCLUDED.department
	`,
		l.ID, l.UserID, l.LecturerID, l.Department,
	)
	return err
}
