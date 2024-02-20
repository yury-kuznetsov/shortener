#!/bin/bash

# Скрипт для генерации полезной нагрузки
#
# Пример вызова для Apache Bench (ab):
# $ ab -n 10000 -c 100 -p bench/encode_data.json -T 'text/plain' http://localhost:8080/
#
# Команда отправляет post-запрос на локальный хост 10_000 раз,
# выполняя 100 запросов одновременно, используя файл `bench/encode_data.json`
# в качестве тела запроса.

for i in {1..10}
do
  # DecodeHandler
  ab -n 10 -c 10 http://localhost:8080/uGjeYcOT
  # EncodeHandler
  ab -n 10 -c 10 -p bench/encode_data.json -T 'text/plain' http://localhost:8080/
  # EncodeBatchHandler
  ab -n 10 -c 10 -p bench/encode_batch_data.json -T 'application/json' http://localhost:8080/api/shorten/batch
  # EncodeJSONHandler
  ab -n 10 -c 10 -p bench/encode_json_data.json -T 'application/json' http://localhost:8080/api/shorten
  # PingHandler
  ab -n 10 -c 10 http://localhost:8080/ping
  # UserUrlsHandler
  ab -n 10 -c 10 http://localhost:8080/api/user/urls

  sleep 1
done