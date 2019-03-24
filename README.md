# docker-compose-deploy

Zero-downtime docker-compose deploy


```bash
curl -O https://raw.githubusercontent.com/ViBiOh/docker-compose-deploy/master/deploy.sh
source deploy.sh

deploy awesome_project $(git rev-parse --short HEAD) path_to_your_compose_default_pwd
```
