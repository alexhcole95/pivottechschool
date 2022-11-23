
# Building a Marvel Client

## Objective

Build an HTTP Client to retrieve a list of Marvel characters.

## Learning Goals

1. Learn how to work with third-party APIs
2. Learn to consume and decode JSON responses from REST+JSON APIs
3. Use the standard library’s `net/http` package
4. Use the [Marvel API documentation](https://developer.marvel.com/docs) to guide your usage of their API
5. Learn how to keep secrets out of Git and off of GitHub!

## Acceptance Criteria

1. You have a package `marvel` at the root of your repository containing a `Client` that will expose the behaviors you need from an executable you will also write.
2. You have an executable located in `cmd/marvel` containing your main.go (and any supporting files you need).
3. Your executable uses your `marvel` package to get a client (e.g. `marvel.NewClient()`), passing in whatever parameters you need to initialize a client.
4. Your executable uses your client to retrieve some characters (e.g. `client.GetCharacters()`), passing in whatever parameters you need to retrieve a list of characters.
5. Your executable lists the `name` and `description` of the retrieved characters to STDOUT (your terminal). In addition, your PR’s comment should include a copy of the output of your program.
6. Your work is done on a separate branch (e.g. `add-marvel-client`) and you provide a link to a GitHub pull request (PR) against your repository’s `main` branch with your solution for evaluation
7. DO NOT COMMIT SECRETS to Git or publish them on GitHub.

## Resources

- [Marvel API docs](https://developer.marvel.com/docs)
- [Class recording where instructor builds a simple client using the Marvel API](https://www.notion.so/7878a2a6edfe498f8d873e0cc2175b22)