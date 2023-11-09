# Choose Your Own Story

_Inspired by [Gophercises](https://courses.calhoun.io/courses/cor_gophercises)_

Application that implements an interactive book where users are given options about how they want to proceed the story

Application has two interfaces: Web and CLI

Story is provided with json file

## Getting Started

```bash
make run
```

### Setup

* Put story into json file (see [example](story.json))
* Run program
  1. Web version

     ```bash
     make runweb
     ```

  2. CLI version

     ```bash
     make runcli
     ```

* You can also specify flags using `ARGS=...`, e.g.

  ```bash
  make runweb ARGS="-port 777 -file story.json"
  ```
