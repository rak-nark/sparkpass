package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB es una variable global que contendrá la conexión a la base de datos.
var DB *gorm.DB

// ConnectDB establece la conexión a la base de datos y maneja errores.
func ConnectDB() {
	// Obtener las credenciales de la base de datos desde variables de entorno
	dsn := getDSN()
	
	// Intentar abrir la conexión con la base de datos
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// Loguear el error con detalles y terminar la ejecución
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Si la conexión es exitosa, imprimir la ruta completa de la base de datos
	fmt.Printf("🚀 Conexión exitosa a la base de datos!\n")
	fmt.Printf("Ruta completa de la base de datos: %s\n", dsn)

	// Asignar la conexión abierta a la variable global DB
	DB = db
}

// getDSN construye la ruta completa para la conexión a la base de datos.
// Se pueden usar variables de entorno para mayor seguridad y flexibilidad.
func getDSN() string {
	// Obtener las credenciales de las variables de entorno (si están definidas)
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_NAME", "sparkpass")
	
	// Construir la cadena DSN (Data Source Name) para la conexión
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", user, password, host, port, dbname)
	return dsn
}

// getEnv obtiene el valor de una variable de entorno, o retorna un valor por defecto si no está definida.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
