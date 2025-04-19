# go-lab-stresstest

### How to use
```
docker run --rm tel33z/fc-stresstest  --url=https://google.com --requests=100 --concurrency=10
```

### How to build local image 
```
docker build -t stresstest .
docker run --rm stresstest --url=http://google.com --requests=100 --concurrency=10
```
