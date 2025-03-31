package main 

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/cors"
)
type Entrada struct {
	Text string `json:"text"`
}

type StatusResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}


func main() {
	//EndPoint 
	http.HandleFunc("/analizar", getCadenaAnalizar)

	// Configurar CORS con opciones predeterminadas
	//Permisos para enviar y recir informacion
	c := cors.Default()

	// Configurar el manejador HTTP con CORS
	handler := c.Handler(http.DefaultServeMux)

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Servidor escuchando en http://localhost:8080")
	http.ListenAndServe(":8080", handler)

}


func getCadenaAnalizar(w http.ResponseWriter, r *http.Request) {
	var respuesta string
	w.Header().Set("Content-Type", "application/json")
	
	var status StatusResponse
	if r.Method == http.MethodPost {
		var entrada Entrada
		if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
			http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
			status = StatusResponse{Message: "Error al decodificar JSON", Type: "unsucces"}
			json.NewEncoder(w).Encode(status)
			return
		}
		lector := bufio.NewScanner(strings.NewReader(entrada.Text))
		for lector.Scan() {
			if lector.Text() != ""{
				linea := strings.Split(lector.Text(), "#") //comentarios
				if len(linea[0]) != 0 {
					respuesta += "Comando: " + linea[0] + "\n"
					respuesta += "Parametro" 
					respuesta += Analizar(linea[0])  + "\n"
				}	
				//Comentarios			
				if len(linea) > 1 && linea[1] != "" {
					fmt.Println("#"+linea[1] +"\n")
					respuesta += "#"+linea[1] +"\n"
				}
			}
			
		}

		w.WriteHeader(http.StatusOK)

		status = StatusResponse{Message: respuesta, Type: "succes"}
		json.NewEncoder(w).Encode(status)

	} else {
		//http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		status = StatusResponse{Message: "Metodo no permitido", Type: "unsucces"}
		json.NewEncoder(w).Encode(status)
	}
}


var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)



func AnalyzeCommnad(command string, params string) {

	if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params)
	} else if strings.Contains(command, "fdisk") {
		fn_fdisk(params)
	} else if strings.Contains(command, "rmdisk") {
		fn_rmdisk(params)
	} else if strings.Contains(command, "mounted") {
		DiskManagement.PrintMountedPartitions()
	} else if strings.Contains(command, "mount") {
		fn_mount(params)

	} else if strings.Contains(command, "mkfs") {
		fn_mkfs(params)
	} else if strings.Contains(command, "login") {
		fn_login(params)
	} else if strings.Contains(command, "rep") {
		fn_rep(params)
	} else if strings.Contains(command, "mkfile") {
		// Implementar la función mkfile aquí
	} else if strings.Contains(command, "cat") {
		// Implementar la función cat aquí
	} else {
		fmt.Println("Error: Commando invalido o no encontrado")
	}

}