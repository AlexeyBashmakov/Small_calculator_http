package application

import (
	"context"
	"os"
	"os/signal"

	"calculator/http/agent"
	"calculator/http/server"
)

func Run(ctx context.Context) int {
	// Создание логгера с настройками для production
	// logger := setupLogger()

	shutDownFunc, err := server.Run(ctx) //, logger, a.Cfg.Height, a.Cfg.Width)
	if err != nil {
		// logger.Error(err.Error())

		return 1 // вернем код для регистрация ошибки системой
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	_, cancel := context.WithCancel(context.Background())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// НАВЕРНО ГДЕ-ТО ЗДЕСЬ НАДО ЗАПУСТИТЬ АГЕНТА
	agent.Run(c)

	_c := <-c
	go func() {
		c <- _c
	}()
	cancel()
	//  завершим работу сервера
	shutDownFunc(ctx)

	return 0
}

/*
func CheckEnvironmentVariables() bool {
	env_vars := []string{
		constants.TimeAdd,
		constants.TimeSub,
		constants.TimeMult,
		constants.TimeDiv,
		constants.CompPow,
	}
	exist := true
	for _, k := range env_vars {
		_, exist = os.LookupEnv(k)
	}
	return exist
}

func SetEnvironmentVariables() error {
	env_vars := map[string]string{
		constants.TimeAdd:  "10",
		constants.TimeSub:  "10",
		constants.TimeMult: "10",
		constants.TimeDiv:  "10",
		constants.CompPow:  fmt.Sprintf("%d", runtime.NumCPU()/2),
	}
	var err error
	for k, v := range env_vars {
		err = os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}*/
