package store

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/config"
	"context"
	"fmt"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"strconv"
)

const storeCollection string = "store"

type Repository interface {
	Get(id string) (Store, error)
	Post(p Store) (Store, error)
	GetAll() ([]Store, error)
	Delete(id string) error
	Ping() error
}
type storeRepository struct {
	configMongo config.Mongodb
	client      *mongo.Client
}

func (sr storeRepository) Ping() error {
	return sr.client.Ping(context.TODO(), readpref.Primary())
}

func NewRepositoryStore(mongodb config.Mongodb) (Repository, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+mongodb.Host+":"+strconv.Itoa(mongodb.Port)))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	r := storeRepository{
		client:      client,
		configMongo: mongodb,
	}
	return r, nil
}

func (sr storeRepository) Get(id string) (Store, error) {
	productCollection := sr.client.Database(sr.configMongo.Dbname).Collection(storeCollection)
	s := Store{}
	result := productCollection.FindOne(context.TODO(), bson.M{"_id": id})
	err := result.Decode(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}
func (sr storeRepository) Post(s Store) (Store, error) {
	storeCollection := sr.client.Database(sr.configMongo.Dbname).Collection(storeCollection)
	result, err := storeCollection.InsertOne(context.TODO(), s)
	if s.Id == "" {
		s.Id = fmt.Sprintf("%v", xid.New())
	}
	if err != nil {
		return s, err
	}
	s.Id = fmt.Sprint(result.InsertedID)
	return s, nil
}
func (sr storeRepository) GetAll() ([]Store, error) {
	storeCollection := sr.client.Database(sr.configMongo.Dbname).Collection(storeCollection)
	var stores []Store
	p, err := storeCollection.Find(context.Background(), bson.D{})
	if err != nil {
		fmt.Println(err.Error())
		return stores, err
	}
	if err = p.All(context.Background(), &stores); err != nil {
		log.Fatal(err)
	}

	return stores, nil
}

func (sr storeRepository) Delete(id string) error {
	storeCollection := sr.client.Database(sr.configMongo.Dbname).Collection(storeCollection)
	s, err := storeCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", s.DeletedCount)
	return nil
}
