package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {
	e := echo.New()
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUserById)
	e.POST("/users/add", addUser)
	e.DELETE("/users/:id", deleteUser)
	e.PUT("/users/:id", editUser)
	e.Logger.Fatal(e.Start(":8000"))

}

type User struct {
	Id    int    `json:"id" form:"id" query:"id" db:"id"`
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

type JsonFile struct {
	Id     int    `json:"id"`
	Users  []User `json:"users"`
}

func initJSON() JsonFile {
	file, _ := ioutil.ReadFile("users.json")
	data := JsonFile{}
	_ = json.Unmarshal(file, &data)
	return data
}

func getUsers(c echo.Context) error {
	data := initJSON()
	return c.JSON(http.StatusOK, data.Users)
}

func getUserById(c echo.Context) error {
	data := initJSON()
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, data.Users[id-1])
}

func addUser(c echo.Context) error {
	data := initJSON()
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	user.Id = data.Id + 1
	data.Id += 1
	data.Users = append(data.Users, *user)
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("users.json", file, 0644)
	return c.JSON(http.StatusOK, user.Id)
}

func deleteUser(c echo.Context) error {
	data := initJSON()
	id, _ := strconv.Atoi(c.Param("id"))
	data.Users = append(data.Users[:id-1], data.Users[id:]...)
	data.Id -= 1
	for i := id - 1; i < len(data.Users); i++ {
		data.Users[i].Id -= 1
	}
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("users.json", file, 0644)
	return c.JSON(http.StatusOK, "Success")
}

func editUser(c echo.Context) error {
	data := initJSON()
	id, _ := strconv.Atoi(c.Param("id"))
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	user.Id = id
	data.Users[id - 1] = *user
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("users.json", file, 0644)
	return c.JSON(http.StatusOK, id)
}
