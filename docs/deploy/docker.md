# Deploying with docker

## building the img manually

```bash
git clone "https://github.com/toggler-io/toggler.git"
cd toggler
docker build . -t toggler
docker run --rm -d -p 8080:8080 -e DATABASE_URL=${DATABASE_URL} --name toggler toggler
```
