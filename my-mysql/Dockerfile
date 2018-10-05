#Derived from official mysql image
FROM mysql

#Adding database
ENV MYSQL_DATABASE redventures

#Container startup script
COPY ./sql-scripts/ /docker-entrypoint-initdb.d/

#Is this needed??
#ENTRYPOINT ["/entrypoint.sh"] #doesn't work
