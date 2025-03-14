package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const uploadDir = "tmp" // Directorio de carga

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el tipo de archivo del encabezado
	fileType := r.Header.Get("X-File-Type")

	// Determinar el nombre del archivo basado en el tipo
	var filename string
	switch fileType {
	case "keywords":
		filename = "data_keywords.json"
	case "rag":
		filename = "data_rag.json"
	default:
		filename = "data.json" // Nombre predeterminado
	}

	// Crear la ruta completa del archivo
	filePath := filepath.Join(uploadDir, filename)

	// Asegurar que el directorio de carga exista
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Error al crear el directorio: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error interno del servidor")
		return
	}

	// Leer el cuerpo de la solicitud
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error al leer el cuerpo de la solicitud: %v", err)
		respondWithError(w, http.StatusBadRequest, "Error al leer el cuerpo de la solicitud")
		return
	}
	defer r.Body.Close()

	// Guardar el archivo
	if err := os.WriteFile(filePath, body, 0644); err != nil {
		log.Printf("Error al guardar el archivo: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error al guardar el archivo")
		return
	}

	// Responder con éxito
	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: fmt.Sprintf("Archivo %s guardado correctamente", filename),
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, Response{
		Status:  "error",
		Message: message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
