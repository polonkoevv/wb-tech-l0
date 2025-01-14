# Сервис предоставляющий информацию о заказе

Демонстрационный сервис с простейшим интерфейсом для получения данных по некоторому uid-заказа.

## Требования ⚙️
Для запуска этого сервиса вас потребуется следующее:
- Docker
- Создать в корневой директории файл .env и заполнить по шаблону:
```
#POSTGRES CONFIG
DBUSER=user
DBPASSWORD=password
DBHOST=postgres
DBPORT=5432
DATABASE=db
DBSSLMODE=disable

#SERVER CONFIG
LEVEL=local
HTTP_HOST=0.0.0.0
HTTP_PORT=8080

#NATS
CLUSTER_ID=test-cluster
CLIENT_ID=test-client
LISTEN_CHANNEL=orders
LISTEN_URL=http://nats-streaming:4222/
```

#### Миграции для базы данных находятся в директории [/migration](./migrations)

## Запуск 🔧

Для запуска выполните в терминале команду ```make compose-up```, после чего сервер будет запущен на localhost на указанном
вами порту.
Для остановки сервера нужно прописать команду ```make compose-down```

## Интерфейс 🌐
После успешного запуска и перехода по пути ```http://localhost:8080/``` (если указан тот же порт, что и в шаблоне) 
у вас откроется страница, где вы можете: 
1) ввести номер заказа и получить информацию по заказу
2) получить информацию о всех имеющихся в БД заказах

Также доступны ручки по путям ```http://localhost:8080/order``` и ```http://localhost:8080/order/${orderUID}``` для получения всех заказов и заказа по определенному uid соответсвенно.

## Публикация в канал nats-streaming
В директории есть скрипт [send_data.go](./send_data.go), с помощью него можно реализовать запись в nats-streaming канал
(считывает данные о заказе из файла [model.json](./model.json))