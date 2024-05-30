package service

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"calculator/internal/constants"
	"calculator/internal/environ_vars"
	"calculator/pkg/my_queue"
	"calculator/pkg/rpn"

	"github.com/go-chi/chi/v5"
)

const wait = "в очереди"
const calculate = "вычисляется"
const finished = "завершено"

// структура используется для десериализации данных HTTP-запроса клиента добавления задачи
type Task_add struct {
	Id   int    `json:"id"`
	Expr string `json:"expression"`
}

// структура используется для десериализации данных HTTP-запроса агента результата вычисления
type Result_get struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
}

// структура используется для передачи задач агенту
type Task_agent struct {
	Id             string  `json:"id"`
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time float64 `json:"operation_time"`
}

type Description struct {
	Expression string
	RPN        []string
	Status     string
	Result     float64
}

type Calc struct {
	Pool  map[int]Description
	Queue my_queue.ConcurrentQueue
}

func NewCalc() *Calc {
	return &Calc{Pool: make(map[int]Description), Queue: my_queue.ConcurrentQueue{}}
}

// полученную ОПН разбиваем на задачи
func (c *Calc) rpn_to_tasks(id int) {
	for i := range c.Pool[id].RPN {
		s0 := c.Pool[id].RPN[i]
		if s0 == "" {
			continue
		}
		if i+1 == len(c.Pool[id].RPN) {
			break
		}
		s1 := c.Pool[id].RPN[i+1]
		a := 0
		// fmt.Printf("rpn_to_tasks, i = %d, s0 = '%s', s1 = '%s'\n", i, s0, s1)
		for s1 == "" {
			a++
			if i+a+1 == len(c.Pool[id].RPN) {
				return
			}
			s1 = c.Pool[id].RPN[i+a+1]
		}
		// fmt.Printf("rpn_to_tasks, i = %d, a = %d\n", i, a)
		if i+a+2 == len(c.Pool[id].RPN) {
			break
		}
		s2 := c.Pool[id].RPN[i+a+2]
		b := 0
		for s2 == "" {
			b++
			if i+a+b+2 == len(c.Pool[id].RPN) {
				return
			}
			s2 = c.Pool[id].RPN[i+a+b+2]
		}
		/* после цикла было
		s0 := c.Pool[id].RPN[i]
		s1 := c.Pool[id].RPN[i+1]
		s2 := c.Pool[id].RPN[i+2]
		*/

		if (s2 == constants.OPl) || (s2 == constants.OMn) || (s2 == constants.OMl) || (s2 == constants.ODv) {
			// если третий символ - операция, то проверяю два предыдущих
			if n0, e0 := strconv.ParseFloat(s0, 64); e0 == nil {
				if n1, e1 := strconv.ParseFloat(s1, 64); e1 == nil {
					// имеем два числа
					task_id := fmt.Sprintf("%d.%d", id, i)
					operation_time := ""
					switch s2 {
					case constants.OPl:
						operation_time = environ_vars.GetValue(constants.TimeAdd)
					case constants.OMn:
						operation_time = environ_vars.GetValue(constants.TimeSub)
					case constants.OMl:
						operation_time = environ_vars.GetValue(constants.TimeMult)
					case constants.ODv:
						operation_time = environ_vars.GetValue(constants.TimeDiv)
					}
					var op_t float64
					if op_t, e0 = strconv.ParseFloat(operation_time, 64); e0 != nil {
						op_t = 1
					}
					c.Queue.Enqueue(Task_agent{Id: task_id, Arg1: n0, Arg2: n1, Operation: s2, Operation_time: op_t})
					c.Pool[id].RPN[i] = ""
					c.Pool[id].RPN[i+a+1] = ""
					c.Pool[id].RPN[i+a+b+2] = ""
					/* после записи в очередь было
					c.Pool[id].RPN[i] = ""
					c.Pool[id].RPN[i+1] = ""
					c.Pool[id].RPN[i+2] = ""
					*/
				}
			}
		}

		if i+3 == len(c.Pool[id].RPN) {
			break
		}
	}
}

// функция проверяет, что данное арифметическое выражение в виде ОПН посчитана
func (c *Calc) rpn_is_finished(id int) bool {
	for i := 1; i < len(c.Pool[id].RPN); i++ {
		if c.Pool[id].RPN[i] != "" {
			return false
		}
	}

	return true
}

// функция для вывода в консоль содержимого очереди задач для агента
func (c *Calc) show_queue() {
	for i := 0; i < c.Queue.Len(); i++ {
		if item, ok := c.Queue.Dequeue().(Task_agent); ok {
			fmt.Println(item)
			c.Queue.Enqueue(item)
		} else {
			break
		}
	}
}

func (c *Calc) Calculate(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path != "/api/v1/calculate" {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var task Task_add
	err := json.NewDecoder(req.Body).Decode(&task)
	defer req.Body.Close()
	if err != nil { // что-то пошло не так
		fmt.Println(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else {
		fmt.Printf("%d, '%s'\n", task.Id, task.Expr)
		rpn_ := rpn.FromInfics(task.Expr)
		if (len(rpn_) == 1) && (rpn_[0] == "") { // невалидные данные
			fmt.Println("Getting data is an invalid JSON encoding")
			http.Error(w, "", http.StatusUnprocessableEntity)
		} else { // нужен механизм для генерации уникальных id
			fmt.Println(rpn_)
			t := Description{Expression: task.Expr, RPN: rpn_, Status: wait, Result: 0}
			c.Pool[task.Id] = t
			c.rpn_to_tasks(task.Id)
			fmt.Println(c.Pool[task.Id].RPN)
			// for i, s := range c.Pool[task.Id].RPN {
			// 	fmt.Printf("%d: '%s'\n", i, s)
			// }
			// c.show_queue()
			// c.Pool[task.Id].Expression = task.Expr ???
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "")
		}
	}
}

/*
Тело ответа
```

	{
	    "expressions": [
	                {
	                      "id": <идентификатор выражения>,
	                      "status": <статус вычисления выражения>,
	                      "result": <результат выражения>
	                },
	                {
	                      "id": <идентификатор выражения>,
	                      "status": <статус вычисления выражения>,
	                      "result": <результат выражения>
	                 }
	            ]
	}

```
*/
type exprs struct {
	Id     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

func (c *Calc) Expressions(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path != "/api/v1/expressions" {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	t := make([]exprs, 0)
	a := exprs{}
	for i := range c.Pool {
		// fmt.Println(i, c.Pool[i])
		a.Id = i
		a.Status = c.Pool[i].Status
		a.Result = c.Pool[i].Result
		t = append(t, a)
	}
	r := map[string][]exprs{"expressions": t}

	// на 28.05.24 работают оба варианта
	// №1
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// №2
	// jsonResp, err := json.Marshal(r)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
	// w.Write(jsonResp)
}

/*
	Коды ответа:

200 - успешно получено выражение
404 - нет такого выражения
500 - что-то пошло не так
Тело ответа
```

	{
	       "expression":
	             {
	                    "id": <идентификатор выражения>,
	                    "status": <статус вычисления выражения>,
	                    "result": <результат выражения>
	              }
	}

```
*/
func (c *Calc) Expression_id(w http.ResponseWriter, req *http.Request) {
	ids := chi.URLParam(req, "id")
	// fmt.Println("id =", ids, "type:", reflect.TypeOf(ids))
	id, err := strconv.Atoi(ids)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	task, ok := c.Pool[id]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	a := exprs{Id: id, Status: task.Status, Result: task.Result}
	err = json.NewEncoder(w).Encode(map[string]exprs{"expression": a})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
	Коды ответа:

200 - успешно получена задача
404 - нет задачи
500 - что-то пошло не так
Тело ответа
```

	{
	       "task":
	             {
	                   "id": <идентификатор задачи>,
	                   "arg1": <имя первого аргумента>,
	                   "arg2": <имя второго аргумента>,
	                   "operation": <операция>,
	                   "operation_time": <время выполнения операции>
	              }
	}

```
*/
func (c *Calc) Task_get(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path != "/internal/task" {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if c.Queue.Len() == 0 { // если очередь задач пуста
		http.Error(w, "", http.StatusNotFound)
		return
	}
	/*for k, v := range c.Pool { // проходим по всем задачам
		if v.Status != wait { // если задача вычисляется или уже завершена, то пропускаем её
			continue
		}
		// здесь уже разбиение задачи на операции
		fmt.Println(k)
	}
	http.Error(w, "", http.StatusNotFound)*/
	if task, ok := c.Queue.Dequeue().(Task_agent); ok {
		w.Header().Set("Content-Type", "application/json")
		r := map[string]Task_agent{"task": task}
		err := json.NewEncoder(w).Encode(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			// изменяем статус ОПН с "в очереди" на "вычисляется"
			_id := strings.Split(task.Id, ".")[0]
			if id, err := strconv.Atoi(_id); err != nil {
				fmt.Println("Ошибка получения номера ОПН из номера задачи")
			} else {
				// c.Pool[id].Status = calculate - так не работает, потому что ... https://stackoverflow.com/questions/42605337/cannot-assign-to-struct-field-in-a-map
				if d, ok := c.Pool[id]; ok {
					d.Status = calculate
					c.Pool[id] = d
				}
			}
		}
	} else {
		http.Error(w, "", http.StatusNotFound)
	}
}

/*
Коды ответа:
200 - успешно записан результат
404 - нет такой задачи
422 - невалидные данные
500 - что-то пошло не так
*/
func (c *Calc) Task_result(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path != "/internal/task" {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var result Result_get
	err := json.NewDecoder(req.Body).Decode(&result)
	defer req.Body.Close()
	if err != nil { // что-то пошло не так
		fmt.Println(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else {
		fmt.Println("Сервер, получен результат:", result)
		ids := strings.Split(result.Id, ".")
		if len(ids) == 1 {
			http.Error(w, "", http.StatusUnprocessableEntity)
		} else {
			_i, _p := ids[0], ids[1]
			if id, err := strconv.Atoi(_i); err != nil {
				fmt.Println("Ошибка получения номера ОПН из номера задачи")
			} else {
				if p, err := strconv.Atoi(_p); err != nil {
					fmt.Println("Ошибка получения позиции результат из номера задачи")
				} else {
					c.Pool[id].RPN[p] = fmt.Sprintf("%f", result.Result)
					fmt.Println(c.Pool[id].RPN)
					if c.rpn_is_finished(id) {
						if d, ok := c.Pool[id]; ok {
							d.Status = finished
							d.Result = result.Result
							c.Pool[id] = d
						}
					} else {
						// перестройка ОПН
						c.rpn_to_tasks(id)
					}
				}
			}
			fmt.Fprintf(w, "")
		}
		/*task, ok := c.Pool[result.Id]
		if !ok {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		task.Status = finished
		task.Result = result.Result
		c.Pool[result.Id] = task
		fmt.Printf("id = %d: %s = %f, %s\n", result.Id, task.Expression, task.Result, task.Status)*/
	}
}

// var tpl = template.Must(template.ParseFiles("index.html"))

// var tpl = template.Must(template.ParseFiles("..\\..\\front\\index.html"))

func Index(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path != "/" {
		http.NotFound(w, req)
		return
	}
	// fmt.Fprint(w, readContent("index.html"))
	// fmt.Fprintf(w, "Привет!\nЭто страница index")

	fmt.Println(os.Getwd())
	tpl := template.Must(template.ParseFiles("index.html"))

	tpl.Execute(w, nil)
}

/*
readContent(filename string) string - принимает на вход путь к файлу, а возвращает его содержимое.
В случае любой ошибки возвращайте пустую строку.
*/
func readContent(filename string) string {
	dir, err := os.Executable() // расположение исполняемого файла
	if err != nil {
		return err.Error()
	}
	// fmt.Println(filepath.Dir(dir) + "\\..\\..\\front")
	// return dir
	dir = filepath.Dir(dir) + "\\..\\..\\front\\"
	// dir = filepath.Dir(dir) + "\\..\\..\\front\\"
	f, err := os.ReadFile(dir + filename)
	if err != nil {
		// обработка ошибки
		// fmt.Println(err)
		return err.Error()
	}
	return string(f)
}
