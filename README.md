# metida - golang service example

Данный проект это небольшой пример реализации REST API сервиса, в котором можно увидеть:
1. Реализацию аутентификации и авторизации ([JWT](https://gist.github.com/zmts/802dc9c3510d79fd40f9dc38a12bccfc))
2. [Генерацию swagger](https://github.com/swaggo/swag) документации
3. [Генерацию sql](https://docs.sqlc.dev/en/latest/index.html) запросов и выполнение их к базе данных 
4. [Логирование](https://github.com/uber-go/zap)
5. [Метрики](https://github.com/prometheus/client_golang)
6. Использование контейнеризации
7. И что то еще по мелочи, вроде реализации graceful shutdown, запуск приложения для dev и prod окружения...

Для реализации основной логики использовался фреймворк Gin, но для следующих версий API можно безболезненно
брать другие библиотеки.

### Описание работы сервиса:
Пользователь может зарегистрироваться под ролью user, войти в свой аккаунт, и обновить Access Token при помощи Refresh Token.
За это отвечают следующие методы:
```go
// @Router /auth/sign-up [post]
func (o *UserHandler) CreateUser(c *gin.Context) {...}

// @Router /auth/sign-in [post]
func (o *UserHandler) AuthenticationUser(c *gin.Context) {...}

// @Router /auth/refresh [post]
func (o *UserHandler) RefreshTokenUser(c *gin.Context) {...}
```

Далее авторизованный пользователь может создавать и просматривать свой [Content](https://github.com/Dsmit05/metida/blob/master/internal/models/user.go#L15):
```go
// @Router /lk/content [POST]
func (o *UserContent) CreateContent(c *gin.Context) {...}

// @Router /lk/content/{id} [GET]
func (o *UserContent) ShowContent(c *gin.Context) {...}
```

и у данного сервиса есть свой [Blog](https://github.com/Dsmit05/metida/blob/master/internal/models/user.go#L8), который может посмотреть каждый:
```go
// @Router /blog/{id} [GET]
func (o *SiteBlog) ShowBlog(c *gin.Context) {...}
```
Но создавать данные записи может только юзер с правами admin
```go
// @Router /lk/blog [POST]
func (o *SiteBlog) CreateBlog(c *gin.Context) {
	middlewares.CheckAccessRights(c, consts.RoleAdmin)
    ...}
```

### Запуск сервиса
Все основные команды можно увидеть в [Makefile](https://github.com/Dsmit05/metida/blob/master/Makefile).
Для быстрого запуска выполните команду `docker-compose up`

и перейдите к swagger документации: http://localhost:8081/swagger/index.html,
для логина под ролью admin войдите как: `email: admin, password: admin`


на порту 8081 расположен debag сервер:

- Профилировщик:
http://localhost:8081/pprof/

- Метрики для Prometheus:
http://localhost:8081/metrics/

- Информация о сборке:
http://localhost:8081/

Так же вы можете запустить для данного сервиса frontend часть: https://github.com/Dsmit05/metida-ui.
И посмотреть полноценную реализацию аутентификации и авторизации.

### Окружение приложения
Для отображения метрик и логов вы можете запустить и настроить Grafana и ELK из папки env-apps

### Что нужно улучшить:
- [ ] для production версии все настройки надо брать из защищенного места(к примеру consul)
- [ ] для работы с postrgre использовать пул и выполнить более качественную обработку ошибок
- [ ] использовать [кеш](https://github.com/Dsmit05/metida/blob/master/pkg/cache/lru/lru-cache.go) для частых запросов к бд
- [ ] не использовать в контейнерах network_mode: host
- [ ] оставлять более подробные комментарии к функциям

### P.S.
В папке scripts приложен пример по работе с миграциями, 
т.к. подходов для миграций много(к примеру через консольные приложения) включать это в основное приложение
я посчитал не нужным.
Данный проект носит ознакомительный характер, используйте его на свой страх и риск.
