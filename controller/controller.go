package controller

import (
	"context"
	model "dbs/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const connnectionString = "mongodb+srv://@gocluster.gdlendu.mongodb.net/?retryWrites=true&w=majority"

const dbName = "netflix"
const colName = "watchlist"

// most important
var collection *mongo.Collection

// connect with mongodb
func init() {
	//only runs once at start
	//client options
	connectionString := os.Getenv("MONGODB_URI")
	if connectionString == "" {
		fmt.Println("MONGODB_URI environment variable not set")
		return
	}
	clientOption := options.Client().ApplyURI(connectionString) //general database connection

	//connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOption) //returns client
	//context todo: when calling db outside. provide a context how long connect is stablished n all
	//what happens when disconnected
	//context.todo no idea which context to use.

	if err != nil {
		fmt.Println("nigg")
		log.Fatal(err)
	}
	fmt.Println("mongo connection stablished")

	collection = client.Database(dbName).Collection(colName)

	//collection reference
	fmt.Println("collection ref is running")

}

//basic method
//check data
//helper: take data will insert and respnose if success or not

//MONGODB helpers -file
//insert 1 record

func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie) //pass conetxt when doing db operation.
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted one movie : ", inserted.InsertedID)
}

//to update databse(mongo) provide a filter acc to that db will be updated.
//pass a flag (set)

//update 1 record

func updateOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId) //objidfromhex converts string to id acc to mongo
	//mongo understands _id
	//find the value (mongo finds)
	//mongo gives bson.

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}} //always setting watched to true.

	result, err := collection.UpdateOne(context.Background(), filter, update) //returns how many values are updated.
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count:", result.ModifiedCount)
}

// to delete in mongo: delete one, delete all
func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("movie got deleted with dltCount: ", deleteCount)
}

// delete all records from mongodb
func deleteAllMovie() int64 {
	filter := bson.D{{}} //not providing anything so everything should be selected
	deleteResult, err := collection.DeleteMany(context.Background(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("nu of movies deleted: ", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

// get all movies form database
func getAllMovies() []primitive.M {
	//episodeCollection gives all data from mongo
	//stored in cursor
	//loop through cursor
	//bson[].m will be type of result.
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	//declare var arr
	var movies []primitive.M //using bson.M gives errors.

	//have cursor and var
	//loop through
	for cursor.Next(context.Background()) { //pass context
		var movie bson.M
		err := cursor.Decode(&movie) //either decode and movie will be fullfilled with values
		//decode value- use my struct to decode -in the case var
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	defer cursor.Close(context.Background())

	return movies
}

// Actual controller -file
func GetMyAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allMovies := getAllMovies()
	//send json response
	json.NewEncoder(w).Encode(allMovies)
}

// create a movie
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode") //property of set
	w.Header().Set("Allow-Control-Allow-Methods", "POST")              //what type of values are allowed
	//routers also allows to set the headers.

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)

}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode") //property of set
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	//need unique id of movie
	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

// delete one movie
func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode") //property of set
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode") //property of set
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)
}
