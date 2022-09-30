package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/athosone/projectraven/tracking/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	UserCollectionName = "users"
)

type userDao struct {
	ID    string `bson:"_id"`
	Email string `bson:"name"`
}

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) (domain.UserRepository, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// err := createUserIndexes(ctx, db.Collection(UserCollectionName))
	// if err != nil {
	// 	return nil, err
	// }
	return &userRepo{
		collection: db.Collection(UserCollectionName),
	}, nil
}

func (r *userRepo) FindByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindById(id uuid.UUID) (*domain.User, error) {
	var userDb userDao
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&userDb)
	if err != nil {
		return nil, err
	}
	return mapUserDbToDomain(&userDb)
}

func (r *userRepo) FindByIdpProviderAndIdpId(idpProvider string, idpId string) (*domain.User, error) {
	var userDb userDao
	err := r.collection.FindOne(context.Background(), bson.M{"idpProvider": idpProvider, "idpId": idpId}).Decode(&userDb)
	if err != nil {
		return nil, err
	}
	return mapUserDbToDomain(&userDb)
}

func (r *userRepo) Save(user *domain.User) error {
	userDb, err := mapUserDomainToDb(user)
	if err != nil {
		return err
	}
	_, err = r.collection.InsertOne(context.Background(), userDb)
	if err != nil {
		return err
	}
	return nil
}

func createUserIndexes(ctx context.Context, collection *mongo.Collection) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.M{
				"idpProvider": 1,
				"idpId":       1,
			},
		},
		{
			Keys: bson.M{
				"email": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

func mapUserDomainToDb(user *domain.User) (*userDao, error) {
	id := user.ID.String()
	return &userDao{
		ID:    id,
		Email: user.Email,
	}, nil
}

func mapUserDbToDomain(userDb *userDao) (*domain.User, error) {
	userId, err := uuid.Parse(userDb.ID)
	if err != nil {
		return nil, fmt.Errorf("map to domain user: %w", err)
	}
	return &domain.User{
		ID:    userId,
		Email: userDb.Email,
	}, nil
}
