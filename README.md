# crawl
How to put it running
-Simply execute:

go run crawl.go

if you only need to see if it works, or

go build crawl.go

to build a binary.


-If you need the running mysql instance:
Get into my-mysql folder:

cd ~/my-mysql/
docker build -t my-mysql .
docker run -d -p 3306:3306 --name my-mysql -e MYSQL_ROOT_PASSWORD=supersecret my-mysql

after that you should have a playful mysql server working. After that, you can:

docker exec -it my-sql bash

and peep into the server.

Running the crawl code should populate the stocks table in redventures database with the top ten stock positions.

