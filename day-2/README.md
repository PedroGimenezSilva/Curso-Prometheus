# Criando nosso primeiro Exporter 

Criando o primeiro exporter em python para o prometheus.

### Primeiro vamos criar o nosso arquivo exporter.py

```python
mport requests # Importa o módulo requests para fazer requisições HTTP
import json # Importa o módulo json para converter o resultado em JSON
import time # Importa o módulo time para fazer o sleep
from prometheus_client import start_http_server, Gauge # Importa o módulo Gauge do Prometheus para criar a nossa métrica e o módulo start_http_server para iniciar o servidor

url_numero_pessoas = 'http://api.open-notify.org/astros.json' # URL para pegar o número de astronautas
url_local_ISS = 'http://api.open-notify.org/iss-now.json' # URL para pegar a localização do ISS

def pega_local_ISS(): # Função para pegar a localização da ISS
    try:
        """
        Pegar o local da estação espacial internacional
        """
        response = requests.get(url_local_ISS) # Faz a requisição para a URL
        data = response.json() # Converte o resultado em JSON
        return data['iss_position'] # Retorna o resultado da requisição
    except Exception as e: # Caso ocorra algum erro
        print("Não foi possível acessar a url!") # Imprime uma mensagem de erro
        raise e # Lança a exceção

def pega_numero_astronautas(): # Função para pegar o número de astronautas
    try: # Tenta fazer a requisição HTTP
        """
        Pegar o número de astronautas no espaço 
        """
        response = requests.get(url) # Faz a requisição HTTP
        data = response.json() # Converte o resultado em JSON
        return data['number'] # Retorna o número de astronautas
    except Exception as e: # Se der algum erro
        print("Não foi possível acessar a url!") # Imprime que não foi possível acessar a url
        raise e # Lança a exceção

def atualiza_metricas(): # Função para atualizar as métricas
    try:
        """
        Atualiza as métricas com o número de astronautas e local da estação espacial internacional
        """
        numero_pessoas = Gauge('numero_de_astronautas', 'Número de astronautas no espaço') # Cria a métrica para o número de astronautas
        longitude = Gauge('longitude_ISS', 'Longitude da Estação Espacial Internacional') # Cria a métrica para a longitude da estação espacial internacional
        latitude = Gauge('latitude_ISS', 'Latitude da Estação Espacial Internacional') # Cria a métrica para a latitude da estação espacial internacional

        while True: # Enquanto True
            numero_pessoas.set(pega_numero_astronautas()) # Atualiza a métrica com o número de astronautas
            longitude.set(pega_local_ISS()['longitude']) # Atualiza a métrica com a longitude da estação espacial internacional
            latitude.set(pega_local_ISS()['latitude']) # Atualiza a métrica com a latitude da estação espacial internacional
            time.sleep(10) # Faz o sleep de 10 segundos
            print("O número atual de astronautas no espaço é: %s" % pega_numero_astronautas()) # Imprime o número atual de astronautas no espaço
            print("A longitude atual da Estação Espacial Internacional é: %s" % pega_local_ISS()['longitude']) # Imprime a longitude atual da estação espacial internacional
            print("A latitude atual da Estação Espacial Internacional é: %s" % pega_local_ISS()['latitude']) # Imprime a latitude atual da estação espacial internacional
    except Exception as e: # Se der algum erro
        print("Problemas para atualizar as métricas! \n\n====> %s \n" % e) # Imprime que ocorreu um problema para atualizar as métricas
        raise e # Lança a exceção
        
def inicia_exporter(): # Função para iniciar o exporter
    try:
        """
        Iniciar o exporter
        """
        start_http_server(8899) # Inicia o servidor do Prometheus na porta 8899
        return True # Retorna True
    except Exception as e: # Se der algum erro
        print("O Servidor não pode ser iniciado!") # Imprime que não foi possível iniciar o servidor
        raise e # Lança a exceção

def main(): # Função principal
    try:
        inicia_exporter() # Inicia o exporter
        print('Exporter Iniciado') # Imprime que o exporter foi iniciado
        atualiza_metricas() # Atualiza as métricas
    except Exception as e: # Se der algum erro
        print('\nExporter Falhou e Foi Finalizado! \n\n======> %s\n' % e) # Imprime que o exporter falhou e foi finalizado
        exit(1) # Finaliza o programa com erro


if __name__ == '__main__': # Se o programa for executado diretamente
    main() # Executa o main
    exit(0) # Finaliza o programa
```

### Após a criação podemos testar localmente para ver se ele está funcionando:

```bash
python3 exporter.py
# ele deverá aparecer o seguinte output
HTTP Server iniciado
O número de Astronautas no espaço nesse momento é: 10
A longitude atual da ISS é: 82.4048
A latitude atual da ISS é: -44.4759
```
### Instalando o Docker se necssário 
```bash
curl -fsSL https://get.docker.com | bash
```
### criando o nosso Dockerfile
```bash
# Vamos utilizar a imagem slim do Python
FROM python:3.8-slim

# Adicionando algumas labels para identificar a imagem
LABEL maintainer Pedro <pedro@email.com>
LABEL description "Dockerfile para criar o nosso primeiro exporter para o Prometheus"

# Adicionando o exporter.py para a nossa imagem
COPY . .

# Instalando as bibliotecas necessárias para o exporter
# através do `requirements.txt`.
RUN pip3 install -r requirements.txt

# Executando o exporter
CMD python3 exporter.py
```

### Após criado o Dockerfile vamos criar a nossa imagem
```bash
docker build -t primeiro_exporter:0.1 .
```

### Agora precisamos criar um container com nosso exporter
```bash
docker run -p 8899:8899 --name primeiro-exporter -d primeiro-exporter:0.1
```
-p 8899:8899 ---> Significa que a porta 8899 do host seja mapeada para a 8899 do container \
--name ---> Define o nome do container \
-d ---> serve para executar o container em segundo plano \
primeiro-exporter:0.1 ---> Define o nome da sua imagem e sua tag

### Validando se o exporter está funcionando
```bash
curl -s http://localhost:8899/metrics

```

### Vendo os targets que estão rodando no Prometheus em formato json com o jq

```bash
curl -s http://localhost:9090/api/v1/targets | jq .

```
### Adicionando um novo target no Prometheus
será necessário abrir o arquivo prometheus.yml que está localizado em /etc/prometheus/ \
Após aberto será necessário apenas adicionar mais um target:
```bash
 - job_name: "Primeiro Exporter" # Nome do job que vai coletar as métricas do primeiro exporter.
    static_configs:
      - targets: ["localhost:8899"] # Endereço do alvo monitorado, ou seja, o nosso primeiro exporter.
```
feito isso é so realizar o restart do prometheus para ele pegar as modificações
```bash
sudo systemctl restart prometheus

```
Caso queira validar os targets executar:

```bash
curl -s http://localhost:9090/api/v1/targets | jq . 
```

### Tipos de dados
gauge: Medidor
O tipo de dado gauge é o tipo de dado utilizado para criar métricas que podem ter seus valores alterados para cima ou para baixo, por exemplo, a ultilização de memória ou cpu. Se quiser trazer para exemplos da vida real, podemos falar que aquelas filas que você odeia é o tipo de dado gauge, ou então a temperatura da sua cidade, ela pode ser alterada para cima ou para baixo, ou seja, é um medidor, é um gauge! :D Um exemplo de métrica do tipo gauge é a métrica memory_usage, que é uma métrica que mostra a utilização de memória.
```bash
memory_usage{instance="localhost:8899",job="Primeiro Exporter"}
 ```


counter: Contador
O tipo de dado counter é o tipo de dado utilizado que vai ser incrementado no decorrer do tempo, por exemplo, quando eu quero contar os erros em uma aplicação no decorrer da última hora. O valor atual do counter quase nunca é importante, pois o que queremos dele são os valores durante uma janela de tempo, por exemplo, quantas vezes a minha aplicação falhou durante o final de semana. Normalmente as métricas counter possuem o sufixo _total para indicar que é o total de valores que foram contados, por exemplo:
```bash
requests_total{instance="localhost:8899",job="Primeiro Exporter"}
  ```


histogram: Histograma
O tipo de dado histogram é o tipo de dado que te permite especificar o seu valor através de buckets predefinidos, por exemplo, o tempo de execução de uma aplicação. Com o histogram eu consigo contar todas as requisições que minha aplicação respondeu entre 0 e 0,5 segundos, ou então as requisições que tiveram respostas entre 1,0 e 2,5 e assim por diante. Por padrão, os buckets predefinidos são até no máximo 10 segundos, se você quiser mais, você pode criar seus próprios buckets personalizados. Um coisa super importante, o Prometheus irá contar cada item em cada bucket, e também a soma dos valores. Uma métrica do tipo histogram inclui alguns itens importantes são adicionados ao final do nome da métrica para indicar o tipo de dado e o tamanho do bucket, por exemplo:
```bash
requests_duration_seconds_bucket{le="0.5"}
 ```
Onde le é o valor limite do bucket, o valor 0.5 indica que o valor do bucket é até 0,5 segundos, ou seja, aqui nesse bucket poderiam estar os requisições que tiveram respostas entre 0 e 0,5 segundos.\
O _bucket é um sufixo que indica que o valor é um bucket. \
Ainda temos alguns sufixos que são importantes e que podem ser úteis para nós: \
_count: Contador \
O sufixo _count indica que o valor é um contador, ou seja, o valor é incrementado a cada vez que a métrica é atualizada. \
_sum: Soma \
O sufixo _sum indica que o valor é uma soma, ou seja, o valor é somado a cada vez que a métrica é atualizada.

O ponto alto do histogram é a excelente flexibilidades, pois percentuais e as janelas de tempos podem definidas durante a criação das queries, o ponto negativo é a precisão é um pouco inferior quando comparado com o summary.
 


summary: Resumo
O tipo de dado summary é bem parecido com o histogram, com a diferença que os buckets, aqui chamados de quantiles, são definidos por um valor entre 0 e 1, ou seja, o valor do bucket é o valor que está entre os quantiles.
Da mesma forma como no histogram, podemos criar métricas do tipo summary com alguns itens importantes adicionados ao final do nome da métrica, por exemplo:
```bash
requests_duration_seconds_sum{instance="localhost:8899",job="Primeiro Exporter"}
 ```
Utilizamos o sufixo _sum indica que o valor é uma soma, ou seja, o valor é somado a cada vez que a métrica é atualizada e o sufixo _count para indicar que o valor é um contador, ou seja, o valor é incrementado a cada vez que a métrica é atualizada.
O ponto alto do summary é a excelente precisão e o ponto baixo é a baixa flexibilidades, pois percentuais e as janelas de tempos precisam ser definidos durante a criação da métrica e não é possível agregar métricas do tipo summary com outras métricas do tipo summary durante a criação das queries.

[LINUXTIPS](https://github.com/badtuxx/DescomplicandoPrometheus)