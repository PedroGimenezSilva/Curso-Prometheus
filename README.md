# Prometheus
Foobar is a Python library for dealing with word pluralization.

## Realizando o download 


```bash
curl -LO https://github.com/prometheus/prometheus/releases/download/v2.38.0/prometheus-2.38.0.linux-amd64.tar.gz

```

## Após o download, vamos extrair o arquivo e acessar o diretório extraído.

```bash
tar -xvf prometheus-2.38.0.linux-amd64.tar.gz
```
## Extrair os binário para o diretório /usr/local/bin.
```bash
sudo mv prometheus-2.38.0.linux-amd64/prometheus /usr/local/bin/prometheus
sudo mv prometheus-2.38.0.linux-amd64/promtool /usr/local/bin/promtool
```

## Vamos ver se o binário está funcionando.

```bash
prometheus --version

prometheus, version 2.38.0 (branch: HEAD, revision: 818d6e60888b2a3ea363aee8a9828c7bafd73699)
  build user:       root@e6b781f65453
  build date:       20220816-13:23:14
  go version:       go1.18.5
  platform:         linux/amd64
```
## Agora vamos criar o diretório de configuração do Prometheus.

```bash
sudo mkdir /etc/prometheus

```
Vamos criar um grupo e um usuário para o Prometheus.
```bash
sudo addgroup --system prometheus
sudo adduser --shell /sbin/nologin --system --group prometheus
```

Precisamos fazer com que o Prometheus seja um serviço em nossa máquina, para isso precisamos criar o arquivo de service unit do SystemD.
```bash
sudo vim /etc/systemd/system/prometheus.service
```

Agora vamos adicionar o seguinte conteúdo ao arquivo de configuração do service unit do Prometheus: 
```bash
[Unit]
Description=Prometheus
Documentation=https://prometheus.io/docs/introduction/overview/
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=prometheus
Group=prometheus
ExecReload=/bin/kill -HUP \$MAINPID
ExecStart=/usr/local/bin/prometheus \
  --config.file=/etc/prometheus/prometheus.yml \
  --storage.tsdb.path=/var/lib/prometheus \
  --web.console.templates=/etc/prometheus/consoles \
  --web.console.libraries=/etc/prometheus/console_libraries \
  --web.listen-address=0.0.0.0:9090 \
  --web.external-url=

SyslogIdentifier=prometheus
Restart=always

[Install]
WantedBy=multi-user.target

```
Mudando a permissão dos diretórios 
```bash
sudo chown -R prometheus:prometheus /var/log/prometheus
sudo chown -R prometheus:prometheus /etc/prometheus
sudo chown -R prometheus:prometheus /var/lib/prometheus
sudo chown -R prometheus:prometheus /usr/local/bin/prometheus
sudo chown -R prometheus:prometheus /usr/local/bin/promtool
```
Vamos fazer um reload no systemd para que o serviço do Prometheus seja iniciado.
```bash
sudo systemctl daemon-reload
```
Vamos iniciar o serviço do Prometheus.


```bash
sudo systemctl start prometheus
```
Temos que deixar o serviço do Prometheus configurado para que seja iniciado automaticamente ao iniciar o sistema.
```bash
sudo systemctl enable prometheus
```
Você pode verificar nos logs se tudo está rodando maravilhosamente.
```bash
sudo journalctl -u prometheus
```

## License

[Linuxtips](https://github.com/badtuxx/DescomplicandoPrometheus)
