package transactions

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database has been connected.")
}

func CreateTable() {
	var err error
	createTb := `
		CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT,
			status TEXT
		);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("Can't create table or table existed.", err)
	}

	log.Println("Table has been created.")
}

func CreateCustomerHandler(c *gin.Context) {
	cust := Customer{}
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row := db.QueryRow("INSERT INTO customers (name, email, status) values ($1,$2,$3) RETURNING id", cust.Name, cust.Email, cust.Status)
	err := row.Scan(&cust.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cust)
}

func GetCustomerByIDHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id,name,email,status FROM customers WHERE id=$1")
	if err != nil {
		log.Fatal("Can't prepare SELECT statement", err)
	}
	row := stmt.QueryRow(id)
	cust := &Customer{}
	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cust)
}

func GetCustomerHandler(c *gin.Context) {
	stmt, err := db.Prepare("SELECT id,name,email,status FROM customers")
	if err != nil {
		log.Fatal("Can't prepare SELECT statement", err)
	}
	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	custs := []Customer{}
	for rows.Next() {
		cust := Customer{}
		err := rows.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		custs = append(custs, cust)
	}
	c.JSON(http.StatusOK, custs)
}

func UpdateCustomerHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id,name,email,status FROM customers WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	row := stmt.QueryRow(id)
	cust := &Customer{}
	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stmt, err = db.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1")
	if err != nil {
		log.Fatal("Can't prepare UPDATE statement", err)
	}
	if _, err := stmt.Exec(id, cust.Name, cust.Email, cust.Status); err != nil {
		log.Fatal("Can't execute UPDATE statement", err)
	}
	c.JSON(http.StatusOK, cust)
}

func DeleteCustomerByIDHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM customers WHERE id=$1")
	if err != nil {
		log.Fatal("Can't prepare DELETE statement", err)
	}
	if _, err := stmt.Exec(id); err != nil {
		log.Fatal("Can't execute DELETE statement", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}
