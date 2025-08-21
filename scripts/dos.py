palabra = input("Ingrese una palabra: ").lower()

conteo = {}

for letra in palabra:
    if letra in conteo:
        conteo[letra] += 1   # Si ya existe la letra, incrementamos  en 1
    else:
        conteo[letra] = 1    # Si no existe, le añadimos  con valor 1


print("Conteo de letras:", conteo)

