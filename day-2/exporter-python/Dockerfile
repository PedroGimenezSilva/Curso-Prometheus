From python:3.8-slim

LABEL maintainer Pedro Gimenez Miranda Silva <pedrogimenez7@hotmail.com>
LABEL description "Dockerfile   para criar as imagens do container"

WORKDIR /app
COPY . /app
RUN pip3 install -r requirements.txt

CMD python3 exporter.py
