# Команда для старта
```NODE_OPTIONS=--openssl-legacy-provider npm start```

docker build -t frontend-app -f frontend/Dockerfile .

docker run -d -p 3000:80 frontend-app
