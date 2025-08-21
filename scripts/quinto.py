entrada_nombres = input("Ingrese nombres separados por comas: ")
entrada_apellidos = input("Ingrese apellidos separados por comas: ")

nombres = entrada_nombres.replace(" ", "").split(",")
apellidos = entrada_apellidos.replace(" ", "").split(",")

nombres_completos = [nombre + " " + apellido for nombre, apellido in zip(nombres, apellidos)]
print("Lista de nombres completos:", nombres_completos)