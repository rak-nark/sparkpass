package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/config"
	"github.com/rak-nark/sparkpass/models"
	"github.com/rak-nark/sparkpass/utils"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context) error {
	// Crear un objeto user para recibir los datos
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos inválidos"})
	}

	// Verificar si el correo ya existe
	var existingUser models.User
	if result := config.DB.Where("email = ?", user.Email).First(&existingUser); result.Error == nil {
		// El correo ya está registrado
		return c.JSON(http.StatusConflict, map[string]string{"error": "El correo ya está registrado"})
	}

	// Generar el hash de la contraseña proporcionada
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al encriptar la contraseña"})
	}
	user.PasswordHash = string(hashedPass)

	// Guardar al usuario en la base de datos
	if result := config.DB.Create(&user); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al crear el usuario"})
	}

	return c.JSON(http.StatusCreated, user)
}

func Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos inválidos"})
	}

	// Imprimir el email recibido para verificar que es correcto
	fmt.Println("Email recibido:", req.Email)

	// Buscar usuario

	var user models.User
	if result := config.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		fmt.Printf("Error buscando usuario: %v\n", result.Error) // Debug detallado
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Credenciales incorrectas"})
	}

	fmt.Printf("Usuario encontrado: %+v\n", user) // Debug: ver toda la estructura

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		fmt.Printf("Error comparando passwords: %v\nHash almacenado: %s\nPassword recibido: %s\n",
			err, user.PasswordHash, req.Password) // Debug detallado
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Credenciales incorrectas"})
	}

	// Generar JWT
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		fmt.Println("Error generando token:", err) // Depurar generación de token
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generando token"})
	}

	// Devuelve el token y el email del usuario
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
		"email": user.Email,
	})
}
