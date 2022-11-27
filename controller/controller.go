package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sp-rahul/mongoapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://sp:1234@cluster0.yvyzbhn.mongodb.net/?retryWrites=true&w=majority"
const dbName = "netflix"
const colName = "watchlist"

//most important
var collection *mongo.Collection

//connect with mongo db
func init() {
	//client option
	clientOption := options.Client().ApplyURI(connectionString)

	//connect{ion to mongo
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")

	//collection = (mongo.Collection)(client).Database(dbName).Collection(colName)
	collection = client.Database(dbName).Collection(colName)
	// collection instance
	fmt.Println("Collection instance is ready")

}

// mngoDB helper file

// insert 1 record

func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted 1 movie in bd with id ", inserted.InsertedID)

}

// update onle record
func updateOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count", result.ModifiedCount)
}

//delete 1 record
func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Movie got with delete count : ", deleteCount)

}

//delete all record from mongoDB
func deleteAllMovie() int64 {
	filter := bson.D{{}}
	deleteResult, err := collection.DeleteMany(context.Background(), filter, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of movie deleted : ", deleteResult.DeletedCount)

	return deleteResult.DeletedCount

}

//  get all from database
func getAllMovies() []primitive.M {
	curr, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var movies []primitive.M
	for curr.Next(context.Background()) {
		var movie bson.M
		err := curr.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)

	}
	defer curr.Close(context.Background())
	return movies

}

// Actual collector ~ file
func GetMyAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)

}
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contest-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)

}
func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contest-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods", "POST")

	params := mux.Vars(r)
	updateOneMovie("id")
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contest-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}
func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contest-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods", "DELETE")

	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)

}
