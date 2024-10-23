package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Обработчик для получения всех задач
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("Ошибка при кодировании задач: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Обработчик для отправки задачи на сервер
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	if err := json.Unmarshal(buf.Bytes(), &newTask); err != nil {
		log.Printf("Ошибка при декодировании задачи: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tasks[newTask.ID] = newTask
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, exists := tasks[id]
	if !exists {
		log.Printf("Задача с ID %s не найдена", id)
		http.Error(w, fmt.Sprintf("Задача с ID %s не найдена", id), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data, err := json.Marshal(task)
	if err != nil {
		log.Printf("Ошибка при кодировании задачи с ID %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Обработчик удаления задачи по ID
func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, exists := tasks[id]
	if !exists {
		log.Printf("Задача для удаления с ID %s не найдена", id)
		http.Error(w, fmt.Sprintf("Задача с ID %s не найдена", id), http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// Регистрация обработчиков
	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskByID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
