# Сборка образа
из главной директории
docker build -t pinger-service -f pinger/Dockerfile .

# Запуск контейнера
docker run -d --name pinger-container pinger-service
! Не работает


