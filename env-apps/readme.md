В данном разделе находятся настройки запуска и конфигурации приложений каким либо образом связанных с нашим апи сервисом.
Это лишь дополнительные возможности которые никак не влияют на работу основного приложения.

### Метрики

Для отображения метрик нашего приложения используем связку Prometheus и Grafana.

Для запуска перейдите в папку metric и выполните команду:
```
docker compose up -d
```

Переходим по http://localhost:3000/ и подключаем Prometheus http://localhost:9090/

И импортим в Grafana model.json для отображения дашборда нашего приложения и подключаем https://grafana.com/grafana/dashboards/6671

<details>
  <summary>полезные ссылки</summary>
  - https://gist.github.com/diolavr/ef6d63288a4244b8f745958041fd3f73 - примеры клиента Prometheus
</details>

### Логи

В папке logs приведена примерная конфигурация ELK и filebeat.

Для начала запустите основное приложение в контейнере, filebeat будет забирать логи из контейнера.
В [конфиге](https://github.com/Dsmit05/metida/blob/master/env-apps/logs/filebeat/config/filebeat.yml#L11) filebeat укажите имя контейнера.

Для запуска перейдите в папку logs и выполните команду:
```
docker compose up -d
```
После перейти по http://localhost:5601/app/management/kibana/dataViews (пароль: metida)
и добавить индекс.

Данные логи можно будет посмотреть по http://localhost:5601/app/discover

Для нормального отображения нужно изменить конфигурацю logstash, т.к. это только пример, корректные настройки здесь не приведены.

Данные скрипты запускали были взяты с https://github.com/deviantony/docker-elk и немного изменены.

<details>
  <summary>полезные ссылки</summary>
  - https://www.elastic.co/guide/en/beats/filebeat/current/configuration-autodiscover-hints.html
</details>
