package main

import (
	"crowdfunding/auth"
	"crowdfunding/handler"
	"crowdfunding/helper"
	"crowdfunding/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	authService := auth.NewService()

	// token, err := authService.ValidateToken("")
	// if err != nil {
	// 	fmt.Println("ERROR")
	// 	fmt.Println("ERROR")
	// 	fmt.Println("ERROR")
	// }
	// if token.Valid {
	// 	fmt.Println("VALID")
	// } else {
	// 	fmt.Println("INVALID")
	// }

	userHandler := handler.NewUserHandler(userService, authService)

	//fmt.Println(authService.GenerateToken(100))

	// token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxN30.VwrFwsrSR5kFpprxLNRubrnEOn7Uqy1GJaxAlZk1aOo")
	// if err != nil {
	// 	fmt.Println("ERROR")
	// 	fmt.Println("ERROR")
	// 	fmt.Println("ERROR")
	// }

	// if token.Valid {
	// 	fmt.Println("VALID")
	// 	fmt.Println("VALID")
	// 	fmt.Println("VALID")
	// } else {
	// 	fmt.Println("INVALID")
	// 	fmt.Println("INVALID")
	// 	fmt.Println("INVALID")
	// }

	router := gin.Default()
	api := router.Group("api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadedAvatar)

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

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized 1", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized 2", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized 3", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized 4", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}

}
