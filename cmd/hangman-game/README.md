
# Capstone Project

## Summary

I have chosen to create a hangman game for my capstone project. The game will be a single player game where the user
will be able to play against the computer. The user will be able to select from three different category options that
will then provide a random word to guess. Those options are cars, dogs, and presidents. There will be a total of seven
mistakes per word.

In this project, I will demonstrate the following:
- Building a REST API in Go.
- Integrating with two remote APIs for categories.
- Using Go's sql package and a sql driver to cache words that have already been fetched for a set period of time.
- Deploying my REST API to a third party PaaS.
- Building a user-friendly CLI client application to interact with my REST API client.
- Good coding practices.
- Testing in Go.

## User Stories

### As a user, I would like to invoke the CLI client application with a category that can be selected. After category selection, the client will then randomly select a word from that category to guess.

**Acceptance Criteria**  

When the program is started, it will list the available categories for a random word.

Example:

```
$ run hangman.go

1. Cars
2. Dogs
3. US Presidents
```

Once the category is selected, I expect the CLI client application to choose a word at random from that category to 
guess. The CLI client application will then print that word to the screen with the letters hidden.

Example:


```
$ 3 (US Presidents)

_ _ _ _ _ _  

you have 7 guesses left
```

For whatever reason, if the category selected was invalid, I expect the CLI app to notify me what the selection was
and allow me to try again.

### As a user, if the letter guessed is in the word, I expect the letter to be revealed in the word. If the letter is not in the word, I expect the letter to be added to a list of incorrect guesses.

**Acceptance Critera**  

When a correct letter is guessed, the client will reveal the letter in the word.

Example:

```
$ e

nice, e is in the word!

_ e _ _ _ e
 
you have 7 guesses left

guessed words: e
```

When an incorrect letter is guessed, the client will add the letter to a list of incorrect guesses and deduct a guess.

Example:

```
$ d

sorry, d is not in the word!

_ e _ _ _ e

you have 6 guesses left

guessed words: e, d
```

If the letter guessed was invalid, or already guessed, I expect the CLI app to notify me what the selection was and allow me to
try again.
