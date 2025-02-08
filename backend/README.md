// ! Удалить

# Сборка образа
из главной директории
docker build -t docker_monitoring -f backend/Dockerfile .

# Запуск контейнера
docker run -p 8080:8080 -e DB_HOST=host.docker.internal docker_monitoring
! Не работает


```ls -R | less```
Показывает рекурсивынй список файлов