# meeting-scheduler-api
A RESTful API built with Golang to help distributed teams find the best meeting time slots based on participant availability. Supports event creation, slot management, and intelligent time slot recommendations. Designed for scalability, clean code, and cloud-native deployment.

## Features
- **Event Management**: Create, update, and delete events with multiple participants.
- **Availability Management**: Participants can set their availability for specific time slots.
- **Intelligent Slot Recommendations**: Automatically suggest optimal meeting times based on participant availability.

- **Scalability**: Designed to handle a large number of participants and events efficiently.
- **Cloud-Native**: Built with cloud-native principles for easy deployment and scaling.
- **Clean Code**: Follows best practices for maintainable and readable code.

## Technologies
- **Golang**: Primary language for building the API.
- **Gorilla/mux**: Web framework for building RESTful APIs.
- **go-migrate**: Database migration tool for managing schema changes.
- **MYSQL**: Relational database for storing events and availability data.
- **Docker**: Containerization for easy deployment and scaling.

## Getting Started
### Prerequisites
- Go 1.18 or later
- Docker (for containerization)
- Docker Compose (for managing multi-container applications)

### Installation
Run the following commands to set up the project:

```bash
docker-compose up --build
```
Note: database migrations will be applied automatically using `go-migrate`.
If you see error related to `mysql_native_password` plugin, then uncomment following line in `docker-compose.yml` file and restart the containers:

```yaml
 #command: --default-authentication-plugin=mysql_native_password
 ```


### Running the API
After the Docker containers are up, you can access the API at `http://localhost:8001`.

### API Documentation
You can find the API documentation in openapi-swagger.yml file. Use Swagger UI or Postman to explore the endpoints.

### Testing
Run the following command to execute the tests:

```bash
go test ./...
```



