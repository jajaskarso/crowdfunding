package main

import (
	"crowdfunding/handler"
	"crowdfunding/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details

	// fmt.Println("Connection worked. Good job!")

	// var users []user.User
	// length := len(users)

	// fmt.Println(length)

	// db.Find(&users)

	// length = len(users)
	// fmt.Println(length)

	// for _, user := range users {
	// 	fmt.Println(user.ID)
	// 	fmt.Println(user.Email)
	// 	fmt.Println(user.Name)
	// }

	dsn := "root:root@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)

	router.Run()
	// userInput := user.RegisterUserInput{}
	// userInput.Name = "Tes dari service"
	// userInput.Email = "dariservice@gmail.com"
	// userInput.Occupation = "Tukang Pos"
	// userInput.Password = "password"

	// userService.RegisterUser(userInput)

	// user := user.User{
	// 	Name: "Test",
	// }

	// userRepository.Save(user)
}
