package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"calculator/internal/constants"
	"calculator/internal/environ_vars"
)

// структура используется для передачи задач агенту
type Task_agent struct {
	Id             string  `json:"id"`
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time float64 `json:"operation_time"`
}

// структура используется для десериализации данных HTTP-запроса агента результата вычисления
type Result_get struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
}

/* Агент
Демон, который получает выражение для вычисления с сервера, вычисляет его и
отправляет на сервер результат выражения.
При старте демон запускает несколько горутин, каждая из которых выступает в роли независимого вычислителя.
Количество горутин регулируется переменной среды COMPUTING_POWER
Агент обязательно общается с оркестратором по http
Агент все время приходит к оркестратору с запросом "дай задачку поработать"
(в ручку GET internal/task для получения задач).
Оркестратор отдаёт задачу.
Агент производит вычисление и в ручку оркестратора (POST internal/task для приёма результатов обработки данных)
отдаёт результат.*/

func Run(c chan os.Signal) {
	_exit := make(chan bool)
	_task := make(chan Task_agent)
	_result := make(chan Result_get)
	g := number_of_goroutines()
	fmt.Printf("Starting %d goroutines\n", g)
	for i := 0; i < g; i++ {
		go evaluator(_task, _result, _exit)
	}

	// экземпляр клиента для общения с сервером-оркестратором по http
	client := &http.Client{}

	go func() {
		for {
			select {
			case _c := <-c:
				fmt.Println("Агент завершает работу")
				for i := 0; i < g; i++ {
					_exit <- true
				}
				go func() {
					c <- _c
				}()
				return
			case <-time.After(1 * time.Second):
				// Агент все время приходит к оркестратору с запросом "дай задачку поработать"
				// (в ручку GET internal/task для получения задач)
				resp, err := client.Get("http://localhost:8080/internal/task")
				if err != nil {
					fmt.Printf("%s, код ответа: %d\n", err.Error(), resp.StatusCode)
				}
				//fmt.Printf("Код ответа: %d\n", resp.StatusCode)
				if resp.StatusCode == 200 {
					task := make(map[string]Task_agent)
					err := json.NewDecoder(resp.Body).Decode(&task)
					if err != nil {
						fmt.Println(err.Error())
					} else {
						// fmt.Println(task)
						go func() {
							_task <- task["task"] // передаём задачу горутине
						}()
					}
				}
			case r := <-_result: // получаем от горутины результат задачи
				// Агент производит вычисление и в ручку оркестратора (POST internal/task для приёма результатов
				// обработки данных) отдаёт результат
				// fmt.Println(r)
				body, err := json.Marshal(r)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					resp, err := client.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(body))
					if err != nil {
						fmt.Println(err.Error())
					}
					resp.Body.Close()
					fmt.Println("Сервер ответил, код:", resp.StatusCode)
					/*req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/internal/task", bytes.NewBuffer(body))
					if err != nil {
						fmt.Println(err.Error())
					} else {
						req.Header.Set("Content-Type", "application/json")
						resp, err := client.Do(req)
						if err != nil {
							fmt.Println(err.Error())
						}
						resp.Body.Close()
					}*/
				}
			}
		}
	}()
}

// горутина - вычислитель (evaluator - оценщик потому, что простые выражения считает)
func evaluator(_task chan Task_agent, _result chan Result_get, _exit chan bool) {
	fmt.Println("Горутина начинает работу")
	time_add, ok := environ_vars.GetValueInt(constants.TimeAdd)
	if !ok {
		time_add = 10000
	}
	time_sub, ok := environ_vars.GetValueInt(constants.TimeSub)
	if !ok {
		time_sub = 10000
	}
	time_mult, ok := environ_vars.GetValueInt(constants.TimeMult)
	if !ok {
		time_mult = 10000
	}
	time_div, ok := environ_vars.GetValueInt(constants.TimeDiv)
	if !ok {
		time_div = 10000
	}
	var add = func(a, b float64) float64 {
		return a + b
	}
	var sub = func(a, b float64) float64 {
		return a - b
	}
	var mult = func(a, b float64) float64 {
		return a * b
	}
	var div = func(a, b float64) float64 {
		return a / b
	}
	for {
		select {
		case <-_exit:
			fmt.Println("Горутина завершает работу")
			return
		case task := <-_task:
			fmt.Println("Получена задача:", task)
			result := Result_get{}
			result.Id = task.Id
			var f func(x, y float64) float64
			var t time.Duration
			switch task.Operation {
			case constants.OPl:
				f = add
				t = time.Duration(time_add) * time.Millisecond
			case constants.OMn:
				f = sub
				t = time.Duration(time_sub) * time.Millisecond
			case constants.OMl:
				f = mult
				t = time.Duration(time_mult) * time.Millisecond
			case constants.ODv:
				f = div
				t = time.Duration(time_div) * time.Millisecond
			}
			result.Result = f(task.Arg1, task.Arg2)
			time.Sleep(t)
			_result <- result
		}
	}
}

// функция получает количество горутин для запуска,
// количество определяется из переменной среды constants.CompPow = "COMPUTING_POWER"
func number_of_goroutines() int {
	goroutines := 1
	g := environ_vars.GetValue(constants.CompPow)
	var err error
	goroutines, err = strconv.Atoi(g)
	if err != nil {
		avail_CPUs := runtime.NumCPU() / 2
		fmt.Printf("Error on conversion environment variable '%s' from string to int! Setting by default in %d\n", constants.CompPow, avail_CPUs)
		goroutines = avail_CPUs
	}

	return goroutines
}
