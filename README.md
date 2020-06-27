# agile-results

Trello Workflows for an [Agile Results](https://gettingresults.com/) lifestyle.

> This project provides workflows for common task creation using Trello.

![Under Construction](https://media1.tenor.com/images/83592060cb2d2cf51e98a5809aeb60d3/tenor.gif)

## How to Use it

Set the following env vars:

- [`TRELLO_API_KEY`](https://trello.com/app-key)
- `TRELLO_API_KEY` (You can gen a token from the API KEY page)
- `TRELLO_LIST_ID` (See instructions below for aquiring this value)

> In order to start creating cards in a list, you need the list ID. To find a
> specific list ID, you need the board ID. See instructions below.

```sh
Run `go run main.go`
```

### Find your Board ID

```sh
curl -s "https://api.trello.com/1/members/me/boards?key=$TRELLO_API_KEY&token=$TRELLO_API_TOKEN" | \
  jq ".[] | { id, name }"
```

### Find your List ID

```sh
curl -s "https://api.trello.com/1/boards/$TRELLO_BOARD_ID/lists?key=$TRELLO_API_KEY&token=$TRELLO_API_TOKEN" | \
  jq ".[] | { id, name }"

```

## MVP

- Create a card in daily column
  - ~~Create card without user input. This should print a link to the user so that
    they can jump right to the card.
  - A verbose mode would be nice to see underlying commands.
  - Modify to accept card name as command line input via prompt
  - 100% test coverage at this point is mandatory.
- Sets a label(s) for new card
  - Cache label data in file (need a way(s) to refresh cache)
  - This feels better as a prompt style program because I don't want to type out
    label names or remember numbers. Also, how do I select multiple in a prompt
    style interface? "Select your labels: 1. Health 2. Relationships 3...." could
    work here via CSV input.
- Append a checklist to said card
  - CSV Input
- Automatically create cards based on checklist
- Appends links for "checklist item cards" to "parent card"
- Appends link to parent card from "checklist item cards"

### Post-MVP

- Set position of card in daily column
  - Or, ask if it's a priorty and that will determine top or bottom. This fits
    the domain model better.
- Provide a description for newly created card
- Create card in any column!
