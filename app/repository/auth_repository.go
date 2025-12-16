package repository

import (
	"database/sql"
	"fmt"

	"bepuas/app/model"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) FindByUsernameOrEmail(identifier string) (model.User, error) {
	var u model.User

	query := `
	SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
	FROM users
	WHERE username = $1 OR email = $1
	LIMIT 1
	`

	err := r.DB.QueryRow(query, identifier).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.RoleID,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
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

func (r *AuthRepository) GetUserPermissions(userID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM users u
		JOIN role_permissions rp ON rp.role_id = u.role_id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE u.id = $1
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		perms = append(perms, code)
	}

	if len(perms) == 0 {
		return nil, fmt.Errorf("permission kosong untuk user %s", userID)
	}

	return perms, nil
}
