package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el tipo de archivo de los parámetros de consulta
	fileType := r.Header.Get("X-File-Type")

	var filename string
	switch fileType {
	case "keywords":
		filename = "data_keywords.json"
	case "rag":
		filename = "data_rag.json"
	default:
		filename = "data.json" // Nombre predeterminado
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Guardar el JSON en el archivo correspondiente
	file, err := os.Create(filename)
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

	fmt.Fprintf(w, "Archivo %s guardado correctamente\n", filename)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
