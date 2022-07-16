package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.TODO()
	opts := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	defer client.Disconnect(ctx)
	fmt.Printf("%T\n", client)

	testDB := client.Database("test")
	fmt.Printf("%T\n", testDB)

	mongoCollection := testDB.Collection("mongo")
	defer mongoCollection.Drop(ctx)

	fmt.Printf("%T\n", mongoCollection)
	example := bson.D{
		{"someString", "Example String"},
		{"someInteger", 12},
		{"someStringSlice", []string{"Example 1", "Example 2", "Example 3"}},
	}
	r, err := mongoCollection.InsertOne(ctx, example)
	if err != nil {
		panic(err)
	}
	fmt.Println(r.InsertedID)

	examples := []interface{}{
		bson.D{
			{"someString", "Second Example String"},
			{"someInteger", 253},
			{"someStringSlice", []string{"Example 15", "Example 42", "Example 83", "Example 5"}},
		},
		bson.D{
			{"someString", "Another Example String"},
			{"someInteger", 54},
			{"someStringSlice", []string{"Example 21", "Example 53"}},
		},
	}

	rs, err := mongoCollection.InsertMany(ctx, examples)

	if err != nil {
		panic(err)
	}

	fmt.Println(rs.InsertedIDs)

	c := mongoCollection.FindOne(ctx, bson.M{"_id": r.InsertedID})

	var exampleResult bson.M
	c.Decode(&exampleResult)

	fmt.Printf("\nItem with ID: %v contains the following:\n", exampleResult["_id"])
	fmt.Println("someString:", exampleResult["someString"])
	fmt.Println("someInteger:", exampleResult["someInteger"])
	fmt.Println("someStringSlice:", exampleResult["someStringSlice"])

	filter := bson.D{{"someInteger", bson.D{{"$lt", 60}}}}
	examplesGT50, err := mongoCollection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	var examplesResult []bson.M
	if err = examplesGT50.All(ctx, &examplesResult); err != nil {
		panic(err)
	}

	for _, e := range examplesResult {
		fmt.Printf("\nItem with ID: %v contains the following:\n", e["_id"])
		fmt.Println("someString:", e["someString"])
		fmt.Println("someInteger:", e["someInteger"])
		fmt.Println("someStringSlice:", e["someStringSlice"])
	}

	rUpdt, err := mongoCollection.UpdateByID(
		ctx,
		r.InsertedID,
		bson.D{
			{"$set", bson.M{"someInteger": 201}},
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of items updated:", rUpdt.ModifiedCount)

	c2 := mongoCollection.FindOne(ctx, bson.M{"_id": r.InsertedID})

	var exampleResult2 bson.M
	c2.Decode(&exampleResult2)

	fmt.Printf("\nItem with ID: %v contains the following:\n", exampleResult2["_id"])
	fmt.Println("someString:", exampleResult2["someString"])
	fmt.Println("someInteger:", exampleResult2["someInteger"])
	fmt.Println("someStringSlice:", exampleResult2["someStringSlice"])

	rUpdt, err = mongoCollection.UpdateOne(
		ctx,
		bson.M{"_id": r.InsertedID},
		bson.D{
			{"$set", bson.M{"someString": "The Updated String"}},
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of items updated:", rUpdt.ModifiedCount)

	c3 := mongoCollection.FindOne(ctx, bson.M{"_id": r.InsertedID})

	var exampleResult3 bson.M
	c3.Decode(&exampleResult3)

	fmt.Printf("\nItem with ID: %v contains the following:\n", exampleResult3["_id"])
	fmt.Println("someString:", exampleResult3["someString"])
	fmt.Println("someInteger:", exampleResult3["someInteger"])
	fmt.Println("someStringSlice:", exampleResult3["someStringSlice"])

	rUpdt2, err := mongoCollection.UpdateMany(
		ctx,
		bson.D{{"someInteger", bson.D{{"$gt", 60}}}},
		bson.D{
			{"$set", bson.M{"someInteger": 60}},
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of items updated:", rUpdt2.ModifiedCount)

	examplesAll, err := mongoCollection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	var examplesResult2 []bson.M
	if err = examplesAll.All(ctx, &examplesResult2); err != nil {
		panic(err)
	}

	for _, e := range examplesResult2 {
		fmt.Printf("\nItem with ID: %v contains the following:\n", e["_id"])
		fmt.Println("someString:", e["someString"])
		fmt.Println("someInteger:", e["someInteger"])
		fmt.Println("someStringSlice:", e["someStringSlice"])
	}

	rDel, err := mongoCollection.DeleteOne(ctx, bson.M{"_id": r.InsertedID})

	if err != nil {
		panic(err)
	}

	fmt.Println("Number of items deleted:", rDel.DeletedCount)

	time.Sleep(20 * time.Second)

}
