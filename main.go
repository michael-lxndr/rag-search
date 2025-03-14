package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RequestData struct {
	Usuario string `json:"usuario"`
	Mensaje string `json:"mensaje"`
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Leer el body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Convertir el JSON a struct
	var data RequestData
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error en el formato JSON", http.StatusBadRequest)
		return
	}

	// Guardar el JSON en un archivo
	file, err := os.Create("datos.json")
	if err != nil {
		http.Error(w, "No se pudo guardar el archivo", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		http.Error(w, "Error al escribir en el archivo", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Archivo JSON recibido y guardado\n")
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
