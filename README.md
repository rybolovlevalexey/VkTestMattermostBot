## Работа с tarantool
#### 1. Предварительно нужно выполнить
```bash
sudo apt-get install libssl-dev pkg-config
```
Чтобы не было проблем с подключением через SSL

#### 2. Должен быть установлен tarantool 
```bash
sudo apt-get install tarantool

```
#### 3. Необходимо запустить файл с настройками базы данных (позже скорее всего будет через Docker)
```bash
tarantool /internal/database/init.lua 
```



## Запуск Mattermost локально через Docker для тестирования

### Требования
- Установленный Docker
- Docker добавлен в группу пользователя либо запускаются через sudo

```bash
docker run -d \
  --name mattermost-preview \
  -p 8065:8065 \
  -p 8067:8067 \
  --restart always \
  -v mm-preview-data:/mattermost/data \
  -v mm-preview-config:/mattermost/config \
  mattermost/mattermost-preview
```

### При необходимости остановить и заново запустить использовать

```bash
docker stop mattermost-preview
docker start mattermost-preview
```
Данный вариант позволит вновь воспользоваться результатами предыдущих действий, например, сохранится информация о созданных ранее ботах и информация о регистрации пользователя.
