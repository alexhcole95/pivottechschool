
# Build a Product Server

## Objective

Build an HTTP server that supports returning product information such as name and price. The server should also support the ability to retrieve the complete list of products or a single product based on its ID. Adding and updating new products are also expected.

## Learning Goals

1. Understand the use of the [net/http](https://pkg.go.dev/net/http) and [encoding/json](https://pkg.go.dev/encoding/json) standard library packages
2. Leverage third-party libraries like [gorilla/mux](https://github.com/gorilla/mux)
3. Work with JSON data (read from file, unmarshal into a Go struct)
4. Practice using GitHub to house your projects

## Acceptance Criteria

1. Product data is externalized in a JSON file (products.json)
2. When server starts, it reads and unmarshals the data into data structure suitable for lookups
    1. Use [JSON-to-Go](https://mholt.github.io/json-to-go/) to convert a product JSON sample into a Go struct
3. Server endpoints:
    1. `GET /products` returns all products
    2. `GET /products/{id}` returns the specified product or a 404 status code if not found
    3. `POST /products` adds a new product to the server which should be listed on subsequent calls to retrieve all products
    4. `PUT /products/{id}` updates a product’s name and price
    5. `DELETE /products/{id}` removes a product from the server
4. Provide a link to a GitHub pull request (PR) against your repository’s main branch with your solution for evaluation.

## Resources

- [Instructor Video - Building an HTTP Server](https://drive.google.com/file/d/1cF6MNqliUzYUvqbliz7j1QRx3y4wx749/view?usp=sharing)
- [Sample products dataset](https://gist.githubusercontent.com/jboursiquot/259b83a2d9aa6d8f16eb8f18c67f5581/raw/9b28998704fb06f127f13540a4f6e3812f50774b/products.json)

## Keep In Mind

1. Run and test your servers BEFORE you submit them for review. Use curl or another HTTP client of your choosing to issue requests against your running server. You’ll catch a lot of issues simply by testing your own work first.
2. You do not need to persist additions, updates, and deletions back to the file — not a requirement of the assignment.
3. Avoid committing your `.idea` folders.
4. The `.gitignore` file goes at the root of your repository.
5. Initialize with `go mod init` at the root of your repo, not in a subfolder. If your `go.mod` and `go.sum` files are not at the root of your repo, delete them and re-initialize with `go mod init` at the root of your repo.
6. Directory names should be lowercased, avoid MixCasedDirectoryNames
7. Avoid directory bloat. No need for separate `handlers` , `objects` , `server` , `models` , etc folders. Keep it simple.
8. Handle your errors. Every time you use `_` to discard an error result you invite trouble and the ire of your teammates.
9. Do not crash your own server with `log.Fatal` within your handlers.