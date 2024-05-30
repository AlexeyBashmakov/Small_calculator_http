package server

import (
	"context"
	"fmt"
	"net/http"

	"calculator/http/server/handler"
	"calculator/internal/service"
)

// маршрутизация
// func new(ctx context.Context,
//
//	logger *zap.Logger,
//	lifeService service.LifeService,
func new(ctx context.Context, calcService service.Calc) (http.Handler, error) {
	muxHandler, err := handler.New(ctx, calcService) //, lifeService)
	if err != nil {
		return nil, fmt.Errorf("handler initialization error: %w", err)
	}
	// // middleware для обработчиков
	// muxHandler = handler.Decorate(muxHandler, loggingMiddleware(logger))

	return muxHandler, nil
}

func Run(ctx context.Context) (func(context.Context) error, error) {
	// сервис с игрой
	// lifeService, err := service.New(height, width)
	// if err != nil {
	// 	return nil, err
	// }
	calcService := service.NewCalc()

	muxHandler, err := new(ctx, *calcService) //, logger, *lifeService)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{Addr: "localhost:8080", Handler: muxHandler}

	go func() {
		// Запускаем сервер
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
			// logger.Error("ListenAndServe",
			// 	zap.String("err", err.Error()))
		}
	}()
	// вернем функцию для завершения работы сервера
	return srv.Shutdown, nil
}
