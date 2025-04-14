# sqlc-example-api

This repository provides starter code for creating a API using sqlc and the Gin web framework in Go. This is part of the Relational Database course as part of the Iknite Space training.


Project Structure

* `api/`: Contains API route definitions and handler functions. You will need to edit this file, if you want to add or modify the apit endpoints.
* `cmd/api/`: Houses the main application entry point. You shouldn't need to edit any files in this directory.
* `db/migrations`: Contains sql files that create/update the database schema (tables and columns) used by the api. If you need to update the database schema make your changes here.
* `db/query`: This folder contains SQL query files. These files define the database queries used by the application, which are processed by sqlc to generate type-safe Go code for interacting with the database.
`db/repo/`: This directory contains repository code that acts as an abstraction layer between the database and the application logic. It provides functions to interact with the database using the generated sqlc code. You shouldn't need to edit any files in this directory.
* `campay_api/Payment/` : Contains campay api endpoint logics for payment request and payment status. You shouldn't need to edit any files in this directory.

# Getting Started
This api is a streamlined solution designed to handle customer and order management with seamless integration for payment processing. It utilizes PostgreSQL as the database, ensuring robust and scalable data storage. The app currently provides two key endpoints:

    Create Customer: This endpoint allows the creation of a new customer record in the database, capturing essential details required for customer management.

    Create Order: This endpoint enables the creation of an order associated with a specific customer. The order details, including the total price and other related information, are securely stored in the database.

In addition to these endpoints, the app integrates with the Campay API for payment processing. After an order is created, the app initiates a payment request through Campay. Once the payment status is received, it updates the corresponding order's status in the database, keeping records consistent and up-to-date.

This api serves as a versatile foundation for e-commerce platforms, payment systems, or order management solutions, offering flexibility and a clear roadmap for future expansion. 



* Install Dependencies:
    * Ensure you have Go installed. Then, install the required Go modules:
    `go mod tidy`
* Generate Code with sqlc:
    * If you make any changes to `db/query` or `db/migrations`, then you need to re-generate the go code in `db/repo`
            `go generate ./...`
    * This command will generate Go code based on the SQL queries defined in the db/queries/ directory.
* Run the Application:
    `go run cmd/api/...`

    The server will start, and it will print the available endpoints.
