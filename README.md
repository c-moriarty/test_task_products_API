# test_task_products_API
Программа предоставляет REST API для управления продуктами и единицами измерения. Она позволяет выполнять операции CRUD как для продуктов, так и для единиц измерения.

Запуск сервера
Для запуска программы необходимо указать путь до сервера в файле main.go. Убедитесь, что сервер запущен и доступен по указанному адресу (по умолчанию используется http://localhost:8080/).


API для продуктов

Продукты имеют следующие параметры:
id: уникальный идентификатор продукта.
name: Название продукта.
quantity: Количество продукта.
unit_cost: Цена за единицу продукта.
measure: ID единицы измерения.

Получить все продукты
Запрос для получения списка всех продуктов:

curl.exe -X GET http://localhost:8080/product/


Получить продукт по ID
Запрос для получения информации о продукте с определённым ID (например, 8):

curl.exe -X GET http://localhost:8080/product/8


Создать новый продукт
Запрос для создания нового продукта:

Invoke-WebRequest -Uri http://localhost:8080/product/ -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name":"New Product", "quantity":20, "unit_cost":15.1, "measure":2}'

Обновить продукт
Запрос для обновления информации о продукте с определённым ID (например, 8):

Invoke-WebRequest -Uri http://localhost:8080/product/8 -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"name":"Updated Product", "quantity":30, "unit_cost":25.5, "measure":1}'

Параметры: Те же, что и при создании нового продукта.

Удалить продукт
Запрос для удаления продукта с определённым ID (например, 8):

Invoke-WebRequest -Uri http://localhost:8080/product/8 -Method DELETE


API для единиц измерения

Единицы измерения имеют следующие параметры:
id: уникальный идентификатор единицы измерения.
name: Название единицы измерения.

Получить все единицы измерения
Запрос для получения списка всех единиц измерения:

curl.exe -X GET http://localhost:8080/measure/

Получить единицу измерения по ID
Запрос для получения информации о единице измерения с определённым ID (например, 2):

curl.exe -X GET http://localhost:8080/measure/2

Создать новую единицу измерения
Запрос для создания новой единицы измерения:

Invoke-WebRequest -Uri http://localhost:8080/measure/ -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name":"New Measure"}'

Обновить единицу измерения
Запрос для обновления информации о единице измерения с определённым ID (например, 1):

Invoke-WebRequest -Uri http://localhost:8080/measure/1 -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"name":"Updated Measure"}'

Удалить единицу измерения
Запрос для удаления единицы измерения с определённым ID (например, 1):

Invoke-WebRequest -Uri http://localhost:8080/measure/1 -Method DELETE
