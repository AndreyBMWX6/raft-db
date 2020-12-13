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
	"strconv"
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
		log.Println("Car already exists:", car)
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
	response := message.DBResponse{
		Owner:    r.Owner,
		Result: nil,
		Status:   0,
		Entry:  nil,
	}

	var path = r.Request.URL.Path[1:]

	id, err := strconv.Atoi(path)
	if err != nil {
		log.Print("invalid path")

		response.Status = http.StatusBadRequest
		response.Result = []byte(
			`{ "error" : "invalid path" }`,
		)
		return &response
	}

	var filter = bson.D{{"id", id}}

	// this logic will be added instead of passing filter in path with fix of routing routes with parameters
	/*
	err := r.Request.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

	params := r.Request.Form

	for param, value :=range params {
		_, err := dbm.Collection.Find(context.TODO(), bson.M{param: bson.M{"$exists": true}})
		if err != nil {
			log.Println(err, "error: invalid parameters")
			response.Status = http.StatusBadRequest
			response.Result = []byte(
				`{ "error" : "invalid parameters" }`,
			)
			return &response
		}
		if param == "id" || param == "number" {
			val, err := strconv.Atoi(value[0])
			if err != nil {
				log.Fatal(err)
			}
			filter = append(bson.D{}, bson.E{param, val})
		} else {
			filter = append(filter, bson.E{param, value})
		}
	}
*/


	cur, err := dbm.Collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	if cur.RemainingBatchLength() == 0 {
		log.Println("Student not found: filter:", filter)
		response.Status = http.StatusNotFound
		response.Result = []byte(
			`{ "error" : "Student not found" }`,
		)
		return &response
	} else {
		var results []Car

		cur, err := dbm.Collection.Find(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}

		for cur.Next(context.TODO()) {
			// create a value into which the single document can be decoded
			var elem Car
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}

			results = append(results, elem)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		// Close the cursor once finished
		err = cur.Close(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		type GetResponse struct {
			Result []Car `json:"result"`
		}
		resp, _ := json.Marshal(GetResponse{Result: results})

		response.Status = http.StatusOK
		response.Result = resp
		return &response
	}
}

func(dbm *MongoDBManager) Update(r message.DBRequest) message.DBMessage {
	response := message.DBResponse{
		Owner:    r.Owner,
		Result: nil,
		Status:   0,
		Entry:  nil,
	}

	var path = r.Request.URL.Path[1:]

	id, err := strconv.Atoi(path)
	if err != nil {
		log.Print("invalid path")

		response.Status = http.StatusBadRequest
		response.Result = []byte(
			`{ "error" : "invalid path" }`,
		)
		return &response
	}

	filter := bson.D{{"id", id}}

	cur, err := dbm.Collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	if cur.RemainingBatchLength() == 0 {
		log.Println("Student not found: filter:", filter)
		response.Status = http.StatusNotFound
		response.Result = []byte(
			`{ "error" : "Student not found" }`,
		)
		return &response
	} else {
		// will use more update methods later now only set, that updates field value
		type PutForm struct {
			Updates map[string]string `json:"updates"`
		}
		decoder := json.NewDecoder(r.Request.Body)
		var gotUpdates PutForm
		if err := decoder.Decode(&gotUpdates); err != nil {
			log.Print("invalid body")

			response.Status = http.StatusBadRequest
			response.Result = []byte(
				`{ "error" : "invalid body" }`,
			)
			return &response
		} else {
			update := bson.D{}
			for field, value := range gotUpdates.Updates {
				update = append(update, bson.E{"$set", bson.D{
					{field, value},
				}})
			}

			_, err := dbm.Collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Car with id", id, "updated")

			response.Status = http.StatusOK
			response.Result = []byte(
				`{ "result" : "Car updated" }`,
			)
			return &response
		}
	}
}

func(dbm *MongoDBManager) Delete(r message.DBRequest) message.DBMessage {
	response := message.DBResponse{
		Owner:    r.Owner,
		Result: nil,
		Status:   0,
		Entry:  nil,
	}

	var path = r.Request.URL.Path[1:]

	id, err := strconv.Atoi(path)
	if err != nil {
		log.Print("invalid path")

		response.Status = http.StatusBadRequest
		response.Result = []byte(
			`{ "error" : "invalid path" }`,
		)
		return &response
	}

	filter := bson.D{{"id", id}}

	cur, err := dbm.Collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	if cur.RemainingBatchLength() == 0 {
		log.Println("Student not found: filter:", filter)
		response.Status = http.StatusNotFound
		response.Result = []byte(
			`{ "error" : "Student not found" }`,
		)
		return &response
	} else {
		_, err := dbm.Collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Car with id", id, "deleted")

		response.Status = http.StatusOK
		response.Result = []byte(
			`{ "result" : "Car deleted" }`,
		)
		return &response
	}

	return &message.DBResponse{}
}

func (dbm *MongoDBManager) ApplyDBMsg(msg message.DBMessage) {
	switch req := msg.(type) {
	case *message.DBRequest:
		switch req.Type {
		case message.PostRequestType:
			dbm.DBOut <- dbm.Post(*req)
		case message.GetRequestType:
			dbm.DBOut <- dbm.Get(*req)
		case message.PutRequestType:
			dbm.DBOut <- dbm.Update(*req)
		case message.DeleteRequestType:
			dbm.DBOut <- dbm.Delete(*req)
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
