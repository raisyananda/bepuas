package repository

import (
	"database/sql"
	"errors"

	"bepuas/app/model"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) FindByUsernameOrEmail(identifier string) (model.User, string, error) {
	var u model.User
	var hash string
	err := r.DB.QueryRow(`
		SELECT id, username, email, full_name, role_id, is_active, created_at, password_hash
		FROM users WHERE username=$1 OR email=$1
	`, identifier).Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, "", errors.New("not found")
		}
		return u, "", err
	}
	return u, hash, nil
}

func (r *AuthRepository) GetPermissionsByRole(roleID string) ([]string, error) {
	rows, err := r.DB.Query(`
		SELECT p.name FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var perms []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		perms = append(perms, name)
	}
	return perms, nil
}
