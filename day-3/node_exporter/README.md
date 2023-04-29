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



[LINUXTIPS](https://github.com/badtuxx/DescomplicandoPrometheus)