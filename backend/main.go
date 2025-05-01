package main

import (
	"os"
	"time"

	"github.com/go-playground/validator/v10" // Añade este import
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rak-nark/sparkpass/config"
	"github.com/rak-nark/sparkpass/routes"
	"gorm.io/gorm"
)

// @title SparkPass API
// @version 1.0
// @description API para gestión de contenido premium y suscripciones
// @host localhost:8080
// @BasePath /api

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// 1. Configuración inicial
	e := echo.New()
	
	// Configurar validador con mensajes en español si lo deseas
	e.Validator = &CustomValidator{validator: validator.New()}
	
	e.HideBanner = false

	// 2. Configurar logger detallado
	configureLogger(e)

	// 3. Conectar a la base de datos
	config.ConnectDB()
	db := config.DB

	// 4. Middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(requestLogger())

	// 5. Configurar rutas
	routes.SetupRoutes(e, db)

	// 6. Mostrar información del servidor al iniciar
	printServerInfo(e, db)

	// 7. Iniciar servidor
	port := getPort()
	e.Logger.Fatal(e.Start(port))
}

// configureLogger establece el nivel de logging y formato
func configureLogger(e *echo.Echo) {
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetHeader("${time_rfc3339} | ${level} | ${short_file}:${line} |")
}

// requestLogger middleware personalizado para registrar solicitudes
func requestLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","method":"${method}","uri":"${uri}",` +
			`"status":${status},"latency":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out},"remote_ip":"${remote_ip}",` +
			`"user_agent":"${user_agent}","error":"${error}"}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:          os.Stdout,
	})
}

// getPort obtiene el puerto de las variables de entorno o usa el predeterminado
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

// printServerInfo muestra información detallada al iniciar el servidor
func printServerInfo(e *echo.Echo, db *gorm.DB) {
	// Obtener información de la base de datos
	var dbVersion string
	db.Raw("SELECT VERSION()").Scan(&dbVersion)

	// Mostrar información del servidor
	e.Logger.Info("🚀 Servidor SparkPass iniciado")
	e.Logger.Infof("📅 Fecha/hora de inicio: %s", time.Now().Format(time.RFC1123))
	e.Logger.Infof("🌐 URL del servidor: http://localhost%s", getPort())
	e.Logger.Infof("🛠️ Modo: %s", "development")
	e.Logger.Infof("📊 Base de datos: %s", dbVersion)
	e.Logger.Info("----------------------------------------------------")
	e.Logger.Info("Rutas disponibles:")
	for _, route := range e.Routes() {
		e.Logger.Infof("%-6s %s", route.Method, route.Path)
	}
	e.Logger.Info("----------------------------------------------------")
}