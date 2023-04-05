package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Usuario struct {
	ID           int       `json:"id"`
	NomeUsuario  string    `json:"nome_usuario"`
	Senha        string    `json:"senha"`
	Email        string    `json:"email"`
	DataInclusao time.Time `json:"data_inclusao"`
	Desativado   bool      `json:"desativado"`
}

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/database")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/usuarios", getUsuarios(db))
	r.GET("/usuarios/:id", getUsuario(db))
	r.POST("/usuarios", createUsuario(db))
	r.PUT("/usuarios/:id", updateUsuario(db))
	r.DELETE("/usuarios/:id", deleteUsuario(db))

	r.Run()
}

func getUsuarios(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM usuarios")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		usuarios := []Usuario{}

		for rows.Next() {
			u := Usuario{}
			err := rows.Scan(&u.ID, &u.NomeUsuario, &u.Senha, &u.Email, &u.DataInclusao, &u.Desativado)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			usuarios = append(usuarios, u)
		}

		c.JSON(http.StatusOK, usuarios)
	}
}

func getUsuario(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		row := db.QueryRow("SELECT * FROM usuarios WHERE id = ?", id)

		u := Usuario{}
		err := row.Scan(&u.ID, &u.NomeUsuario, &u.Senha, &u.Email, &u.DataInclusao, &u.Desativado)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, u)
	}
}

func createUsuario(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u Usuario
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := db.Exec("INSERT INTO usuarios (nome_usuario, senha, email, data_inclusao, desativado) VALUES (?, ?, ?, ?, ?)", u.NomeUsuario, u.Senha, u.Email, time.Now(), u.Desativado)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		u.ID = int(lastInsertId)

		c.JSON(http.StatusOK, u)
	}
}

func updateUsuario(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var u Usuario
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("UPDATE usuarios SET nome_usuario = ?, senha = ?, email = ?, desativado = ? WHERE id = ?", u.NomeUsuario, u.Senha, u.Email, u.Desativado, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		u.ID = idInt

		c.JSON(http.StatusOK, u)
	}
}

func deleteUsuario(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM usuarios WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}
