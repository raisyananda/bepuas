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
		    submitted_at = NOW()
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

func (r *AchievementRefRepository) FindByMongoID(mongoID string) (*model.AchievementReference, error) {
	var ref model.AchievementReference

	err := r.DB.QueryRow(`
		SELECT id, student_id, mongo_achievement_id, status, created_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1
		  AND status != 'deleted'
	`, mongoID).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &ref, nil
}

func (r *AchievementRefRepository) Reject(mongoID, RejectionNote string) error {
	res, err := r.DB.Exec(`
		UPDATE achievement_references
		SET status = 'rejected',
		    rejection_note = $2,
		    updated_at = NOW()
		WHERE mongo_achievement_id = $1
		  AND status = 'rejected'
	`, mongoID, RejectionNote)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("prestasi belum disubmit atau sudah diverifikasi")
	}

	return nil
}

func (r *AchievementRefRepository) Verify(mongoID, userID string) error {
	res, err := r.DB.Exec(`
		UPDATE achievement_references
		SET status = 'verified',
		    verified_at = NOW(),
		    verified_by = $2,
		    updated_at = NOW()
		WHERE mongo_achievement_id = $1
		  AND status = 'submitted'
	`, mongoID, userID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("prestasi belum disubmit atau sudah diverifikasi")
	}

	return nil
}

func (r *AchievementRefRepository) FindHistoryByMongoID(mongoID string) ([]model.AchievementReference, error) {
	rows, err := r.DB.Query(`
		SELECT status, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1
		ORDER BY updated_at ASC
	`, mongoID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []model.AchievementReference
	for rows.Next() {
		var h model.AchievementReference
		if err := rows.Scan(
			&h.Status,
			&h.RejectionNote,
			&h.CreatedAt,
			&h.UpdatedAt,
		); err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, nil
}
