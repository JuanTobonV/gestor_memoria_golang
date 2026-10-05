package main

import "fmt"

const marcosTotales = 16

type Pagina struct {
	PID    int
	Numero int
}

type Proceso struct {
	PID     int
	Paginas int
	EnSwap  bool
}

type Gestor struct {
	Memoria      [marcosTotales]*Pagina
	Procesos     map[int]Proceso
	SiguientePID int
}

func NuevoGestor() *Gestor {
	return &Gestor{Procesos: make(map[int]Proceso), SiguientePID: 1}
}

func (g *Gestor) CrearProceso(cantidad int) (int, error) {
	if cantidad <= 0 {
		return 0, fmt.Errorf("la cantidad de páginas debe ser mayor que cero")
	}
	if cantidad > g.MarcosLibres() {
		return 0, fmt.Errorf("no hay suficientes marcos libres")
	}

	pid := g.SiguientePID
	g.SiguientePID++
	g.Procesos[pid] = Proceso{PID: pid, Paginas: cantidad}
	g.cargarPaginas(pid, cantidad)
	return pid, nil
}

func (g *Gestor) TerminarProceso(pid int) error {
	if _, existe := g.Procesos[pid]; !existe {
		return fmt.Errorf("el proceso P%d no existe", pid)
	}
	g.liberarPaginas(pid)
	delete(g.Procesos, pid)
	return nil
}

func (g *Gestor) Compactar() {
	var compactada [marcosTotales]*Pagina
	indice := 0
	for _, pagina := range g.Memoria {
		if pagina != nil {
			compactada[indice] = pagina
			indice++
		}
	}
	g.Memoria = compactada
}

func (g *Gestor) EnviarASwap(pid int) error {
	proceso, existe := g.Procesos[pid]
	if !existe {
		return fmt.Errorf("el proceso P%d no existe", pid)
	}
	if proceso.EnSwap {
		return fmt.Errorf("el proceso P%d ya está en SWAP", pid)
	}
	g.liberarPaginas(pid)
	proceso.EnSwap = true
	g.Procesos[pid] = proceso
	return nil
}

func (g *Gestor) RecuperarDeSwap(pid int) error {
	proceso, existe := g.Procesos[pid]
	if !existe || !proceso.EnSwap {
		return fmt.Errorf("el proceso P%d no está en SWAP", pid)
	}
	if proceso.Paginas > g.MarcosLibres() {
		return fmt.Errorf("no hay suficientes marcos libres")
	}
	g.cargarPaginas(pid, proceso.Paginas)
	proceso.EnSwap = false
	g.Procesos[pid] = proceso
	return nil
}

func (g *Gestor) MarcosLibres() int {
	libres := 0
	for _, pagina := range g.Memoria {
		if pagina == nil {
			libres++
		}
	}
	return libres
}

func (g *Gestor) cargarPaginas(pid, cantidad int) {
	numero := 1
	for i := range g.Memoria {
		if g.Memoria[i] == nil && numero <= cantidad {
			g.Memoria[i] = &Pagina{PID: pid, Numero: numero}
			numero++
		}
	}
}

func (g *Gestor) liberarPaginas(pid int) {
	for i, pagina := range g.Memoria {
		if pagina != nil && pagina.PID == pid {
			g.Memoria[i] = nil
		}
	}
}
