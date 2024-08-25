// Ref: https://github.com/gin-gonic/gin
package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
    ID   string `gorm:"primary_key"`
    Name string
	CreatedAt time.Time
}


func dbConnect() *gorm.DB {
	dsn := "local:password@tcp(mysql:3306)/go_tutrial?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

type UserRegisterRequest struct {
	ID string `json:"id"`
	Name string `json:"name"`
}


type UserUpdateRequest struct {
	Name string `json:"name"`
}

func main() {
	db := dbConnect()
	
	
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
    	c.JSON(http.StatusOK, gin.H{
      	"message": "pong",
    	})
  	})

	// READ
	r.GET("/users/:id", func(c *gin.Context) {
		user := &User{}
		id := c.Param("id")
		result := db.First(&user, "ID = ?", id)

		if result.Error != nil {
			c.String(http.StatusNotFound, "User not found")
		}

		c.JSON(http.StatusOK, user)
  	})

	// CREATE
	r.POST("/users", func(c *gin.Context) {
		userRegisterRequest := &UserRegisterRequest{}

		if err := c.BindJSON(&userRegisterRequest); err != nil {
			c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
		}
		
		user := &User{ID: userRegisterRequest.ID, Name: userRegisterRequest.Name}

		db.Create(&user)
		c.JSON(http.StatusOK, user)
  	})


	//UPDATE
	r.PUT("/users/:id", func(c *gin.Context) {
		userUpdateRequest := &UserUpdateRequest{}

		if err := c.BindJSON(&userUpdateRequest); err != nil {
			c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
		}

		user := &User{}
		id := c.Param("id")
		result := db.First(&user, "ID = ?", id)
		if result.Error != nil {
			c.String(http.StatusNotFound, "User not found")
		}

		db.Model(&user).Update("Name", userUpdateRequest.Name)

		c.JSON(http.StatusOK, user)
	})
		
	//DELETE
	r.DELETE("/users/:id", func(c *gin.Context) {
		user := &User{}
		id := c.Param("id")

		db.Where("id = ?", id).Delete(&user)
		c.Status(http.StatusNoContent)
	})

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}