# sverk

> /ˈsvɜːrk/, personal checkin script, power by golang

## Features
- [✓] [Hifini](https://www.hifini.com/)
- [✓] [V2ex](https://www.v2ex.com/)

## Docker

### Usage

```bash
docker run -d --name sverk svefn/sverk:latest
mkdir -p $PWD/sverk && docker cp sverk:/app/conf/config.yml $PWD/sverk/config.yml
docker stop sverk && docker rm sverk
docker run -d \
  --name sverk \
  --restart=always \
  -v $PWD/sverk:/app/conf \
  svefn/sverk:latest
```

### Update

```bash
docker stop sverk && docker rm sverk && docker pull svefn/sverk:latest
docker run -d \
  --name sverk \
  --restart=always \
  -v $PWD/sverk:/app/conf \
  svefn/sverk:latest
```

### Delete

```bash
docker stop sverk && docker rm sverk && docker rmi svefn/sverk && rm -rf $PWD/sverk
```

## Manual

### Usage

```bash
git clone h github.com/svefnz/sverk
cd sverk
cp conf/config.yml.example conf/config.yml
# edit conf/config.yml
go build .
./sverk serve # run service background
./sverk start # run service once
./sverk -s xxx # run named service
```
