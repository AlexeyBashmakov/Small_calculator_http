package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"calculator/internal/service"

	"github.com/go-chi/chi/v5"
)

func New(ctx context.Context, calcService service.Calc) (http.Handler, error) {
	serveMux := chi.NewRouter()

	// serveMux.Get("/", service.Index)

	// Добавление вычисления арифметического выражения
	serveMux.Post("/api/v1/calculate", calcService.Calculate)
	// Получение списка выражений
	serveMux.Get("/api/v1/expressions", calcService.Expressions)
	// Получение выражения по его идентификатору
	serveMux.Get("/api/v1/expressions/:{id}", calcService.Expression_id)
	// Получение задачи агентом для выполнения
	serveMux.Get("/internal/task", calcService.Task_get)
	// Прием результата обработки данных от агента
	serveMux.Post("/internal/task", calcService.Task_result)

	// os.Chdir("..\\..\\front\\")
	// Добавьте следующие две строки
	// fs := http.FileServer(http.Dir("/"))
	// serveMux.Handle("/", http.StripPrefix("/", fs))
	// Настройка раздачи статических файлов
	staticPath, _ := filepath.Abs("../../front/")
	fs := http.FileServer(http.Dir(staticPath))
	serveMux.Handle("/*", fs)

	return serveMux, nil
}
