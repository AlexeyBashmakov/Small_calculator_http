rem два следующих синтаксиса не работают при перечаде json в винде
rem curl --location "http://localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{"id": 1, "expression": "1+1"}" -v
::curl --location "http://localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{'id': 1, 'expression': '1+1'}" -v

echo "Getting result of calculation"
curl --location "http://localhost:8080/internal/task" --header "Content-Type: application/json" --data "{\"id\": 1, \"result\": 2.5}"
