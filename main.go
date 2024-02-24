// package main

// import "fmt"

// // when
// // 1. when we need to update state
// // 2. when we want to optimize the memory for large objects that are getting called A lot

// type User struct {
// 	email string
// 	username string
// 	age int
// }

// //method receiver on structures, this is just a syntactic sugar allowed by go lang
// func (u User) Email() string {
// 	return u.email
// }

// // x amount of bytes => sizeof(user)
// func Email(u User) string {
// 	return	u.email
// }

// func (u *User) updateEmail(email string) {
// 	u.email = email
// }

// func main() {
// 	user := User{
// 		email: "ajay@gmail.com",
// 	}

// 	user.updateEmail("ajay2@gmail.com")
// 	fmt.Println(user.Email())
// }

package main

import (
	"fmt"
	"log"

	"main.go/api"
)

func main() {

	store, err := api.NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":3000", store)
	server.Run()
}
