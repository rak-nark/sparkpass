entrada = input("ingrese letras seleccionadas por comas: ")
letras = entrada.replace(" ", "").split(",")
salto = int(input("Ingrese el número de salto: "))

nueva_lista = letras[::salto]
print("Lista original:", letras)
print("Salto:", salto)
print("Nueva lista con salto:", nueva_lista)
