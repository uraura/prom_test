# prom_test

prometheusのremote_writeの挙動を確認したいだけのやつ

```shell
docker compose up -d
```

gen <-(scrape)- left -(remote_write)-> right