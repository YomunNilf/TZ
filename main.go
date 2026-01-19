package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

// NumberRequest представляет запрос с числом для сохранения
type NumberRequest struct {
	Number int `json:"number"`
}

// NumbersResponse представляет ответ со списком отсортированных чисел
type NumbersResponse struct {
	Numbers []int `json:"numbers"`
}

// App содержит состояние приложения, включая подключение к базе данных
type App struct {
	DB *sql.DB
}

// main запускает HTTP сервер и инициализирует подключение к базе данных
func main() {
	// Инициализация подключения к базе данных
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	app := &App{DB: db}

	// Регистрация обработчика для эндпоинта /numbers
	http.HandleFunc("/numbers", app.handleNumbers)

	// Получение порта из переменной окружения или использование порта по умолчанию
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// initDB инициализирует подключение к PostgreSQL и создает таблицу, если она не существует
func initDB() (*sql.DB, error) {
	// Получение строки подключения из переменной окружения или использование значения по умолчанию
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost/numbersdb?sslmode=disable"
	}

	// Открытие подключения к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Создание таблицы, если она не существует
	createTable := `
	CREATE TABLE IF NOT EXISTS numbers (
		id SERIAL PRIMARY KEY,
		value INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		return nil, err
	}

	return db, nil
}

// handleNumbers обрабатывает HTTP запросы к эндпоинту /numbers
// Поддерживает POST для добавления числа и GET для получения всех чисел
func (app *App) handleNumbers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Маршрутизация по HTTP методу
	switch r.Method {
	case http.MethodPost:
		app.addNumber(w, r)
	case http.MethodGet:
		app.getNumbers(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// addNumber обрабатывает POST запрос для добавления числа в базу данных
// Поддерживает как JSON формат, так и query параметры
// Возвращает отсортированный список всех чисел
func (app *App) addNumber(w http.ResponseWriter, r *http.Request) {
	var req NumberRequest

	// Попытка сначала распарсить JSON
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		// Попытка распарсить из query параметра
		numberStr := r.URL.Query().Get("number")
		if numberStr == "" {
			http.Error(w, "Number is required", http.StatusBadRequest)
			return
		}
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			http.Error(w, "Invalid number format", http.StatusBadRequest)
			return
		}
		req.Number = number
	}

	// Вставка числа в базу данных
	_, err := app.DB.Exec("INSERT INTO numbers (value) VALUES ($1)", req.Number)
	if err != nil {
		log.Printf("Error inserting number: %v", err)
		http.Error(w, "Failed to save number", http.StatusInternalServerError)
		return
	}

	// Получение всех чисел отсортированными
	numbers, err := app.getAllNumbers()
	if err != nil {
		log.Printf("Error getting numbers: %v", err)
		http.Error(w, "Failed to retrieve numbers", http.StatusInternalServerError)
		return
	}

	// Формирование и отправка ответа
	response := NumbersResponse{Numbers: numbers}
	json.NewEncoder(w).Encode(response)
}

// getNumbers обрабатывает GET запрос для получения всех отсортированных чисел из базы данных
func (app *App) getNumbers(w http.ResponseWriter, r *http.Request) {
	numbers, err := app.getAllNumbers()
	if err != nil {
		log.Printf("Error getting numbers: %v", err)
		http.Error(w, "Failed to retrieve numbers", http.StatusInternalServerError)
		return
	}

	// Формирование и отправка ответа
	response := NumbersResponse{Numbers: numbers}
	json.NewEncoder(w).Encode(response)
}

// getAllNumbers получает все числа из базы данных, отсортированные по возрастанию
func (app *App) getAllNumbers() ([]int, error) {
	// Выполнение SQL запроса для получения всех чисел, отсортированных по возрастанию
	rows, err := app.DB.Query("SELECT value FROM numbers ORDER BY value ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Сканирование результатов запроса в срез
	var numbers []int
	for rows.Next() {
		var num int
		if err := rows.Scan(&num); err != nil {
			return nil, err
		}
		numbers = append(numbers, num)
	}

	return numbers, rows.Err()
}
