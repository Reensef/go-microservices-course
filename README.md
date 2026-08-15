# go-microservices-course

![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/Reensef/5f118e68acb9cb3783668ea0905b45a2/raw/coverage_badge_info.json)

Для того чтобы вызывать команды из Taskfile, необходимо установить [Taskfile CLI](https://taskfile.dev)

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`)
  - Линтинг кода
  - Проверка безопасности
  - Выполняется автоматическое извлечение версий из Taskfile.yml
