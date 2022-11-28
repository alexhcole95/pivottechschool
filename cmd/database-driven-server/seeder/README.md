
# Database-Driven Product Server

## Objective

Refactor the product server to use a database instead of a JSON file.

## Learning Goals

1. Using SQLite to create databases and tables
2. Learn how to seed databases
3. Learn how to execute DDL and DML SQL statements
4. Using the Go database/sql package and SQLite driver to interact with databases
5. Refactoring and extending existing code

## Part 1: Seeding a Database

You assignment will comprise two executables, one for seeding a products database and the other (which you should already have) to serve up the REST API that uses the seeded database. The following requirements are for the seeder.

### Acceptance Criteria

1. Have a `cmd/seeder` executable in your repository (hereafter referred to as **seeder**)
2. When launched, seeder looks for and deletes a `products.db` file so as to start with a fresh database every time it is executed. Hint: Use the `os` package to do this.
3. Seeder uses `database/sql` along with a SQLite driver of your choice ([github.com/mattn/go-sqlite3](http://github.com/mattn/go-sqlite3) is recommended) to create a new database (`products.db`).
4. Seeder executes the appropriate SQL statement to create a table called `products` in the database. The columns of the products table must at least include an `id` (preferably an autoincrementing integer) that serves as the table’s primary key, a `name` of type `text`, and a `price` of type `real` (to support floats). You may have additional columns if you wish.
5. Seeder loads your existing products using the `products.json` file from your previous product server assignment.
6. Seeder iterates through your list of products and inserts them into the `products` table according to the following:
    1. Each products must have a unique ID
    2. A transaction is used to encapsulate all of the product inserts
    3. A prepared statement (with `?` symbols to accept `name`, `price`, etc) is used for the inserts inside the transaction
    4. When executed, the prepared statement is given the parameters that replace the `?` symbols in the SQL statement with the appropriate values of each product that you’re iterating over.
7. Seeder confirms all the data was imported by performing a SELECT query with a LIMIT of 5 and outputs the resulting products.
    1. You must populate a `Product` struct for each row of the query result and append to a list
    2. You must then iterate through your list of products, outputting each product’s `id`, `name`, and `price`
8. Your executable lists the `id`, `name`, and `price` of the products to STDOUT (your terminal). **Not a dump**.
9. Your work is done on a separate branch (e.g. `add-product-database`) and you provide a link to a GitHub pull request (PR) against your repository’s `main` branch with your solution for evaluation.
10. Your PR’s comment includes a copy of the output of your program.

## Part 2: Refactor HTTP Handlers

This part of your assignment calls for refactoring your existing product server’s HTTP handlers so that they use the products database you seeded in Part 1 above.

### Acceptance Criteria

1. Your server is initialized with a connection to the `products` database using the `database/sql` package and a suitable SQLite driver.
2. Your server verifies its connection to the database with a Ping (see [https://pkg.go.dev/database/sql@go1.19.3#DB.Ping](https://pkg.go.dev/database/sql@go1.19.3#DB.Ping)) and stops initialization cleanly (i.e. no `panic`) if the connection is invalid.
3. Your server supports the following operations
    1. `GET /products?limit=<numeric>&sort=<column>`
    2. `GET /products/<id>`
    3. `POST /products` - with a request body containing a `product` JSON payload
    4. `PUT /products/<id>` - with a request body containing a `product` JSON payload
    5. `DELETE /products/<id>`
4. You have a `getProducts` handler that retrieves a list of products from the DB to return as a JSON list to the client, limiting the number of results based on the `limit` parameter and sorting the results based on the `sort` parameter.
    1. Unknown sort fields should be ignored to prevent your SQL statements from causing an error.
    2. You must use parameterized prepared statements.
5. You have a `getProduct` handler that retrieves a single product from the database based on the passed in `id`.
    1. Return a `404 Not Found` to the client for unknown product IDs
    2. Return a `400 Bad Request` if a non-numeric ID is passed in
    3. Returns a `200 OK` and a JSON representation of the `product` back to the client.
6. You have an `addProduct` handler that inserts a new product in the DB.
    1. Return a `400 Bad Request` if any required fields are missing in the request body, including `name` and `price`.
    2. Return a `400 Bad Request` if a non-numeric `price` is passed in.
    3. Returns a `201 Created` and no response body upon success.
7. You have an `updateProduct` handler that updates an existing product in the DB.
    1. Return a `404 Not Found` to the client for unknown product IDs
    2. Return a `400 Bad Request` if a non-numeric ID is passed in
    3. Return a `400 Bad Request` if any required fields are missing in the request body, including `name` and `price`.
    4. Returns a `200 OK` and no response body upon success.
8. You have a `deleteProduct` handler that removes an existing product from the DB.
    1. Return a `404 Not Found` to the client for unknown product IDs
    2. Return a `400 Bad Request` if a non-numeric ID is passed in
    3. Returns a `200 OK` and no response body upon success.
9. Your work is done on a separate branch (e.g. `refactor-product-server-handlers`) and you provide a link to a GitHub pull request (PR) against your repository’s `main` branch with your solution for evaluation.
10. Your PR’s comment includes a copy of the output of your program for each endpoint you call via a client such as `curl`. You must show the `curl` calls and the server’s response underneath.

## Resources

1. [SQLite tutorial](https://www.youtube.com/watch?v=zLQ03DeH04c&list=PL-1QdJ8od_eyxntzYQhwCkcVZlqWVrmSf&index=1)
2. [SQLite data types](https://www.sqlite.org/datatype3.html)
3. Go [database/sql](https://pkg.go.dev/database/sql@go1.19.3) package documentation
4. SQLite driver usage [examples](https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go) in Go with the https://github.com/mattn/go-sqlite3 package

## Keep In Mind

1. Directory names should be lowercased, avoid MixCasedDirectoryNames.
2. Avoid directory bloat. No need for separate `handlers` , `objects` , `server` , `models` , etc folders. Keep it simple.
3. Handle your errors. Do not discard them with the blank identifier ( `_` ).
4. Do not crash your own server with `log.Fatal` within your handlers.
