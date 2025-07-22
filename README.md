# go-microservices-course

![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/Reensef/5f118e68acb9cb3783668ea0905b45a2/raw/coverage_badge_info.json)

Этот репозиторий содержит проект из курса [Микросервисы, как в BigTech 2.0](https://olezhek28.courses/microservices) от [Олега Козырева](http://t.me/olezhek28go).

Для того чтобы вызывать команды из Taskfile, необходимо установить [Taskfile CLI](https://taskfile.dev)

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Линтинг кода
  - Проверка безопасности
  - Выполняется автоматическое извлечение версий из Taskfile.yml
