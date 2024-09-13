package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func Connection() (*sqlx.DB, error) {
	dsn := "host=localhost port=5432 user=postgres password=123321 dbname=replica sslmode=disable"
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return nil, err
	}

	// Check the database connection
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		db.Close()
		return nil, err
	}

	return db, nil
}

type User struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Age   int    `json:"age" db:"age"`
}

type Methods struct {
	db *sqlx.DB
}

func main() {
	db, err := Connection()
	if err != nil {
		log.Fatal(err)
	}

	m := &Methods{db: db}

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"pong": "success"}) })
	router.POST("/register", m.Register)
	router.GET("/get/:id", m.GetByID)
	router.GET("/get", m.GetAll)

	log.Fatal(router.Run(":8080"))

}

func (m *Methods) Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := m.db.QueryRow("INSERT INTO users(name, email, age) VALUES($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Age).Scan(&user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Prepare error: %v", err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (m *Methods) GetByID(c *gin.Context) {
	id := c.Param("id")

	var res User

	err := m.db.Get(&res, "select * from users where id = $1 UNION ALL select * from users where id = $2", id, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Prepare error: %v", err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (m *Methods) GetAll(c *gin.Context) {
	var id string
	if err := c.ShouldBindQuery(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var res []User

	err := m.db.Select(&res, "select * from users UNION ALL select * from users_server2")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Prepare error: %v", err)
		return
	}

	c.JSON(http.StatusCreated, res)
}
