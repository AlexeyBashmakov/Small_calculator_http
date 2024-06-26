# Small_calculator_http
## Распределенный вычислитель арифметических выражений

Создан на Go как один из этапов обучения в Yandex лицее.

Идея вычислителя:<br>
Пользователь хочет считать арифметические выражения. Он вводит строку `2 + 2 * 2` и хочет получить в ответ `6`.<br>
Но наши операции сложения и умножения (также деления и вычитания) выполняются "очень-очень" долго.<br>
Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, невозможна.<br>
Более того, вычисление каждой такой операции в нашей "альтернативной реальности" занимает "гигантские" вычислительные мощности.<br>
Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "машин".<br>
Поэтому пользователь может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"?<br>
Если выражение наконец будет вычислено - то он получит результат.<br>
Некоторые части арфиметического выражения мы вычисляем параллельно.

Вычислитель должен состоять из двух частей: front-end'а и back-end'а.<br>
Front-end на данном этапе не реализован.

Вычислитель умеет считать арифметические выражения содержащие четыре арифметические операции (+, -, *, /) и скобки.
## Установка проекта
Для установки проекта на локальный компьютер возможны два варианта:
1) склонировать репозиторий командой<br>
git clone https://github.com/AlexeyBashmakov/Small_calculator_http.git<br>
появится папка `Small_calculator_http`.
2) скачать проект с репозитория как архив zip и разархивировать, появится папка `Small_calculator_http`.

Далее работа с проектом в обоих вариантах установки совпадает, но различается в зависимости от операционной системы.<br>
В директории `Small_calculator_http/calculator/cmd/calculator` находится файл `main.go`, который является
точкой входа в приложение. Также здесь находятся файлы с расширениями `.bat` и `.sh`. Первые используются для работы с приложением в ОС Windows, вторые - в ОС Linux.
### ОС Windows
Для того чтобы проект можно было запустить на исполнение необходимо в папке `Small_calculator_http/calculator/cmd/calculator`
запустить скрипт `build.bat` при этом go докачает пакет `github.com/go-chi/chi/v5` и соберёт исполняемый файл `main.exe`<br>
Для запуска проекта и работы с ним потребуются две консоли (запустить их можно с помощью комбинаций клавиш Win+R и, 
в появившемся окне "Выполнить", в строке "Открыть" ввести `cmd`).<br> 
В одной консоли, в директории `Small_calculator_http/calculator/cmd/calculator`, запускаем скрипт `run.bat`. 
(В этом скрипте устанавливаются переменные окружения необходимые для работы программы.) 
При старте серверной части в консоли отобразятся приветствие Чебурашки и информационные сообщения о 
старте вычислительных горутин агента. После старта сервер ждет запросов. Для завершения 
работы сервера в консоли достаточно нажать комбинацию клавиш `Ctrl+C`.<br>
В другой консоли, в той же директории, предлагается запуск скриптов (для первичного, демонстрационного 
ознакомления) в следующей последовательности:
 - `calculate.bat`
 - `expression_id.bat`
 - `expression.bat`

Скрипт `calculate.bat` содержит набор арифметических выражений для передачи серверу с помощью утилиты `curl`.<br>
Скрипт `expression_id.bat` содержит запрос к серверу о статусе выполнения отдельной задачи, также с помощью утилиты `curl`..<br>
Скрипт `expressions.bat` содержит запрос к серверу о статусе выполнения всех задач переданных серверу, также с помощью утилиты `curl`.
### ОС Linux
Для того чтобы проект можно было запустить на исполнение необходимо в папке `Small_calculator_http/calculator/cmd/calculator`
запустить скрипт `build.sh` при этом go докачает пакет `github.com/go-chi/chi/v5` и соберёт исполняемый файл `main`<br>
Для запуска проекта и работы с ним потребуются две консоли.<br> 
В одной консоли, в директории `Small_calculator_http/calculator/cmd/calculator`, запускаем скрипт `run.sh`. 
(В этом скрипте устанавливаются переменные окружения необходимые для работы программы.) 
При старте серверной части в консоли отобразятся приветствие Чебурашки и информационные сообщения о 
старте вычислительных горутин агента. После старта сервер ждет запросов. Для завершения 
работы сервера в консоли достаточно нажать комбинацию клавиш `Ctrl+C`.<br>
В другой консоли, в той же директории, предлагается запуск скриптов (для первичного, демонстрационного 
ознакомления) в следующей последовательности:
 - `calculate.sh`
 - `expression_id.sh`
 - `expression.sh`

Скрипт `calculate.sh` содержит набор арифметических выражений для передачи серверу с помощью утилиты `wget`.<br>
Скрипт `expression_id.sh` содержит запрос к серверу о статусе выполнения отдельной задачи, также с помощью утилиты `wget`..<br>
Скрипт `expressions.sh` содержит запрос к серверу о статусе выполнения всех задач переданных серверу, также с помощью утилиты `wget`.
## Примеры использования
Файлы `calculate.bat` и `calculate.sh` содержат по заданию намеренно написанному с ошибкой для демонстрации обработки сервером синтаксически неправильных выражений:<br>
`--data "{\"id\": 2, \"expression\": \"1-1*3+3*\"}"`<br>
здесь выражение завершается символом операции и считается неправильным (невалидным). Также эти файлы содержат по заданию показывающему ограниченность текущего алгортима 
параллельных вычислений задач:<br>
`--data "{\"id\": 7, \"expression\": \"(8+2*5)/(1+3*2-4)\"}"`<br>
Текущий алгоритм расчёта возвращает результат этого выражения 2, а правильный ответ - 6.
### ОС Windows
Запросы к серверу можно осуществлять с помощью утилиты `curl`, которая входит в состав ОС Windows (по крайней мере до версии 10, включительно).
1) Постановка задачи серверу - добавление вычисления арифметического выражения<br>
`curl --location "http://localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"id\": 1, \"expression\": \"1-1*3+3*9\"}" -v`<br>
Объяснение параметров:<br>
--location - адрес, к которому отправить запрос<br>
--header - заголовок HTTP-запроса<br>
--data - данные, которые передаются с запросом<br>
2) Получение списка выражений<br>
`curl --location "http://localhost:8080/api/v1/expressions" -v`
3) Получение выражения по его идентификатору<br>
`curl --location "http://localhost:8080/api/v1/expressions/:8" -v`<br>
Здесь после двоеточия в конце строки адреса указывается идентификатор задачи (в данном случае 8), статус которой нам интересен.

### ОС Linux
В Linux'е нет, по-умолчанию, curl'а, поэтому для демонстрации работы проекта можно использовать wget, которая включается почти в каждый дистрибутив данной ОС.
1) Постановка задачи серверу - добавление вычисления арифметического выражения<br>
`wget http://localhost:8080/api/v1/calculate --header="Content-Type: application/json" --post-data="{\"id\": 3, \"expression\": \"3+4\"}" -O -`<br>
Объяснение параметров:<br>
адрес, к которому отправляется запрос, располагается сразу за именем утилиты, поэтому не выделяется никаким специальным параметром<br>
--header - заголовок HTTP-запроса<br>
--post-data - данные, которые передаются с запросом<br>
-O - полученную информацию выводить в консоль
2) Получение списка выражений<br>
`wget http://localhost:8080/api/v1/expressions -O -`
3) Получение выражения по его идентификатору<br>
`wget http://localhost:8080/api/v1/expressions/:3 -O -`<br>
Здесь после двоеточия в конце строки адреса указывается идентификатор задачи (в данном случае 3), статус которой нам интересен.

## Структура проекта
Структура каталогов:<br>
Small_calculator_http/<br>
|<br>
-calculator/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;-cmd/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;-calculator/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-build.bat<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-build.sh<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-calculate.bat<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-calculate.sh<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-expression_id.bat<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-expression_id.sh<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-expressions.bat<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-expressions.sh<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-main<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-main.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-main_test.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-run.bat<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-run.sh<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-task_get.bat<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-task_result.bat<br>
|&nbsp;&nbsp;&nbsp;&nbsp;-http/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;-agent/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;-agent.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;-server/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-handler/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;-handler.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-server.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;-internal/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;-application/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;-application.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;-constants/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;-environ_vars.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;-operations.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;-environ_vars/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;-environ_vars.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;-service/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-service.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;-pkg/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;-my_queue/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;-my_queue.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;-my_stack/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;|&nbsp;&nbsp;-my_stack.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;-rpn/<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|<br>
|&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-rpn.go<br>
|&nbsp;&nbsp;&nbsp;&nbsp;-go.mod<br>
|&nbsp;&nbsp;&nbsp;&nbsp;-go.sum<br>
-LICENSE<br>
-README.md

Пакет `calculator/cmd/calculator/main.go` является точкой входа в приложение. Здесь происходит вызов функции для проверки установки переменных среды<br>
`if !environ_vars.CheckEnvironmentVariables()`<br>
и если они не установлены, то происходит их установка значениями по-умолчанию<br>
`if environ_vars.SetEnvironmentVariables() != nil`<br>
если и это не удаётся, то приложение завершается с кодом 2. В случае успеха переходим к запуску сервера и агента<br>
`application.Run(ctx)`<br>

Остальная часть описания ждёт своего часа...

## Back-end часть
Состоит из 2 элементов:
 - Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения. Далее будем называть его "оркестратором".
 - Вычислитель, который может получить от "оркестратора" задачу, выполнить его и вернуть серверу результат. Далее будем называть его "агентом".

### Оркестратор
Сервер, который имеет следующие endpoint-ы:
 - Добавление вычисления арифметического выражения
```
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
      "id": <уникальный идентификатор выражения>,
      "expression": <строка с выражением>
}'
```
Коды ответа:<br>
201 - выражение принято для вычисления<br>
422 - невалидные данные<br>
500 - что-то пошло не так<br>
Тело ответа
```
{
}

```

 - Получение списка выражений
```
curl --location 'localhost/api/v1/expressions'
```
Коды ответа:<br>
200 - успешно получен список выражений<br>
500 - что-то пошло не так<br>
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

 - Получение выражения по его идентификатору
```
curl --location 'localhost/api/v1/expressions/:id'
```
Коды ответа:<br>
200 - успешно получено выражение<br>
404 - нет такого выражения<br>
500 - что-то пошло не так<br>
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

 - Получение задачи для выполнения.
```
curl --location 'localhost/internal/task'
```
Коды ответа:<br>
200 - успешно получена задача<br>
404 - нет задачи<br>
500 - что-то пошло не так<br>
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

 - Прием результата обработки данных.
```
curl --location 'localhost/internal/task' \
--header 'Content-Type: application/json' \
--data '{
      "id": 1,
      "result": 2.5
}'
```
Коды ответа:<br>
200 - успешно записан результат<br>
404 - нет такой задачи<br>
422 - невалидные данные<br>
500 - что-то пошло не так

Время выполнения операций задается переменными среды в милисекундах
- **TIME_ADDITION_MS** - время выполнения операции сложения в милисекундах
- **TIME_SUBTRACTION_MS** - время выполнения операции вычитания в милисекундах
- **TIME_MULTIPLICATIONS_MS** - время выполнения операции умножения в милисекундах
- **TIME_DIVISIONS_MS** - время выполнения операции деления в милисекундах

### Агент
Демон, который получает выражение для вычисления с сервера, вычисляет его и отправляет на сервер результат выражения.<br>
При старте демон запускает несколько горутин, каждая из которых выступает в роли независимого вычислителя. <br>
Количество горутин регулируется переменной среды **COMPUTING_POWER**<br>
Агент общается с оркестратором по http.<br>
Агент все время приходит к оркестратору с запросом "дай задачку поработать" (в ручку GET internal/task для получения задач).<br>
Оркестратор отдаёт задачу.<br>
Агент производит вычисление и в ручку оркестратора (POST internal/task для приёма результатов обработки данных) отдаёт результат.


