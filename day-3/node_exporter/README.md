## node_exporter


Intalando o node_exporter e movendo ele para o caminho correto dos binários 

```bash
wget https://github.com/prometheus/node_exporter/releases/download/v1.5.0/node_exporter-1.5.0.linux-amd64.tar.gz
cd node_exporter
mv node_exporter /usr/local/bin
```  

### Configurando o node_exporter para se  tornar um serviço 
```bash
#Criando o grupo
sudo addgroup --system node_exporter

#criando o usuário 
sudo adduser --shell /bin/nologin --system --group node_exporter

# Criando e configurando  o arquivo do serviço do node_exporter
vim /etc/systemd/system/node_exporter.service
```  

Arquivo node_exporter.service

```bash
[Unit] # Inicio do arquivo de configuração do serviço
Description=Node Exporter # Descrição do serviço
Wants=network-online.target # Define que o serviço depende da rede para iniciar
After=network-online.target # Define que o serviço deverá ser iniciado após a rede estar disponível

[Service] # Define as configurações do serviço
User=node_exporter # Define o usuário que irá executar o serviço
Group=node_exporter # Define o grupo que irá executar o serviço
Type=simple # Define o tipo de serviço
ExecStart=/usr/local/bin/node_exporter # Define o caminho do binário do serviço

[Install] # Define as configurações de instalação do serviço
WantedBy=multi-user.target # Define que o serviço será iniciado utilizando o target multi-user
```

Agora precisamos dar um daemon-reload e realizar o start do serviço node_exporter
```bash
systemctl daemon-reload
systemctl start node_exporter.service
systemctl enable node_exporter #habilitar para realizar o start automatico 

```

precisamos adicionar o novo target no prometheus e realizar o restart do prometheus
```bash
- job_name: 'node_exporter'
	static_configs:
	  - targets: ['localhost:9100']
```

```bash
sudo systemctl restart prometheus
```	  

## adicionando um novo collector

os collectors  são os responsáveis por capturar as métricas do sistema operacional. Por padrão, o Node Exporter vem com um monte de coletores habilitados, mas você pode habilitar outros, caso queira.

[Collectors já habilitados](https://github.com/prometheus/node_exporter#enabled-by-default)

[Collectors que podem ser habilitados](https://github.com/prometheus/node_exporter#disabled-by-default)

### Adicionando um novo collector

 criar o arquivo /etc/node_exporter/node_exporter_options e o diretório /etc/node_exporter/ caso ele não exista:
```bash
sudo mkdir /etc/node_exporter
sudo vim /etc/node_exporter/node_exporter_options
   ```

Agora vamos adicionar a variável de ambiente OPTIONS no arquivo /etc/node_exporter/node_exporter_options:
```bash
OPTIONS="--collector.systemd"
  ```

Vamos ajustar as permissões do arquivo /etc/node_exporter/node_exporter_options:

```bash
sudo chown -R node_exporter:node_exporter /etc/node_exporter/
 ```

E no arquivo de configuração do serviço do Node Exporter para o SystemD, vamos adicionar a variável de ambiente OPTIONS e o arquivo vai ficar assim:

```bash
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
EnvironmentFile=/etc/node_exporter/node_exporter_options
ExecStart=/usr/local/bin/node_exporter $OPTIONS

[Install]
WantedBy=multi-user.target
``` 

Pronto, adicionamos o nosso novo arquivo que contém a variável de ambiente OPTIONS e agora vamos reiniciar o serviço do Node Exporter para que ele leia as novas configurações:
```bash
sudo systemctl daemon-reload
sudo systemctl restart node_exporter
```
[LINUXTIPS](https://github.com/badtuxx/DescomplicandoPrometheus)