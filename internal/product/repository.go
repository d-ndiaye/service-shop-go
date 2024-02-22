package product

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

const (
	productCollection string = "product"
)

type Repository interface {
	Get(id string) (Product, error)
	Post(p Product) (Product, error)
	GetAll() ([]Product, error)
	GetAllByName(name string) ([]Product, error)
	Delete(id string) error
	Ping() error
}
type productRepository struct {
	configMongo config.Mongodb
	client      *mongo.Client
}

func (pr productRepository) Ping() error {
	return pr.client.Ping(context.TODO(), readpref.Primary())
}

func NewRepository(mongodb config.Mongodb) (Repository, error) {
	fmt.Println("mongodb://" + mongodb.Host + ":" + strconv.Itoa(mongodb.Port))
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+mongodb.Host+":"+strconv.Itoa(mongodb.Port)))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	r := productRepository{
		client:      client,
		configMongo: mongodb,
	}
	return r, nil
}

func (pr productRepository) Get(id string) (Product, error) {
	productCollection := pr.client.Database(pr.configMongo.Dbname).Collection(productCollection)
	p := Product{}
	result := productCollection.FindOne(context.TODO(), bson.M{"_id": id})
	err := result.Decode(&p)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (pr productRepository) Post(p Product) (Product, error) {
	productCollection := pr.client.Database(pr.configMongo.Dbname).Collection(productCollection)
	AddTva(&p)
	result, err := productCollection.InsertOne(context.TODO(), p)
	if p.Id == "" {
		p.Id = fmt.Sprintf("%v", xid.New())
	}
	if err != nil {
		return p, err
	}
	p.Id = fmt.Sprint(result.InsertedID)
	return p, nil
}
func (pr productRepository) GetAll() ([]Product, error) {
	productCollection := pr.client.Database(pr.configMongo.Dbname).Collection(productCollection)
	var products []Product
	p, err := productCollection.Find(context.Background(), bson.D{})
	if err != nil {
		fmt.Println(err.Error())
		return products, err
	}
	if err = p.All(context.Background(), &products); err != nil {
		log.Fatal(err)
	}
	return products, nil
}

func (pr productRepository) GetAllByName(name string) ([]Product, error) {
	productCollection := pr.client.Database(pr.configMongo.Dbname).Collection(productCollection)
	var products []Product
	p, err := productCollection.Find(context.Background(), bson.M{"name": name})
	if err != nil {
		fmt.Println(err.Error())
		return products, err
	}
	if err = p.All(context.Background(), &products); err != nil {
		log.Fatal(err)
	}
	return products, nil
}

func (pr productRepository) Delete(id string) error {
	productCollection := pr.client.Database(pr.configMongo.Dbname).Collection(productCollection)
	p, err := productCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", p.DeletedCount)
	return nil
}
