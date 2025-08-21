notas_estudiantes = {
    "juli": 4.5,
    "camilo": 3.8,
    "noe": 4.0
}

print("Diccionario de notas:", notas_estudiantes)

print("\nNotas por  estudiante:")
for nombre, nota in notas_estudiantes.items():
    print(f"{nombre}: {nota}")