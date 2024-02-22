package storeProduct

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/product"
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
	storeProductCollection string = "storeProduct"
	productCollection      string = "product"
)

type Repository interface {
	GetByCategory(categoryName product.Category) ([]StoreProduct, error)
	Get(idStore string, idProduct string) (StoreProduct, error)
	Delete(idStore string, idProduct string) error
	Post(sp StoreProduct) (StoreProduct, error)
	GetAll() ([]StoreProduct, error)
}
type storeProductRepository struct {
	configMongo config.Mongodb
	client      *mongo.Client
}

func (spr storeProductRepository) Ping() error {
	return spr.client.Ping(context.TODO(), readpref.Primary())
}
func NewRepositoryStoreProduct(mongodb config.Mongodb) (Repository, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+mongodb.Host+":"+strconv.Itoa(mongodb.Port)))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	r := storeProductRepository{
		client:      client,
		configMongo: mongodb,
	}
	return r, nil
}
func (spr storeProductRepository) GetByCategory(categoryName product.Category) ([]StoreProduct, error) {
	storeProductCollection := spr.client.Database(spr.configMongo.Dbname).Collection(storeProductCollection)
	productCollection := spr.client.Database(spr.configMongo.Dbname).Collection(productCollection)

	sp := make([]StoreProduct, 0)
	name := string(categoryName)
	filter := bson.M{
		"category": name,
	}
	productResult, err := productCollection.Find(context.TODO(), filter)
	var products []product.Product
	if err = productResult.All(context.Background(), &products); err != nil {
		return sp, err
	}
	prIds := make([]string, 0)
	for _, p := range products {
		prIds = append(prIds, p.Id)

	}
	productIdFilter := bson.M{"productID": bson.M{"$in": prIds}}
	result, _ := storeProductCollection.Find(context.TODO(), productIdFilter) //result := storeProductCollection.FindOne(context.TODO(), bson.M{"category": categoryName})
	if err = result.All(context.Background(), &sp); err != nil {
		return sp, err
	}
	return sp, nil
}

func (spr storeProductRepository) Get(idStore string, idProduct string) (StoreProduct, error) {
	storeProductCollection := spr.client.Database(spr.configMongo.Dbname).Collection(storeProductCollection)
	sp := StoreProduct{}
	result := storeProductCollection.FindOne(context.TODO(), bson.M{"storeID": idStore, "productID": idProduct})
	err := result.Decode(&sp)
	if err != nil {
		return sp, err
	}
	return sp, nil
}
func (spr storeProductRepository) Post(sp StoreProduct) (StoreProduct, error) {
	storeProductCollection := spr.client.Database(spr.configMongo.Dbname).Collection(storeProductCollection)
	result, err := storeProductCollection.InsertOne(context.TODO(), sp)
	if sp.StoreID == "" {
		sp.StoreID = fmt.Sprintf("%v", xid.New())
	}
	if sp.ProductID == "" {
		sp.ProductID = fmt.Sprintf("%v", xid.New())
	}
	if err != nil {
		return sp, err
	}
	sp.StoreID = fmt.Sprint(result.InsertedID)
	sp.ProductID = fmt.Sprint(result.InsertedID)
	return sp, nil
}
func (spr storeProductRepository) GetAll() ([]StoreProduct, error) {
	storeProductCollection := spr.client.Database(spr.configMongo.Dbname).Collection(storeProductCollection)
	var storesProducts []StoreProduct
	sp, err := storeProductCollection.Find(context.Background(), bson.D{})
	if err != nil {
		fmt.Println(err.Error())
		return storesProducts, err
	}
	if err = sp.All(context.Background(), &storesProducts); err != nil {
		log.Fatal(err)
	}
	return storesProducts, nil
}

func (spr storeProductRepository) Delete(idStore string, idProduct string) error {
	storeProductCollection := spr.client.Database(spr.configMongo.Dbname).Collection(storeProductCollection)
	sp, err := storeProductCollection.DeleteOne(context.TODO(), bson.M{"storeID": idStore, "productID": idProduct})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", sp.DeletedCount)
	return nil
}
