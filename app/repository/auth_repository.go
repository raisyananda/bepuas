package repository

import (
	"database/sql"

	"bepuas/app/model"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) FindByUsernameOrEmail(username string) (model.User, error) {
	query := `
		SELECT id, username, email, password_hash, role_id, is_active
		FROM users
		WHERE username = $1 OR email = $1
	`

	var user model.User
	err := r.DB.QueryRow(query, username).
		Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.RoleID,
			&user.IsActive,
		)

	return user, err
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
		SELECT p.code
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		JOIN roles r ON r.id = rp.role_id
		JOIN users u ON u.role_id = r.id
		WHERE u.id = $1
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}

	return perms, nil
}
