package user

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	mongodb "samgates.io/server/database"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}

func Find(username string) (User, error) {
	MongoDatabase, MongoCtx := mongodb.GetMongoConnection()
	userCollection := MongoDatabase.Collection("users")

	var user User
	err := userCollection.FindOne(MongoCtx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func Authenticate(username string, password string) (User, error) {
	user, err := Find(username)

	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func Create(username string, password string) error {
	MongoDatabase, MongoCtx := mongodb.GetMongoConnection()
	userCollection := MongoDatabase.Collection("users")

	hashedPassword, err := hashPassword(password)

	if err != nil {
		return err
	}

	doc := bson.D{{Key: "username", Value: username}, {Key: "password", Value: hashedPassword}}

	result, err := userCollection.InsertOne(MongoCtx, doc)
	_ = result

	if err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
