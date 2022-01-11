package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Products struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Productname string             `json:"productname,omitempty" bson:"productname,omitempty"`
}

var collection = connectdb()

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var products []Products
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
		// return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var product Products
		err := cur.Decode(&product)
		if err != nil {
			log.Fatal(err)
		}

		products = append(products, product)
	}

	json.NewEncoder(w).Encode(products)

}
func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product Products

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&product)

	if err != nil {
		log.Fatal(err)
		// return
	}

	json.NewEncoder(w).Encode(product)
}
func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product Products

	_ = json.NewDecoder(r.Body).Decode(&product)
	result, err := collection.InsertOne(context.TODO(), product)

	if err != nil {
		log.Fatal(err)
		// return
	}

	json.NewEncoder(w).Encode(result)
}
func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	var product Products

	filter := bson.M{"_id": id}

	_ = json.NewDecoder(r.Body).Decode(&product)

	update := bson.D{
		{"$set", bson.D{
			{"productname", product.Productname},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&product)

	if err != nil {
		log.Fatal(err)
		// return
	}

	product.ID = id

	json.NewEncoder(w).Encode(product)
}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Fatal(err, w)
		// return
	}

	json.NewEncoder(w).Encode(deleteResult)
}
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))

}
func connectdb() *mongo.Collection {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("go_crud_api").Collection("products")
	return collection
}
