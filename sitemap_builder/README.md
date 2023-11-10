# Sitemap Builder

_Inspired by [Gophercises](https://courses.calhoun.io/courses/cor_gophercises)_

Program that builds a sitemap for given url

## Getting Started

### Quick Start

```bash
go run main.go -url https://bebest.pro -depth 2
```

### Setup

* Choose url and max depth for page traversal (in order to limit recursion)
* Run program

  ```bash
  go run main.go -url YOUR_URL -depth YOUR_MAX_DEPTH
  ```
  or
  ```bash
  make run ARGS="-url YOUR_URL -depth YOUR_MAX_DEPTH"
  ```
