# docker-compose-deploy

Zero-downtime docker-compose deploy


```bash
curl -O https://raw.githubusercontent.com/ViBiOh/docker-compose-deploy/master/deploy.sh
chmod +x deploy.sh

./deploy.sh awesome_project $(git rev-parse --short HEAD) path_to_your_compose_default_pwd
```
