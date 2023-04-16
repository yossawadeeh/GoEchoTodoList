package main

import (
	"echo-todolist/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func main() {
	Db, err := gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/todolist?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	Db.AutoMigrate(&model.User{}, &model.ToDoList{}, &model.Status{})

	statusToDo := model.Status{StatusName: "To Do"}
	Db.Create(&statusToDo)
	statusPending := model.Status{StatusName: "Pending"}
	Db.Create(&statusPending)
	statusDone := model.Status{StatusName: "Done"}
	Db.Create(&statusDone)

	Db.Model(&model.User{}).Create([]map[string]interface{}{
		{"UserName": "yeolowbatt", "FirstName": "Yossawadee", "LastName": "Hoymala", "Email": "yeolowbatt@gmail.com", "Password": "123"},
		{"UserName": "pongkimuji", "FirstName": "Pong", "LastName": "Kimuji", "Email": "pongkimuji@gmail.com", "Password": "123"},
	})

}
