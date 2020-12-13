package manager

import (

	"../message"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type MongoDBManager struct {
	DBIn  <-chan message.DBMessage
	DBOut chan<- message.DBMessage

	DataBase *mongo.Database
	Collection *mongo.Collection
}

type Specification struct {
	Speed uint32				`json:"speed"`
	Color string				`json:"color"`
	SpecialAbility string		`json:"special ability"`
}

type Car struct {
	Id uint32   				`json:"id"`
	Name string                 `json:"name"`
	Rarity string				`json:"rarity"`
	Number uint32               `json:"number"`
	Specification Specification `json:"specification"`
	Catchword string			`json:"catchword"`
}

func(dbm *MongoDBManager) Post(r message.DBRequest) message.DBMessage {
	response := message.DBResponse{
		Owner:    r.Owner,
		Result: nil,
		Status:   0,
		Entry:  nil,
	}

	bytes, _ := ioutil.ReadAll(r.Request.Body)

	var car Car
	if err := json.Unmarshal(bytes, &car); err != nil {
		log.Print("invalid body")
		response.Status = http.StatusBadRequest
		response.Result = []byte(
			`{ "error" : "invalid body" }`,
		)
		response.Entry = nil
		return &response
	}

	cur, err := dbm.Collection.Find(context.TODO(), bson.D{{"id", car.Id}}, options.Find().SetLimit(1))
	if err != nil {
		log.Fatal(err)
	}


	if cur.RemainingBatchLength() != 0 {
		log.Printf("Car already exists: %s", car)
		response.Status = http.StatusConflict
		response.Result = []byte(
			`{ "error" : "already exists" }`,
		)
		response.Entry = nil
		return &response
	}

	insertResult, err := dbm.Collection.InsertOne(context.TODO(), car)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted a single document: ", insertResult.InsertedID)
	response.Status = http.StatusOK
	response.Result = []byte(
		`{ "result" : "Car posted" }`,
	)
	query := []byte("Post\n")
	query = append(query, bytes...)
	response.Entry = &message.Entry{
		Term:  0,
		Query: query,
	}
	return &response
}

func(dbm *MongoDBManager) Get(r message.DBRequest) message.DBMessage {
return &message.DBResponse{}
}

func(dbm *MongoDBManager) Update(r message.DBRequest) message.DBMessage {
	return &message.DBResponse{}
}

func(dbm *MongoDBManager) Delete(r message.DBRequest) message.DBMessage {
	return &message.DBResponse{}
}

func (dbm *MongoDBManager) ApplyDBMsg(msg message.DBMessage) {
	switch req := msg.(type) {
	case *message.DBRequest:
		switch req.Type {
		case message.PostRequestType:
			dbm.DBOut <- dbm.Post(*req)
		case message.GetRequestType:
			dbm.DBOut <- dbm.Post(*req)
		case message.PutRequestType:
			dbm.DBOut <- dbm.Post(*req)
		case message.DeleteRequestType:
			dbm.DBOut <- dbm.Post(*req)
		}
	default:
		log.Print("`DBRequestMessage` expected, got another type")
	}
	return
}

func (dbm *MongoDBManager) ProcessMessage() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://AndreyBMWX6:rAft-db)@cluster0.fpnaa.mongodb.net/transport?retryWrites=true&w=majority",
	))
	if err != nil { log.Fatal(err) }

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	dbm.DataBase = client.Database("transport")
	dbm.Collection = dbm.DataBase.Collection("cars")

	for {
		select {
		case msg := <-dbm.DBIn:
			dbm.ApplyDBMsg(msg)
		default:
		}
	}
}
