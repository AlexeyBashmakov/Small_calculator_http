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
Коды ответа:
201 - выражение принято для вычисления
422 - невалидные данные
500 - что-то пошло не так
Тело ответа
```
{
}

```

 - Получение списка выражений
```
curl --location 'localhost/api/v1/expressions'
```
Коды ответа:
200 - успешно получен список выражений
500 - что-то пошло не так
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

 - Получение задачи для выполнения.
```
curl --location 'localhost/internal/task'
```
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

 - Прием результата обработки данных.
```
curl --location 'localhost/internal/task' \
--header 'Content-Type: application/json' \
--data '{
      "id": 1,
      "result": 2.5
}'
```
Коды ответа:
200 - успешно записан результат
404 - нет такой задачи
422 - невалидные данные
500 - что-то пошло не так

Время выполнения операций задается переменными среды в милисекундах
TIME_ADDITION_MS - время выполнения операции сложения в милисекундах
TIME_SUBTRACTION_MS - время выполнения операции вычитания в милисекундах
TIME_MULTIPLICATIONS_MS - время выполнения операции умножения в милисекундах
TIME_DIVISIONS_MS - время выполнения операции деления в милисекундах

### Агент
Демон, который получает выражение для вычисления с сервера, вычисляет его и отправляет на сервер результат выражения.
При старте демон запускает несколько горутин, каждая из которых выступает в роли независимого вычислителя. 
Количество горутин регулируется переменной среды
COMPUTING_POWER
Агент обязательно общается с оркестратором по http
Агент все время приходит к оркестратору с запросом "дай задачку поработать" (в ручку GET internal/task для получения задач).
Оркестратор отдаёт задачу.
Агент производит вычисление и в ручку оркестратора (POST internal/task для приёма результатов обработки данных) отдаёт результат.
