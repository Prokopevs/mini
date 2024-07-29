# <center>Messages api</center>
---

### Description
- User create a message and send it to the server. The server processes the message and places it in the messages table with `status idle`. 
- To ensure delivery to Kafka, the `Outbox pattern` was applied. The Outbox table did not need to be created because the message entity already has the outbox property. 
- Next, for visibility, a worker is launched every `5 seconds` to observe how message statuses are updated. The worker reads up to `3 messages` with status='idle' per limit. Then, the message IDs are sent to Kafka and marked as sent `(status='send')`. 
- Kafka consumers read the message IDs and then update their status to `received`.

### Endpoints
`http://localhost:5555/api/v1/getMessages` (Get) Returns a list of all messages, or an empty array if there are none.  

`http://localhost:5555/api/v1/create` (Post) Creates a new message. Request body example: { "message": "msg-1" }.  

`http://localhost:5555/api/v1/swagger/index.html` (Get) Swagger documentation.  

### How to start
1. _cd mini; make build_image_
2. _cd server_ 
3. _sudo PG_USER=postgres PG_PASSWORD=postgres KAFKA_BROKER=broker:9092 KAFKA_TOPIC=topicname KAFKA_GROUP_ID=group BUFFER_SIZE=3 docker compose up_
4. Error minidb doesn't exist will appear. Stop server using _ctrl+c_
5. Create database minidb in docker postgres container
    - _docker start miniPostgresCont_
    - _docker exec -it miniPostgresCont psql -U postgres_
    - _create database minidb;_
    - _\q_
6. Start migration
    - _docker inspect miniPostgresCont_ and copy field <IPAddress> (example 172.21.0.2)
    - _migrate -path db/migrations -database "postgres://postgres:postgres@<IPAddress>:5432/minidb?sslmode=disable" -verbose up_    
7. Start server with command number 3
8. Open new terminal and create topic
    - _docker-compose exec broker kafka-topics --create --topic topicname --partitions 1 --replication-factor 1 --bootstrap-server broker:9092_
    - _docker exec broker kafka-topics --bootstrap-server broker:9092 --describe todos_ check if the topic was created
9. restart server _ctrl+c + command number 3_

### Logs
message-api-1     | {"level":"info","ts":1722179582.900519,"caller":"db/message.go:20","msg":"Successfully added to db message with id","id":128}
message-api-1     | {"level":"info","ts":1722179589.2348382,"caller":"worker/sender.go:73","msg":"Successfully sended messages to kafka","ids":[128]}
message-api-1     | {"level":"info","ts":1722179589.2356353,"caller":"worker/sender.go:82","msg":"Changed statuses to send in db for","ids":[128]}
message-api-1     | {"level":"info","ts":1722179589.2397115,"caller":"worker/listener.go:58","msg":"Successfully readed from kafka","id":128}
message-api-1     | {"level":"info","ts":1722179589.2397833,"caller":"worker/listener.go:59","msg":"Changed statuses to received in db for","id":128}

### Stack:
- Go
- Gin
- Docker
- PostgreSQL
- Kafka
- Swagger
---

