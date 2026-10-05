package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	gestor := NuevoGestor()
	entrada := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(`
==== GESTOR DE MEMORIA ====
1. Mostrar memoria
2. Crear proceso
3. Terminar proceso
4. Compactar memoria
5. Mostrar páginas
6. Enviar proceso a SWAP
7. Recuperar proceso de SWAP
8. Salir`)

		switch leerEntero(entrada, "Opción: ") {
		case 1:
			mostrarMemoria(gestor)
		case 2:
			cantidad := leerEntero(entrada, "Páginas: ")
			pid, err := gestor.CrearProceso(cantidad)
			mostrarResultado(fmt.Sprintf("Proceso P%d creado.", pid), err)
		case 3:
			pid := leerEntero(entrada, "PID: ")
			mostrarResultado(fmt.Sprintf("Proceso P%d terminado.", pid), gestor.TerminarProceso(pid))
		case 4:
			gestor.Compactar()
			fmt.Println("Memoria compactada.")
		case 5:
			mostrarProcesos(gestor)
		case 6:
			pid := leerEntero(entrada, "PID: ")
			mostrarResultado(fmt.Sprintf("Proceso P%d enviado a SWAP.", pid), gestor.EnviarASwap(pid))
		case 7:
			pid := leerEntero(entrada, "PID: ")
			mostrarResultado(fmt.Sprintf("Proceso P%d recuperado de SWAP.", pid), gestor.RecuperarDeSwap(pid))
		case 8:
			fmt.Println("Hasta luego.")
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

func mostrarMemoria(g *Gestor) {
	fmt.Printf("\nMemoria: %d/%d marcos libres\n", g.MarcosLibres(), marcosTotales)
	for i, pagina := range g.Memoria {
		if pagina == nil {
			fmt.Printf("[%02d libre] ", i)
		} else {
			fmt.Printf("[%02d P%d-p%d] ", i, pagina.PID, pagina.Numero)
		}
		if (i+1)%4 == 0 {
			fmt.Println()
		}
	}
}

func mostrarProcesos(g *Gestor) {
	if len(g.Procesos) == 0 {
		fmt.Println("No hay procesos.")
		return
	}
	for pid := 1; pid < g.SiguientePID; pid++ {
		proceso, existe := g.Procesos[pid]
		if existe {
			estado := "memoria"
			if proceso.EnSwap {
				estado = "SWAP"
			}
			fmt.Printf("P%d: %d página(s), %s\n", pid, proceso.Paginas, estado)
		}
	}
}

func mostrarResultado(exito string, err error) {
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(exito)
}

func leerEntero(entrada *bufio.Reader, mensaje string) int {
	for {
		fmt.Print(mensaje)
		linea, err := entrada.ReadString('\n')
		if err != nil && len(linea) == 0 {
			return 0
		}
		valor, err := strconv.Atoi(strings.TrimSpace(linea))
		if err == nil {
			return valor
		}
		fmt.Println("Introduce un número válido.")
	}
}
