# Gestor de memoria en Go

## Guion para una exposición de 5 minutos

## 1. Objetivo del programa

Este proyecto es un gestor de memoria sencillo escrito en Go. Simula una
memoria principal dividida en 16 marcos. Los procesos se dividen en páginas y
cada página ocupa un marco.

El programa permite:

1. Mostrar el estado de la memoria.
2. Crear un proceso.
3. Terminar un proceso.
4. Compactar la memoria.
5. Mostrar los procesos y sus páginas.
6. Enviar un proceso a SWAP.
7. Recuperar un proceso de SWAP.
8. Salir.

La simulación busca ser fácil de entender. No representa bytes reales ni usa
un archivo de SWAP: mantiene el estado en memoria mientras el programa está
ejecutándose.

## 2. Organización de los archivos

### `main.go`

Contiene la función `main`, el menú y las funciones que muestran información o
leen datos del usuario. Este archivo se encarga de la interfaz del programa.

### `memoria.go`

Contiene las estructuras de datos y la lógica del gestor: creación y
terminación de procesos, asignación de páginas, compactación y SWAP.

### `go.mod`

Define el módulo `gestor-memoria` y especifica que se usa Go 1.22.

## 3. Modelo de datos

En `memoria.go` se define la constante:

```go
const marcosTotales = 16
```

Esto significa que la memoria principal tiene 16 posiciones.

La estructura `Pagina` identifica qué proceso ocupa una posición y cuál es el
número de esa página:

```go
type Pagina struct {
    PID    int
    Numero int
}
```

La estructura `Proceso` guarda el PID, la cantidad de páginas y si está en
SWAP:

```go
type Proceso struct {
    PID     int
    Paginas int
    EnSwap  bool
}
```

El gestor reúne todo el estado:

```go
type Gestor struct {
    Memoria      [marcosTotales]*Pagina
    Procesos     map[int]Proceso
    SiguientePID int
}
```

- `Memoria` es un arreglo de 16 posiciones. Un puntero `nil` significa que el
  marco está libre.
- `Procesos` es un mapa que relaciona cada PID con su proceso.
- `SiguientePID` genera automáticamente los identificadores `P1`, `P2`, etc.

La función `NuevoGestor` crea el mapa vacío y comienza los PID en 1.

## 4. Operaciones principales

### Crear proceso

`CrearProceso` recibe la cantidad de páginas solicitada. Primero comprueba que
sea mayor que cero y que existan suficientes marcos libres. Después genera un
PID, guarda el proceso y llama a `cargarPaginas`.

`cargarPaginas` recorre la memoria desde el primer marco y coloca cada página
en la primera posición libre. Por ejemplo, un proceso de 3 páginas puede
aparecer como `P1-p1`, `P1-p2` y `P1-p3`.

### Mostrar memoria

La función `mostrarMemoria` de `main.go` recorre el arreglo `Memoria` y muestra
cada marco como libre o con el proceso y la página que contiene. También usa
`MarcosLibres` para mostrar cuántos marcos quedan disponibles.

### Mostrar páginas

`mostrarProcesos` recorre los procesos existentes y muestra su cantidad de
páginas y su estado: `memoria` o `SWAP`.

### Terminar proceso

`TerminarProceso` comprueba que el PID exista, llama a `liberarPaginas` y
elimina el proceso del mapa. `liberarPaginas` recorre toda la memoria y deja
libres las posiciones pertenecientes a ese PID.

### Compactar memoria

`Compactar` crea un arreglo nuevo vacío y copia en él todas las páginas que no
son `nil`, empezando desde la posición cero. Así elimina los huecos entre
páginas y deja todos los marcos libres juntos al final.

La compactación no cambia el PID ni el número de página; solo cambia la
posición física que ocupa cada página.

### Enviar a SWAP

`EnviarASwap` busca el proceso, verifica que no esté ya en SWAP, libera sus
páginas y cambia `EnSwap` a `true`.

En esta versión, SWAP es una simulación lógica: el proceso sigue registrado en
`Procesos`, pero ya no ocupa memoria principal.

### Recuperar de SWAP

`RecuperarDeSwap` comprueba que el proceso exista y tenga `EnSwap == true`.
También verifica que haya suficientes marcos libres. Si todo es correcto,
vuelve a cargar sus páginas y cambia `EnSwap` a `false`.

## 5. Flujo del menú

La función `main` crea el gestor y un `bufio.Reader` para leer la entrada del
teclado. Después ejecuta un bucle infinito:

1. Muestra el menú.
2. Lee la opción con `leerEntero`.
3. Ejecuta el método correspondiente mediante `switch`.
4. Repite hasta seleccionar la opción 8.

`leerEntero` convierte el texto introducido a entero con `strconv.Atoi`. Si el
usuario escribe algo que no es un número, muestra un error y vuelve a pedir el
valor.

`mostrarResultado` evita repetir el mismo patrón de control de errores: si hay
un error, lo imprime; si no, muestra el mensaje de éxito.

## 6. Ejemplo para demostrarlo

Durante la exposición se puede hacer esta secuencia:

1. Elegir `2` y crear un proceso de 3 páginas.
2. Elegir `2` y crear otro proceso de 2 páginas.
3. Elegir `1` para observar los marcos ocupados.
4. Elegir `6`, enviar `P1` a SWAP y comprobar que sus marcos quedan libres.
5. Elegir `5` para ver que `P1` aparece en SWAP y `P2` en memoria.
6. Elegir `7` y recuperar `P1`.
7. Elegir `4` para compactar la memoria.
8. Elegir `3` y terminar `P2`.
9. Elegir `8` para salir.

## 7. Cierre de la exposición

La idea central es separar la interfaz de la lógica. `main.go` pregunta y
muestra resultados, mientras `memoria.go` administra los datos. La memoria se
modela como un arreglo de marcos, las páginas identifican qué ocupa cada
marco, y el mapa de procesos permite consultar el estado de cada proceso.

Para ejecutar el proyecto:

```bash
go run .
```

Para comprobar que compila:

```bash
go test ./...
go vet ./...
```