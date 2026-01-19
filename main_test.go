package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq" // Драйвер PostgreSQL для тестов
)

// setupTestDB создает тестовую базу данных и возвращает функцию очистки
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// Подключение к тестовой базе данных
	connStr := "postgres://postgres:postgres@localhost/numbersdb_test?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skipf("Skipping test: could not connect to test database: %v", err)
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		t.Skipf("Skipping test: could not ping test database: %v", err)
	}

	// Создание таблицы для тестов
	createTable := `
	CREATE TABLE IF NOT EXISTS numbers (
		id SERIAL PRIMARY KEY,
		value INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Функция очистки данных после тестов
	cleanup := func() {
		db.Exec("DELETE FROM numbers")
		db.Close()
	}

	return db, cleanup
}

// TestAddNumber тестирует добавление чисел и проверяет, что они возвращаются отсортированными
func TestAddNumber(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	app := &App{DB: db}

	// Тестовые случаи: добавление чисел 3, 2, 1 и проверка сортировки
	tests := []struct {
		name           string
		number         int
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Add first number 3",
			number:         3,
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "Add second number 2",
			number:         2,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Add third number 1",
			number:         1,
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
	}

	var allNumbers []int

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := NumberRequest{Number: tt.number}
			jsonBody, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/numbers", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			app.addNumber(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response NumbersResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if len(response.Numbers) != tt.expectedCount {
				t.Errorf("Expected %d numbers, got %d", tt.expectedCount, len(response.Numbers))
			}

			// Проверка того, что числа отсортированы
			for i := 1; i < len(response.Numbers); i++ {
				if response.Numbers[i] < response.Numbers[i-1] {
					t.Errorf("Numbers are not sorted: %v", response.Numbers)
				}
			}

			allNumbers = response.Numbers
		})
	}

	// Финальная проверка: должен быть результат [1, 2, 3]
	expectedFinal := []int{1, 2, 3}
	if len(allNumbers) != len(expectedFinal) {
		t.Errorf("Expected final numbers %v, got %v", expectedFinal, allNumbers)
	}
	for i, num := range expectedFinal {
		if allNumbers[i] != num {
			t.Errorf("Expected final numbers %v, got %v", expectedFinal, allNumbers)
			break
		}
	}
}

// TestGetNumbers тестирует получение всех чисел и проверяет сортировку результата
func TestGetNumbers(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	app := &App{DB: db}

	// Сначала добавляем несколько чисел в произвольном порядке
	numbers := []int{5, 1, 3, 2, 4}
	for _, num := range numbers {
		app.DB.Exec("INSERT INTO numbers (value) VALUES ($1)", num)
	}

	req := httptest.NewRequest(http.MethodGet, "/numbers", nil)
	w := httptest.NewRecorder()

	app.getNumbers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response NumbersResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expected := []int{1, 2, 3, 4, 5}
	if len(response.Numbers) != len(expected) {
		t.Errorf("Expected %d numbers, got %d", len(expected), len(response.Numbers))
	}

	for i, num := range expected {
		if response.Numbers[i] != num {
			t.Errorf("Expected sorted numbers %v, got %v", expected, response.Numbers)
			break
		}
	}
}

// TestAddNumberInvalidInput тестирует обработку невалидного JSON при добавлении числа
func TestAddNumberInvalidInput(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	app := &App{DB: db}

	// Отправка невалидного JSON
	req := httptest.NewRequest(http.MethodPost, "/numbers", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.addNumber(w, req)

	// Ожидается статус Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
