# Variable service

Variable service in golang for vs code

## Quick start for Windows

```
git clone https://github.com/sladeshka/variable-go.service.git
cd variable-go.service
cp docker/local/.env.example docker/local/.env
docker-compose -f /docker/local/Dokerfile --env-file /docker/local/.env build --no-cache
docker-compose -f /docker/local/Dokerfile --env-file /docker/local/.env up -d --force-recreate 
```