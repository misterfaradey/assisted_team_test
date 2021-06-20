## Тестовое задание в команду ассистеда (Python)

В папке два XML – это ответы на поисковые запросы, сделанные к одному из наших партнёров.
В ответах лежат варианты перелётов (тег `Flights`) со всей необходимой информацией,
чтобы отобразить билет на Aviasales.

На основе этих данных, нужно сделать вебсервис,
в котором есть эндпоинты, отвечающие на следующие запросы:

* Какие варианты перелёта из DXB в BKK мы получили? 
* Самый дорогой/дешёвый, быстрый/долгий и оптимальный варианты
* В чём отличия между результатами двух запросов (изменение маршрутов/условий)?

Язык реализации: `python3`
Формат ответа: `json`
Используемые библиотеки и инструменты — всё на твой выбор.

Оценивать будем умение выполнять задачу имея неполные данные о ней,
умение самостоятельно принимать решения и качество кода.

## Ответ:

```
---для запуска под Linux---
1.  dep ensure -vendor-only
2.  cp configs/config.toml.dist configs/config.toml
3.  go build -o ./assisted_team_api cmd/assisted_team/main.go
4.  ./assisted_team_api

---для запуска под Windows---
1.  dep ensure -vendor-only
2.  copy configs\config.toml.dist configs\config.toml
3.  go build -o ./assisted_team_api.exe cmd/assisted_team/main.go
4.  ./assisted_team_api.exe

```
* Какие варианты перелёта из DXB в BKK мы получили?  
```http request
GET http://0.0.0.0:8080/api/flights/all?from=DXB&to=BKK
GET http://0.0.0.0:8080/api/flights/all?from=DXB&to=BKK&oneway=true

```
сортировка по убыванию цены SingleAdult
```http request
GET http://0.0.0.0:8080/api/flights/all?from=DXB&to=BKK&price=1&type=SingleAdult&oneway=true
```
сортировка по возрастанию цены SingleChild
```http request
GET http://0.0.0.0:8080/api/flights/all?from=DXB&to=BKK&price=1&type=SingleChild&oneway=true
```
* Самый дорогой/дешёвый, быстрый/долгий и оптимальный варианты  

**самый дорогой вариант**
```http request
GET http://0.0.0.0:8080/api/flight?from=DXB&to=BKK&sort=max_price&oneway=false
```
**самый дешевый вариант**
```http request
GET http://0.0.0.0:8080/api/flight?from=DXB&to=BKK&sort=min_price&oneway=true
```
**максимальное время**
```http request
GET http://0.0.0.0:8080/api/flight?from=DXB&to=BKK&sort=max_time&oneway=false
```
**минимальное время**
```http request
GET http://0.0.0.0:8080/api/flight?from=DXB&to=BKK&sort=min_time&oneway=true
```
**оптимальный вариант: примем за формулу оптимальности стоимость часа перелета как 5 SGD (1h = 5SGD ≈ 235RUB)**
```http request
GET http://0.0.0.0:8080/api/flight?from=DXB&to=BKK&sort=optimal&oneway=false
```
* В чём отличия между результатами двух запросов (изменение маршрутов/условий)?  

**в файле `RS_ViaOW.xml` перелеты в одну сторону. а в `RS_Via-3.xml` с возвратом**
