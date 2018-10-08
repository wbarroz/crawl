# crawl

### How to put it running

Simply execute:

>go run crawl.go

if you only need to see if it works, or

>go build crawl.go

to build a binary.


### If you need the running mysql instance
Get into my-mysql folder:

>cd ~/my-mysql/
>docker build -t my-mysql .
>docker run -d -p 3306:3306 --name my-mysql -e MYSQL_ROOT_PASSWORD=supersecret my-mysql

after that you should have a playful mysql server working. After that, you can:

>docker exec -it my-sql bash

and peep into the server.

Running the crawl code should populate the stocks table in redventures database with the top ten stock positions.

### The container
Yes, the app is now containeirized!
In this folder just fire it up:

>docker build -t crawl .

to build the container, and

>docker run --net=host --name=crawler -d crawl:latest

to run it, it's a one-off run. You could check the db for the updated stocks table.

### Miscellaneous
In my-mysql folders there's an alternate my.cnf file, just in case you need to connect
to mysql from outside without much hassle.

### Known bugs
At the time of this release, the crawler couldn't fetch the very first stock option, named AALR3.
Partly because the access instability of the site, I couldn't dig further to solve this, yet.
