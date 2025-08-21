entrada = input("Ingrese números separados por comas: ")
numeros = entrada.replace(" ", "").split(",")
sin_duplicados = list(set(numeros))
ordenados = sorted(sin_duplicados, reverse=True)

print("Lista original:", numeros)
print("Sin duplicados:", sin_duplicados)
print("Ordenados de mayor a menor:", ordenados)