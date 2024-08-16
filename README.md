# Go Serverless CRUD App

This project is a serverless CRUD application built using AWS Lambda and DynamoDB. It provides a simple API for managing user data using API Gateway , DynamoDB , Lambda Comlpete serveless stack.

This project is a serverless CRUD application built using AWS Lambda and DynamoDB. It provides a simple API for managing user data.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Middleware](#middleware)
- [Environment Variables](#environment-variables)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Reaper1994/go-serverless-app.git
   cd go-serverless-app
   ```

2. Set up your Go environment and install dependencies.

## Start the application
go run cmd/main.go

## API Endpoints

The application exposes the following API endpoints:

- **GET https://localhost:8000** : Retrieve a list of users.
- **POST https://localhost:8000**: Create a new user.
- **PUT https://localhost:8000/{id}**: Update an existing user.
- **DELETE https://localhost:8000/{id}**: Delete a user.

## Middleware

The application uses middleware for logging and monitoring. The `TreblleMiddleware` is included to track API requests and responses.

## Environment Variables

The following environment variables are required:

- `TREBLLE_API_KEY`: Your Treblle API key for monitoring.
- `TREBLLE_PROJECT_ID`: Your Treblle project ID.
- `AWS_REGION`: The AWS region where your DynamoDB is hosted. // auto setup by AWS


## Running the Application

1. Install dependencies:
   ```bash
   go get ./...
   ```

2. Run the application:
   ```bash
   go run cmd/main.go
   ```

3. Test the endpoints using a tool like Postman or curl.
## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

