docker-compose down 
docker container stop $(docker container ls -aq)
docker container stop $(docker container ls -aq)
docker system prune
docker volume prune
