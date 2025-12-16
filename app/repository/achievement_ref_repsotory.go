package repository

import (
	"bepuas/app/model"
	"database/sql"
	"errors"
	"log"
)

type AchievementRefRepository struct {
	DB *sql.DB
}

func NewAchievementRefRepository(db *sql.DB) *AchievementRefRepository {
	return &AchievementRefRepository{DB: db}
}

func (r *AchievementRefRepository) Create(id string, studentID string, mongoID string, status string) error {

	query := `
		INSERT INTO achievement_references (id, student_id, mongo_achievement_id, status)
		VALUES ($1, $2, $3, $4)
	`

	res, err := r.DB.Exec(query, id, studentID, mongoID, status)
	if err != nil {
		log.Println("PG INSERT ERROR:", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("PG ROWS AFFECTED ERROR:", err)
		return err
	}

	if r.DB == nil {
		return errors.New("database postgres belum diinisialisasi")
	}

	log.Println("PG INSERT OK, rows affected:", rows)
	return nil
}

func (r *AchievementRefRepository) FindByStudentID(studentID string) ([]model.AchievementReference, error) {
	rows, err := r.DB.Query(`
		SELECT id, student_id, mongo_achievement_id, status, created_at
		FROM achievement_references
		WHERE student_id = $1
		  AND status != 'deleted'
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.CreatedAt,
		); err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func (r *AchievementRefRepository) UpdateStatus(mongoID, status string) error {
	_, err := r.DB.Exec(`
		UPDATE achievement_references
		SET status = $1,
		    updated_at = NOW()
		WHERE mongo_achievement_id = $2
		  AND status = 'draft'
	`, status, mongoID)

	return err
}

func (r *AchievementRefRepository) DeleteDraft(mongoID string) error {
	res, err := r.DB.Exec(`
		UPDATE achievement_references
		SET status = 'deleted',
		    updated_at = NOW()
		WHERE mongo_achievement_id = $1
		  AND status = 'draft'
	`, mongoID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("prestasi bukan draft atau sudah diproses")
	}

	return nil
}

func (r *AchievementRefRepository) FindAll() ([]model.AchievementReference, error) {
	rows, err := r.DB.Query(`
		SELECT id, student_id, mongo_achievement_id, status, created_at
		FROM achievement_references
		WHERE status != 'deleted'
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.CreatedAt,
		); err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}
