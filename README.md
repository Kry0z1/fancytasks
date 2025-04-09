# fancytasks

## Что это
Это попытка переписать приложение для трекинга задач на более сложном стэке

## Используемые технологии
 - redis для кэширования
 - postgresql для хранения данных
 - docker-compose
 - 'голая' библиотека `net/http` для обработки хэндлеров и мидлварей
 - 'голая' библиотека `database/sql` для работы с базой данных

## Запуск
```bash
git clone https://github.com/Kry0z1/fancytasks.git
cd fancytasks
docker compose up --build
```

Для работы необходимо в переменных окружения или в `.env` в корне проекта выставить следующие переменные:
 - REDIS_PASS=пароль для редиса 
 - POSTGRES_PASS=пароль для постгреса
 - POSTGRES_DB=имя бд
 - PG_ADMIN_DEFAULT_EMAIL=почта для pgadmin
 - PG_ADMIN_DEFAULT_PASS=пароль от pgadmin
 - DATA_PATH=директория для хранения данных

По адресу `localhost:15433` будет доступна панель доступа pgAdmin

По адресу `localhost:8000` будет доступно само приложение

По адресу `localhost:15432` будет доступна база данных

# .env добавлен для более легкого запуска