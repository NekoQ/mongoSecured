package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	// "golang.org/x/crypto/bcrypt"
)

type User struct {
	Login    string `bson:"login"`
	Password string `bson:"password"`
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func secureUser(user User) User {
	var securedUser User
	securedUser.Login = user.Login
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	check(err)
	securedUser.Password = string(password)

	return securedUser
}

func main() {

	fmt.Println("Connected to the Database: secred, Collection: users")
	var user User
	fmt.Println("Create new user:")
	fmt.Println("Login:")
	fmt.Scanln(&user.Login)
	fmt.Println("Password:")
	fmt.Scanln(&user.Password)
	/*
	   Connect to my cluster
	*/
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://test:zAAjgMQ4vZlXEuYu@cluster0.ke1qd.mongodb.net/secured?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	collection := client.Database("secured").Collection("users")

	secured := secureUser(user)
	_, err = collection.InsertOne(ctx, secured)
	check(err)
	fmt.Printf("\nStored in database\nLogin: %v\nPassword: %v\n", secured.Login, secured.Password)

	for {
		var newUser User
	Login:
		fmt.Println("\nLogin:")
		fmt.Scanln(&newUser.Login)
		filter := bson.D{{"login", newUser.Login}}
		err = collection.FindOne(ctx, filter).Decode(&newUser)
		if err != nil {
			fmt.Println("Error: No such user")
			goto Login
		}
		var password string
	L:
		fmt.Println("Password: ")
		fmt.Scanln(&password)
		err = bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(password))
		if err == nil {
			fmt.Println("\nSuccessful login!\n ")
			break
		} else {
			fmt.Println("Error: Wrong password")
			goto L
		}

	}
	// cur, err := collection.Find(ctx, bson.D{})
	// check(err)
	// defer cur.Close(ctx)

	// var users []User
	// if err = cur.All(ctx, &users); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(users)
}

// email1, err := bcrypt.GenerateFromPassword([]byte("marius@test.com"), 10)
// check(err)
// password1, err := bcrypt.GenerateFromPassword([]byte("testpassword"), 10)
// check(err)
// email2, err := bcrypt.GenerateFromPassword([]byte("secured@mail.com"), 10)
// check(err)
// password2, err := bcrypt.GenerateFromPassword([]byte("securedpassword"), 10)
// check(err)

// users := []interface{}{
// 	bson.D{{"email", string(email1)}, {"login", "testlogin"}, {"password", string(password1)}},
// 	bson.D{{"email", string(email2)}, {"login", "secured"}, {"password", string(password2)}},
// }

// res, insertErr := collection.InsertMany(ctx, users)
// check(insertErr)
// fmt.Println(res)
