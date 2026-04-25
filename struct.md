```text
cmd/<service>/main.go            # точка входа процесса

internal/
  app/                         # bootstrap, wiring, routing
  behavior/                    # lifecycle/readiness-контракты
  config/                      # env-конфигурация
  model/
    request.go                 # transport input
    response.go                # transport output
    dto/                       # межслойные структуры
    dao/                       # структуры для интеграций
  provider/                    # infra-клиенты (HTTP/Auth/Broker)
  repository/
    domain/                    # адаптеры к сервисам проекта
    integrations/              # адаптеры к внешним системам
  server/
    container.go
    http/
      server.go                # HTTP server lifecycle
      handler/                 # публичные HTTP handlers
      middleware/              # version/cors/logger/counter/auth
  service/                     # use-case orchestration
  vars/                        # константы и общие ошибки

pkg/
  utils/                       # signal/sleep и др.

deploy/                        # окружения и k8s values
docs/                          # документация/диаграммы
```