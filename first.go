package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//type JSONTime time.Time

// func (t JSONTime) MarshalJSON() ([]byte, error) {
// 	//do your serializing here
// 	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
// 	return []byte(stamp), nil
//}

type Article struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title,omitempty"`
	SubTitle string             `json:"subtitle" bson:"subtitle,omitempty"`
	Content  string             `json:"content" bson:"content,omitempty"`
	Creation time.Time          `json:"creation" bson:"creation,omitempty"`
}

//Connection mongoDB
var collection = ConnectDB()

func main() {
	//Init Router
	r := mux.NewRouter()

	// arrange our route
	r.HandleFunc("/articles", getArticles).Methods("GET")
	r.HandleFunc("/articles", createArticle).Methods("POST")
	r.HandleFunc("/articles/{id}", getArticle).Methods("GET")
	r.HandleFunc("/articles/search?q=title", getArticle).Methods("GET")

	// set our port address
	log.Fatal(http.ListenAndServe(":8000", r))
}

func getArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Book array
	var Articles []Article

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var article Article
		// & character returns the memory address of the following variable.
		err := cur.Decode(&article) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		Articles = append(Articles, article)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(Articles) // encode similar to serialize process.
}

// ConnectDB : This is helper function to connect mongoDB
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func ConnectDB() *mongo.Collection {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://hariharan:hariharan@cluster0.p7nty.mongodb.net/<dbname>?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("go_rest_api").Collection("Articles")

	return collection
}

// ErrorResponse : This is error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError : This is helper function to prepare error model.
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var article Article

	// decode the body request params
	_ = json.NewDecoder(r.Body).Decode(&article)

	// insert our book model.
	result, err := collection.InsertOne(context.TODO(), article)

	if err != nil {
		GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}
func getArticle(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var article Article
	// we get params with mux.
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// create filter. If it is unnecessary to sort data
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&article)

	if err != nil {
		GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(article)
}

func searchArticle(w http.ResponseWriter, r *http.Request) {
	var article Article

	title := string(r.URL.Query().Get("q"))

	filter := bson.M{"title": title}
	err := collection.FindOne(context.TODO(), filter).Decode(&article)
	if err != nil {
		GetError(err, w)
		return
	}
	// Use the fmt.Fprintf() function to interpolate the id value with our response
	// and write it to the http.ResponseWriter.
	json.NewEncoder(w).Encode(article)
}
