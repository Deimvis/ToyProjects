# QUIZ GAME
_Inspired by [Gophercises](https://courses.calhoun.io/courses/cor_gophercises)_

Play quiz game with custom questions and time limit

## Getting Started

### Quick Start

```bash
go run main.go [options]
```

### Setup

* Put questions into a csv file with `{question},{answer}` format (see [example](problems.csv))
* Answer is considered to be correct regardless of letter case and leading/trailing spaces
* Run program specifying path to the csv file using flag `-problems` and the number of seconds to limit the time of quiz using flag `-time-limit`

  ```go run main.go -problems=PATH_TO_THE_CSV_FILE -time-limit=NUMBER_OF_SECONDS```
