package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

// Database Connection (mejor inicializar una sola vez)
func conexionBD() {
	godotenv.Load()

	driver := os.Getenv("Driver")
	usuario := os.Getenv("Usuario")
	contrasena := os.Getenv("Contrasena")
	nombre := os.Getenv("Nombre")

	var err error
	db, err = sql.Open(driver, usuario+":"+contrasena+"@tcp(0.0.0.0:3306)/"+nombre)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

// Structs
type Servicio struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Correo      string `json:"correo"`
	Telefono    string `json:"telefono"`
	Equipo      string `json:"equipo"`
	Diagnostico string `json:"diagnostico"`
	Resultados  string `json:"resultados"`
	Decision    string `json:"decision"`
	Taller      string `json:"taller"`
	Servicio    string `json:"servicio"`
	Entrega     string `json:"entrega"`
	Fecha       string `json:"fecha"`
}

type Servicios struct {
	Servicios []Servicio `json:"servicios"`
}

func main() {
	godotenv.Load()

	conexionBD()

	router := gin.Default()

	// CORS
	router.Use(cors.Default())

	port := os.Getenv("Port")

	if port == "" {
		port = "Port"
	}

	// Routes
	router.GET("/servicio", Get)
	router.GET("/servicio/:id", GetServicio)
	router.POST("/servicio", Post)
	router.PUT("/servicio/:id", Put)
	router.DELETE("/servicio/:id", Delete)

	log.Println("Server running on http://0.0.0.0:" + port + "/servicio")
	router.Run("0.0.0.0" + port)
}

// GET ALL
func Get(c *gin.Context) {
	rows, err := db.Query("SELECT id, nombre, correo, telefono, equipo, diagnostico, resultados, decision, taller, servicio, entrega, fecha FROM taller")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	result := Servicios{}

	for rows.Next() {
		var s Servicio
		if err := rows.Scan(&s.ID, &s.Nombre, &s.Correo, &s.Telefono, &s.Equipo, &s.Diagnostico, &s.Resultados, &s.Decision, &s.Taller, &s.Servicio, &s.Entrega, &s.Fecha); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result.Servicios = append(result.Servicios, s)
	}

	c.JSON(http.StatusOK, result)
}

// GET ONE
func GetServicio(c *gin.Context) {
	id := c.Param("id")

	rows, err := db.Query("SELECT id, nombre, correo, telefono, equipo, diagnostico, resultados, decision, taller, servicio, entrega, fecha FROM taller WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	result := Servicios{}

	for rows.Next() {
		var s Servicio
		if err := rows.Scan(&s.ID, &s.Nombre, &s.Correo, &s.Telefono, &s.Equipo, &s.Diagnostico, &s.Resultados, &s.Decision, &s.Taller, &s.Servicio, &s.Entrega, &s.Fecha); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result.Servicios = append(result.Servicios, s)
	}

	c.JSON(http.StatusOK, result)
}

// POST
func Post(c *gin.Context) {
	var s Servicio

	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(
		"INSERT INTO taller (nombre, correo, telefono, equipo, diagnostico, resultados, decision, taller, servicio, entrega, fecha) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		s.Nombre, s.Correo, s.Telefono, s.Equipo, s.Diagnostico, s.Resultados, s.Decision, s.Taller, s.Servicio, s.Entrega, s.Fecha,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}

// PUT
func Put(c *gin.Context) {
	id := c.Param("id")

	var s Servicio
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error al parsear"})
		return
	}

	result, err := db.Exec(
		"UPDATE taller SET nombre=?,correo=?,telefono=?,equipo=?,diagnostico=?,resultados=?,decision=?,taller=?,servicio=?,entrega=?,fecha=? WHERE id=?",
		s.Nombre, s.Correo, s.Telefono, s.Equipo, s.Diagnostico, s.Resultados, s.Decision, s.Taller, s.Servicio, s.Entrega, s.Fecha, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al actualizar"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Registro no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Actualizado correctamente"})
}

// DELETE
func Delete(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM taller WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al eliminar"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Registro no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Eliminado correctamente"})
}
