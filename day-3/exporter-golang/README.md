# Criando nosso Segundo Exporter 

Criando o psegundo exporter em Golang para o prometheus.


### Primeiro vamos criar o nosso arquivo segundo-exprter.go

```go
package main

import ( // importando as bibliotecas necessárias
	"log"      // log
	"net/http" // http

	"github.com/pbnjay/memory"                                // biblioteca para pegar informações de memória
	"github.com/prometheus/client_golang/prometheus"          // biblioteca para criar o nosso exporter
	"github.com/prometheus/client_golang/prometheus/promhttp" // biblioteca criar o servidor web
)

func memoriaLivre() float64 { // função para pegar a memória livre
	memoria_livre := memory.FreeMemory() // pegando a memória livre através da função FreeMemory() da biblioteca memory
	return float64(memoria_livre)        // retornando o valor da memória livre
}

func totalMemory() float64 { // função para pegar a memória total
	memoria_total := memory.TotalMemory() // pegando a memória total através da função TotalMemory() da biblioteca memory
	return float64(memoria_total)         // retornando o valor da memória total
}

var ( // variáveis para definir as nossas métricas do tipo Gauge
	memoriaLivreBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{ // métrica para pegar a memória livre em bytes
		Name: "memoria_livre_bytes",                  // nome da métrica
		Help: "Quantidade de memória livre em bytes", // descrição da métrica
	})

	memoriaLivreMegabytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{ // métrica para pegar a memória livre em megabytes
		Name: "memoria_livre_megabytes",                  // nome da métrica
		Help: "Quantidade de memória livre em megabytes", // descrição da métrica
	})

	totalMemoryBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{ // métrica para pegar a memória total em bytes
		Name: "total_memoria_bytes",                  // nome da métrica
		Help: "Quantidade total de memória em bytes", // descrição da métrica
	})

	totalMemoryGigaBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{ // métrica para pegar a memória total em gigabytes
		Name: "total_memoria_gigabytes",                  // nome da métrica
		Help: "Quantidade total de memória em gigabytes", // descrição da métrica
	})
)

func init() { // função para registrar as métricas

	prometheus.MustRegister(memoriaLivreBytesGauge)     // registrando a métrica de memória livre em bytes
	prometheus.MustRegister(memoriaLivreMegabytesGauge) // registrando a métrica de memória livre em megabytes
	prometheus.MustRegister(totalMemoryBytesGauge)      // registrando a métrica de memória total em bytes
	prometheus.MustRegister(totalMemoryGigaBytesGauge)  // registrando a métrica de memória total em gigabytes
}

func main() { // função principal
	memoriaLivreBytesGauge.Set(memoriaLivre())                        // setando o valor da métrica de memória livre em bytes
	memoriaLivreMegabytesGauge.Set(memoriaLivre() / 1024 / 1024)      // setando o valor da métrica de memória livre em megabytes
	totalMemoryBytesGauge.Set(totalMemory())                          // setando o valor da métrica de memória total em bytes
	totalMemoryGigaBytesGauge.Set(totalMemory() / 1024 / 1024 / 1024) // setando o valor da métrica de memória total em gigabytes

	http.Handle("/metrics", promhttp.Handler()) // criando o servidor web para expor as métricas

	log.Fatal(http.ListenAndServe(":7788", nil)) // iniciando o servidor web na porta 7788
}

```

### Após a criação podemos ralizar o build do arquivo:

```bash
#Instalando as bibliotecas que utilizamos em nosso código
go mod init segundo-exporter
go mod tidy

#compilando o nosso código
go build segundo-exporter.go
```

## Criando a imagem com o Dockerfile e subindo o container
```bash
From golang:1.19.0-alpine3.16 AS buildando

WORKDIR /app
COPY . /app

RUN  go build segundo-exporter.go

FROM alpine:3.16

COPY --from=buildando /app/segundo-exporter /app/segundo-exporter
EXPOSE 7788
WORKDIR /app
CMD ["./segundo-exporter"]
```
Após termos criado o docker file vamos criar a imagem e subir o container
```bash
#Criando a imagem através do dockerfile
docker build -t segundo-exporter:1.0 .


#após criarmos a imagem vamos subir o container
docker run -d --name segundo-exporter -p 7788:7788 segundo-exporter:1.0
```  

### Feito isso não podemos esquecer de alterar no /etc/prometheus/prometheus.yml para adicionarmos o novo target 
```bash
 - job_name: 'segundo-exporter'
    static_configs:
      - targets: ['localhost:7788']

``` 
## Funções

### rate
A função rate representa a taxa de crescimento por segundo de uma determinada métrica como média, durante um intervalo de tempo.

```bash
rate(metrica)[5m]
```
 Onde metrica é a métrica que você deseja calcular a taxa de crescimento durante um intervalo de tempo de 5 minutos. Você pode utilizar a função rate para trabalhar com métricas do tipo gauge e counter.

Vamos para um exemplo real:
```bash
rate(prometheus_http_requests_total{job="prometheus",handler="/api/v1/query"}[5m])
```
### irate
A função irate representa a taxa de crescimento por segundo de uma determinada métrica, mas diferentemente da função rate, a função irate não faz a média dos valores, ela pega os dois últimos pontos e calcula a taxa de crescimento. Quando representado em um gráfico, é possível ver a diferença entre a função rate e a função irate, enquanto o gráfico com o rate é mais suave, o gráfico com o irate é mais "pontiagudo", você consegue ver quedas e subidas mais nítidas.

```bash
irate(metrica)[5m]
```
Onde metrica é a métrica que você deseja calcular a taxa de crescimento, considerando somente os dois últimos pontos, durante um intervalo de tempo de 5 minutos.

Vamos para um exemplo real:

```bash
irate(prometheus_http_requests_total{job="prometheus",handler="/api/v1/query"}[5m])
```

### delta
A função delta representa a diferença entre o valor atual e o valor anterior de uma métrica. Quando estamos falando de delta estamos falando por exemplo do consumo de um disco. Vamos imaginar que eu queira saber o quando eu usei de disco em um determinado intervalo de tempo, eu posso utilizar a função delta para calcular a diferença entre o valor atual e o valor anterior.

```bash
delta(metrica[5m])
 ```

Onde metrica é a métrica que você deseja calcular a diferença entre o valor atual e o valor anterior, durante um intervalo de tempo de 5 minutos.

 

Vamos para um exemplo real:

```bash
delta(prometheus_http_response_size_bytes_count{job="prometheus",handler="/api/v1/query"}[5m])
```
Agora estou calculando a diferença entre o valor atual e o valor anterior da métrica prometheus_http_response_size_bytes_count, filtrando por job e handler e durante um intervalo de tempo de 5 minutos. Nesse caso eu quero saber o quanto de bytes eu estou consumindo nas queries que estão sendo feitas no Prometheus.

  


### increase

Da mesma forma que a função delta, a função increase representa a diferença entre o primeiro e último valor durante um intervalo de tempo, porém a diferença é que a função increase considera que o valor é um contador, ou seja, o valor é incrementado a cada vez que a métrica é atualizada. Ela começa com o valor 0 e vai somando o valor da métrica a cada atualização. Você já pode imaginar qual o tipo de métrica que ela trabalha, certo? Qual? Counter!


```bash
increase(metrica[5m])
 ```

Onde metrica é a métrica que você deseja calcular a diferença entre o primeiro e último valor durante um intervalo de tempo de 5 minutos.

 

Vamos para um exemplo real:



```bash
increase(prometheus_http_requests_total{job="prometheus",handler="/api/v1/query"}[5m])
 ```

Aqui estou calculando a diferença entre o primeiro e último valor da métrica prometheus_http_requests_total, filtrando por job e handler e durante um intervalo de tempo de 5 minutos.			

### sum
A função sum representa a soma de todos os valores de uma métrica. Você pode utilizar a função sum nos tipos de dados counter, gauge, histogram e summary. Um exemplo de uso da função sum é quando você quer saber o quanto de memória está sendo utilizada por todos os seus containers, 
ou o quanto de memória está sendo utilizada por todos os seus pods.

```bash
sum(metrica)
 ```

Onde metrica é a métrica que você deseja somar.

 

Vamos para um exemplo real:

```bash
sum(go_memstats_alloc_bytes{job="prometheus"})
 ```
 ### count
 
 Outra função bem utilizada é função count representa o contador de uma métrica. Você pode utilizar a função count nos tipos de dados counter, gauge, histogram e summary.
 Um exemplo de uso da função count é quando você quer saber quantos containers estão rodando em um determinado momento ou quantos de seus pods estão em execução.

```bash
count(metrica)
 ```

Onde metrica é a métrica que você deseja contar.

 

Vamos para um exemplo real:

```bash
count(prometheus_http_requests_total)
 ```
 
 ### avg
A função avg representa o valor médio de uma métrica. Você pode utilizar a função avg nos tipos de dados counter, gauge, histogram e summary. Essa é uma das funções mais utilizadas, pois é muito comum você querer saber o valor médio de uma métrica, por exemplo, o valor médio de memória utilizada por um container.

```bash
avg(metrica)
  ```
Onde metrica é a métrica que você deseja calcular a média.

###
min
A função min representa o valor mínimo de uma métrica. Você pode utilizar a função min nos tipos de dados counter, gauge, histogram e summary. Um exemplo de uso da função min é quando você quer saber qual o menor valor de memória utilizada por um container.

```bash
min(metrica)
```
Onde metrica é a métrica que você deseja calcular o mínimo.

### max
A função max representa o valor máximo de uma métrica. Um exemplo de uso da função max é quando você quer saber qual o maior valor de memória pelos nodes de um cluster Kubernetes.

```bash
max(metrica)
   ```
 Onde metrica é a métrica que você deseja calcular o máximo.


### avg_over_time
A função avg_over_time representa a média de uma métrica durante um intervalo de tempo. Normalmente utilizada para calcular a média de uma métrica durante um intervalo de tempo, como por exemplo, a média de requisições por segundo durante um intervalo de tempo ou ainda as pessoas que estão no espaço durante o último ano. :D

```bash
avg_over_time(metrica[5m])
 ```
   
Onde metrica é a métrica que você deseja calcular a média durante um intervalo de tempo de 5 minutos.

 

Vamos para um exemplo real:

```bash
avg_over_time(prometheus_http_requests_total{handler="/api/v1/query"}[5m])
  ```

Agora estou calculando a média da métrica prometheus_http_requests_total, filtrando por handler e durante um intervalo de tempo de 5 minutos.

  


### sum_over_time
Também temos a função sum_over_time, que representa a soma de uma métrica durante um intervalo de tempo. Vimos a avg_over_time que representa a média, a sum_over_time representa a soma dos valores durante um intervalo de tempo. Imagina calcular a soma de uma métrica durante um intervalo de tempo, como por exemplo, a soma de requisições por segundo durante um intervalo de tempo ou ainda a soma de pessoas que estão no espaço durante o último ano.

```bash
sum_over_time(metrica[5m])
 ```

Onde metrica é a métrica que você deseja calcular a soma durante um intervalo de tempo de 5 minutos.

 

Vamos para um exemplo real:

```bash
sum_over_time(prometheus_http_requests_total{handler="/api/v1/query"}[5m])
 ```

Agora estou calculando a soma da métrica prometheus_http_requests_total, filtrando por handler e durante um intervalo de tempo de 5 minutos.

  


### max_over_time
A função max_over_time representa o valor máximo de uma métrica durante um intervalo de tempo.

```bash
max_over_time(metrica[5m])
   ```
Onde metrica é a métrica que você deseja calcular o valor máximo durante um intervalo de tempo de 5 minutos.

 
Vamos para um exemplo real:

```bash
max_over_time(prometheus_http_requests_total{handler="/api/v1/query"}[5m])
  ```

Agora estamos buscando o valor máximo da métrica prometheus_http_requests_total, filtrando por handler e durante um intervalo de tempo de 5 minutos.

  


min_over_time
A função min_over_time representa o valor mínimo de uma métrica durante um intervalo de tempo.

```bash
min_over_time(metrica[5m])
 ```

Onde metrica é a métrica que você deseja calcular o valor mínimo durante um intervalo de tempo de 5 minutos.

 

Vamos para um exemplo real:

```bash
min_over_time(prometheus_http_requests_total{handler="/api/v1/query"}[5m])
  ``` 

Agora estamos buscando o valor mínimo da métrica prometheus_http_requests_total, filtrando por handler e durante um intervalo de tempo de 5 minutos.

  


### stddev_over_time
A função stddev_over_time representa o desvio padrão, que são os valores que estão mais distantes da média, de uma métrica durante um intervalo de tempo. Um bom exemplo seria para o calculo de desvio padrão para saber se houve alguma anomalia no consumo de disco, por exemplo.

```bash
stddev_over_time(metrica[5m])
  ```
Onde metrica é a métrica que você deseja calcular o desvio padrão durante um intervalo de tempo de 5 minutos.

 

Vamos para um exemplo real:

```bash
stddev_over_time(prometheus_http_requests_total{handler="/api/v1/query"}[10m])
 ```

Agora estamos buscando os desvios padrões da métrica prometheus_http_requests_total, filtrando por handler e durante um intervalo de tempo de 10 minutos. Vale a pena verificar o gráfico, pois facilita a visualização dos valores.
  

  
### by
A sensacional e super utilizada função by é utilizada para agrupar métricas. Com ela é possível agrupar métricas por labels, por exemplo, 
se eu quiser agrupar todas as métricas que possuem o label job eu posso utilizar a função by da seguinte forma:

``bash
sum(metrica) by (job)
``` 

Onde metrica é a métrica que você deseja agrupar e job é o label que você deseja agrupar.

 

Vamos para um exemplo real:

``bash
sum(prometheus_http_requests_total) by (code)
```  

Agora estamos somando a métrica prometheus_http_requests_total e agrupando por code, assim sabemos quantas requisições foram feitas por código de resposta.

  

### without
A função without é utilizada para remover labels de uma métrica. Você pode utilizar a função without nos tipos de dados counter, gauge, histogram e summary e frequentemente usado em conjunto com a função sum.

Por exemplo, se eu quiser remover o label `job` de uma métrica, eu posso utilizar a função `without` da seguinte forma:


```bash
sum(metrica) without (job)
 ```

Onde metrica é a métrica que você deseja remover o label job.

 
Vamos para um exemplo real:

```bash
sum(prometheus_http_requests_total) without (handler)
 ``` 
 
Agora estamos somando a métrica prometheus_http_requests_total e removendo o label handler, 
assim sabemos quantas requisições foram feitas por código de resposta, sem saber qual handler foi utilizado para ter uma visão mais geral e focado no código de resposta.

### histogram_quantile e quantile
As funções histogram_quantile e quantile são muito parecidas, porém a histogram_quantile é utilizada para calcular o percentil de uma métrica do tipo histogram e a quantile é utilizada para calcular o percentil de uma métrica do tipo summary. Basicamente utilizamos esses funções para saber qual é o valor de uma métrica em um determinado percentil.

```bash
quantile(0.95, metrica)
 ```  
 
Onde metrica é a métrica do tipo histogram que você deseja calcular o percentil e 0.95 é o percentil que você deseja calcular.

 

Vamos para um exemplo real:

```bash
quantile(0.95, prometheus_http_request_duration_seconds_bucket)
```  

Agora estamos calculando o percentil de 95% da métrica prometheus_http_request_duration_seconds_bucket, assim sabemos qual é o tempo de resposta de 95% das requisições.





[LINUXTIPS](https://github.com/badtuxx/DescomplicandoPrometheus)