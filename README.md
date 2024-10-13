# Transaction Routine Application 

This repository contains the source code for transaction routine application built using Golang. The system is responsible for creating an account. Once the account is created, users can create the transactions.

## Prerequisites

Before running the Transaction Routine Application, make sure you have the following prerequisites installed on your system:

- Go programming language (go1.23.1)
- PostgreSQL(14.8)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/ankitchahal20/transaction-routine.git
   ```

2. Navigate to the project directory:

   ```bash
   cd transaction-routine
   ```

3. Install the required dependencies:

   ```bash
   go mod tidy
   ```

4. Defaults.toml

    Add the values to defaults.toml and execute `go run ./cmd/main.go` from the root directory.
    
5. Swagger docs are accesible at 
    ```
    http://localhost:8080/docs/ui/#/
    ```

## APIs
There are 4 API's which this repo currently supports.

Create Account API
```
curl -i -k -X POST \
   http://localhost:8080/v1/accounts \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
        "document_number": "1000"
    }'
```

Get Account By Id API

```
curl -i -k -X POST \
  http://localhost:8080/v1/accounts/:account_id \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json"
```

Create Transaction API

```
curl -i -k -X POST \
  http://localhost:8080/v1/transactions \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
	    "account_id": "7f6b2bc0-62cc-45b3-aa86-37ab635f4c9f",
	    "operation_type_id": 2,
	    "amount": -1223.45
    }'
```

Get Transaction By Id API

```
curl -i -k -X GET \
  http://localhost:8080/v1/transactions/:transaction_id \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json"
```

## Project Structure

The project follows a standard Go project structure:

- `config/`: Configuration file for the application.
- `internal/`: Contains the internal packages and modules of the application.
  - `config/`: Global configuration which can be used anywhere in the application.
  - `constants/`: Contains constant values used throughout the application.
  - `db/`: Contains the database package for interacting with PostgreSQL.
  - `middleware`: Contains the logic to validate the incoming request
  - `models/`: Contains the data models used in the application.
  - `transactionroutineerror`: Defines the errors in the application
  - `service/`: Contains the business logic and services of the application.
  - `server/`: Contains the server logic of the application.
  - `utils/`: Contains utility functions and helpers.
- `cmd/`:  Contains command you want to build.
    - `main.go`: Main entry point of the application.
- `README.md`: README.md contains the description for the Transaction Routine Application.

## Contributing

Contributions to the Transaction Routine Application are welcome. If you find any issues or have suggestions for improvement, feel free to open an issue or submit a pull request.

## License

The Transaction Routine Application is open-source and released under the [MIT License](LICENSE).

## Contact

For any inquiries or questions, please contact:

- Ankit Chahal
- ankitchahal20@gmail.com

Feel free to reach out with any feedback or concerns.
