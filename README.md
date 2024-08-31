# A File-System-Based Database with CLI using Golang

This project is a file-system-based database written in Go, featuring a Command-Line Interface (CLI) for managing collections and records. It provides basic CRUD (Create, Read, Update, Delete) operations using JSON files for data storage.



## Installation

1. Clone the repository:
   ```
   git clone https://github.com/sakarbaral/database-cli.git
   cd database-cli
   ```
2. Build the project

    ```go build -o db-cli```


## Usage

### Write a record
./db-cli write users "John Doe" "30" "555-1234" "Acme Corp" "New York"
### Read a record
./db-cli read users "John Doe"
### Read all records
./db-cli readall users
### Delete a record
./db-cli delete users "John Doe"
