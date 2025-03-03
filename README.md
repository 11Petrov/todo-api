Убедитесь, что у вас установлен Go (версия 1.20 или выше)

## Настройка базы данных
Вам необходимо создать базу данных PostgreSQL и настроить подключение к ней. Для этого создайте файл .env в корневой директории проекта и добавьте в него следующие переменные:

DB_CONN_STR=postgres://username:password@localhost:port/dbname?sslmode=disable
PORT=port
Замените username, password, localhost, port и dbname на соответствующие значения для вашей базы данных.

## Запуск приложения
Для запуска приложения выполните команду:

go run cmd/todo-api/main.go

Сервер запустится на порту, указанном в переменной PORT.

## Скриншоты работоспособности:
![Скриншот успешного запроса POST](screenshots/post.png) 
![Скриншот успешного запроса UPDATE](screenshots/update.png) 
![Скриншот успешного запроса DELETE](screenshots/delete.png) 
![Скриншот успешного запроса GET](screenshots/get.png) 

![Скриншот наличия записей в базе данных](screenshots/db.png)
