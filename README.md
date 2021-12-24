request with init data and write it to DB:

```
http://localhost:8080/v1/pullrate/28/10/2021
```

application front:

```
http://localhost:8080

```

".config" file example:

```
postgres:
  PostgresqlHost: db
  PostgresqlPort: "5432"
  PostgresqlUser: "db_user"
  PostgresqlPassword: "<changeme>"
  PostgresqlDbname: "currency"
  PostgresqlSSLMode: false
redis:
  RedisHost: redis
  RedisPort: "6379"
  RedisUsername: ""
  RedisPassword: "<changeme>"
```

ENVIRONMENT variables:

```
ENV=PRODUCTION
```

uses .config

```
ENV=STAGING
```

uses .stageconfig

```
ENV=DEV or ENV=""
```

uses .devconfig

.env file example for docker-compose:

```
REDIS_PASSWORD=
POSTGRES_PASSWORD=
ACCESS_SECRET=
REFRESH_SECRET=
```

wait script for docker-compose:
```
mkdir -p ./scripts && cd ./scripts && wget https://github.com/adrian-gheorghe/wait-go/releases/download/1.0.0/wait-go-linux -P ./scripts -O wait
```
