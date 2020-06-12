package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = func() context.Context {
	return context.Background()
}()

type User struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

var user User

//main
func main() {
	route := gin.Default()
	route.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "helloworld",
		})
	})

	route.GET("users", findAll)
	route.GET("/user/:id", find)
	route.POST("/user", insert)
	route.PUT("/user/:id", updateOne)
	route.DELETE("/user/:id", deleteOne)
	route.Run()
}

//connect database
func connect() (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	// fmt.Println("connect database successfully")
	// db := client.Database("belajar_golang")
	return client.Database("belajar_golang"), nil
}

//error handler
func Error(c *gin.Context, err error) bool {
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error(), "data": ""})
		return true
	}
	return false
}

//findAll
func findAll(c *gin.Context) {
	db, err := connect()
	if Error(c, err) {
		return //exit
	}

	cur, err := db.Collection("user").Find(ctx, bson.D{})
	if Error(c, err) {
		return //exit
	}
	defer cur.Close(ctx)

	result := make([]User, 0)
	for cur.Next(ctx) {
		var row User
		err := cur.Decode(&row)
		if Error(c, err) {
			return //exit
		}
		result = append(result, row)
	}

	fmt.Println(result)
	c.JSON(200, gin.H{
		"message": "success",
		"data":    result,
	})

}

//find
func find(c *gin.Context) {
	db, err := connect()
	if Error(c, err) {
		return //exit
	}
	id := c.Param("id")
	fmt.Println(id)

	_id, _ := primitive.ObjectIDFromHex(id)

	var res User
	err = db.Collection("user").FindOne(ctx, bson.M{"_id": _id}).Decode(&res)
	if Error(c, err) {
		return //exit
	}

	fmt.Println(res)
	c.JSON(200, gin.H{
		"message": "success",
		"data":    res,
	})
}

//insert
func insert(c *gin.Context) {
	db, err := connect()
	if Error(c, err) {
		return //exit
	}

	c.BindJSON(&user)
	fmt.Println(user.Name)
	res, err := db.Collection("user").InsertOne(ctx, bson.M{"name": user.Name, "age": user.Age})
	if Error(c, err) {
		return //exit
	}

	_ = db.Collection("user").FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&user)

	c.JSON(200, gin.H{
		"message": "success",
		"data":    user,
	})
}

//updateOne
func updateOne(c *gin.Context) {
	db, err := connect()
	if Error(c, err) {
		return //exit
	}

	c.BindJSON(&user)

	id := c.Param("id")
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}

	_, err = db.Collection("user").UpdateOne(ctx, filter, bson.M{"$set": user})
	if Error(c, err) {
		return //exit
	}
	// var ress User
	_ = db.Collection("user").FindOne(ctx, bson.M{"_id": _id}).Decode(&user)

	c.JSON(200, gin.H{
		"message": "success",
		"data":    user,
	})
}

//deleteOne
func deleteOne(c *gin.Context) {
	db, err := connect()
	if Error(c, err) {
		return //exit
	}

	c.BindJSON(&user)

	id := c.Param("id")
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}

	_, err = db.Collection("user").DeleteOne(ctx, filter)
	if Error(c, err) {
		return //exit
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    gin.H{"_id": id},
	})
}
