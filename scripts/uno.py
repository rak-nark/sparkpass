diccionario ={
    "car" :"carro",
    "mot" :"moto",
    "bic" :"bicicleta",
    "cam" :"camion",
    "avi" :"avion",
}
palabra = input("Ingrese una palabra en inglés: ").lower()

if palabra in diccionario:
    print("La traducción de", palabra, "es:", diccionario[palabra])
else:
    print("La palabra no está en el diccionario.")