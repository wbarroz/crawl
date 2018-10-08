FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go get -u github.com/gocolly/colly/...
RUN go get -u github.com/go-sql-driver/mysql
RUN go build -o main . 
CMD ["/app/main"]
