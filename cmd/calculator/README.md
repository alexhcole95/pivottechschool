
# Writing Tests

## Objective

Create a calculator library and write the tests that ensure all of its operations account for failure cases.

## Learning Goals

1. Create a library that you import into an executable (all within your homework repository)
2. Use the standard library’s `testing` package
3. Write tests from the “outside in”

## Acceptance Criteria

1. You have a package named `calculator` that exposes functions with the following signatures
    1. `Add(a, b int) int`
    2. `Subtract(a, b int) int`
    3. `Multiply(a, b int) int`
    4. `Divide(a, b int) (int, error)` — Note the special case for this function, handle “divide by zero” errors gracefully by returning an error
2. You have an executable that uses your package (i.e. it imports your package `github.com/<your-username>/<your-pivottechschool-repo>/calculator`)
3. You have a `calculator_test` package (inside a file `caltulator_test.go` to go along with your `calculator.go` file)
4. You make use of the “Table-Driven Testing” testing technique to test multiple scenarios
5. You use subtests to clearly delineate which scenarios are being tested
6. Provide a link to a GitHub pull request (PR) against your repository’s `main` branch with your solution for evaluation

## Resources

- [Recording from 10/31/22 class](https://www.notion.so/Build-a-Product-Server-f20d8d35d0ed424ba30de523b31a81ea)