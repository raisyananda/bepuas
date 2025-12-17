package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"bepuas/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementMongoRepository struct {
	Collection *mongo.Collection
}

func NewAchievementMongoRepository(db *mongo.Database) *AchievementMongoRepository {
	return &AchievementMongoRepository{
		Collection: db.Collection("achievements"),
	}

}

func (r *AchievementMongoRepository) Create(ctx context.Context, a *model.Achievement) error {

	a.ID = primitive.NewObjectID()
	a.Status = "draft"
	a.IsDeleted = false
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	if a.Attachments == nil {
		a.Attachments = []model.Attachment{}
	}

	res, err := r.Collection.InsertOne(ctx, a)
	if err != nil {
		log.Println("MONGO INSERT ERROR:", err)
		return err
	}

	log.Println("MONGO INSERT OK")
	log.Println("DB:", r.Collection.Database().Name())
	log.Println("COLLECTION:", r.Collection.Name())
	log.Println("ID:", res.InsertedID)
	return nil
}

func (r *AchievementMongoRepository) SubmitForVerification(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.Collection.UpdateOne(ctx,
		bson.M{
			"_id":        id,
			"status":     "draft",
			"is_deleted": false,
		},
		bson.M{
			"$set": bson.M{
				"status":     "submitted",
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

func (r *AchievementMongoRepository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.Collection.UpdateOne(ctx,
		bson.M{
			"_id":    id,
			"status": "draft",
		},
		bson.M{
			"$set": bson.M{
				"is_deleted": true,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

func (r *AchievementMongoRepository) FindByStudent(ctx context.Context, studentID string) ([]model.Achievement, error) {

	filter := bson.M{
		"student_id": studentID,
		"is_deleted": false,
	}

	cur, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var res []model.Achievement
	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *AchievementMongoRepository) FindAll(ctx context.Context) ([]model.Achievement, error) {

	filter := bson.M{"is_deleted": false}

	cur, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var res []model.Achievement
	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *AchievementMongoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Achievement, error) {

	var a model.Achievement
	err := r.Collection.FindOne(ctx, bson.M{
		"_id":        id,
		"is_deleted": false,
	}).Decode(&a)

	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AchievementMongoRepository) UpdateDraft(ctx context.Context, id primitive.ObjectID, studentID string, data bson.M) error {

	filter := bson.M{
		"_id":        id,
		"student_id": studentID,
		"status":     "draft",
		"is_deleted": false,
	}

	update := bson.M{
		"$set": data,
	}

	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *AchievementMongoRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.Achievement, error) {

	cursor, err := r.Collection.Find(ctx, bson.M{
		"_id": bson.M{"$in": ids},
	})
	if err != nil {
		return nil, err
	}

	var result []model.Achievement
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AchievementMongoRepository) AddAttachment(ctx context.Context, id primitive.ObjectID, studentID string,
	att model.Attachment) error {

	filter := bson.M{
		"_id":        id,
		"student_id": studentID,
		"status":     "draft",
		"is_deleted": false,
	}

	update := bson.M{
		"$push": bson.M{
			"attachments": att,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	res, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("prestasi tidak ditemukan atau tidak dalam draft")
	}

	return nil
}
