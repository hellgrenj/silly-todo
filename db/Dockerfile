FROM golang:latest

# get migrate tool
RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

COPY start.sh start.sh
COPY migrations migrations 

RUN apt-get update && apt-get install -y dos2unix 
RUN dos2unix ./start.sh && apt-get --purge remove -y dos2unix && rm -rf /var/lib/apt/lists/*
RUN ["chmod", "+x", "./start.sh"]

ENTRYPOINT ["./start.sh"]
