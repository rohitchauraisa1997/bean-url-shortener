# bean-url-shortener

This is a simple URL shortener service that allows you to shorten long URLs into shorter, more manageable links. It's built using react as frontend, bean as the backend server and mysql/redis as the database.

You must have commonly observed these shortened links when you share links accross social apps or posts on LinkedIn.

Features

1. Shorten long URLs into concise, easy-to-share links.
2. Redirect users to the original URL when they access the shortened link.
3. Track the number of clicks on each shortened link.
4. Prevents DDOS and bot attacks by allocating a QUOTA for each user.
5. Admin UI to track all users analytics.

Setup

Since the project is setup using docker compose. Once the keys are added, the project should ideally work with a simple

```
    docker compose up --build
```

To Connect to redis and mysql databases

```
docker compose exec -it database-redis redis-cli -p 6379

docker compose exec -it database-mysql  mysql -u root local -p
```

Quick Demo:-

https://github.com/user-attachments/assets/8fe6e243-73f9-4dc3-9bbc-0d21d308a25c
