persona = {
    "nombre": "Julieth",
    "edad": 30,
    "ciudad": "Bogotá",
    "profesion": "Estudiante"
}
print("Información personal:")
print(f"Nombre: {persona['nombre']}")
print(f"Edad: {persona['edad']}")
print(f"Ciudad: {persona['ciudad']}")
print(f"Profesión: {persona['profesion']}")

persona["email"] = "julieth@gmail.com"

print("\nDiccionario actualizado con email:")
for clave, valor in persona.items():
    print(f"{clave.capitalize()}: {valor}")
