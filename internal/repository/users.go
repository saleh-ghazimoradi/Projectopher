package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository/mongoDTO"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"strings"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	GetUsers(ctx context.Context, offset, limit int64) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
	CountUser(ctx context.Context) (int64, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func (u *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	userDTO, err := mongoDTO.FromUserCoreToDTO(user)
	if err != nil {
		return err
	}

	result, err := u.collection.InsertOne(ctx, userDTO)
	if err != nil {
		switch {
		case u.isDuplicateEmailError(err):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		user.Id = oid.Hex()
	}

	return nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userDTO mongoDTO.UserDTO

	err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&userDTO)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return mongoDTO.FromUserDTOToCore(&userDTO), nil
}

func (u *userRepository) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	uId, _ := u.oId(id)
	var userDTO mongoDTO.UserDTO

	if err := u.collection.FindOne(ctx, bson.M{"_id": uId}).Decode(&userDTO); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return mongoDTO.FromUserDTOToCore(&userDTO), nil
}

func (u *userRepository) GetUsers(ctx context.Context, offset, limit int64) ([]domain.User, error) {
	cursor, err := u.collection.Find(ctx, bson.M{}, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	update := bson.M{
		"$set": bson.M{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	}
	result, err := u.collection.UpdateOne(ctx, bson.M{"_id": user.Id}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (u *userRepository) DeleteUser(ctx context.Context, id string) error {
	uId, _ := u.oId(id)
	_, err := u.collection.DeleteOne(ctx, bson.M{"_id": uId})
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) CountUser(ctx context.Context) (int64, error) {
	return u.collection.CountDocuments(ctx, bson.M{})
}

func (u *userRepository) oId(id string) (bson.ObjectID, error) {
	oid, err := bson.ObjectIDFromHex(id)
	return oid, err
}

func (u *userRepository) isDuplicateEmailError(err error) bool {
	var we mongo.WriteException
	if errors.As(err, &we) {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 || e.Code == 11001 {
				if strings.Contains(e.Message, "index: email_1") || strings.Contains(e.Message, "email dup key") {
					return true
				}
			}
		}
	}
	return false
}

func NewUsersRepository(database *mongo.Database, collectionName string) UserRepository {
	return &userRepository{
		collection: database.Collection(collectionName),
	}
}
