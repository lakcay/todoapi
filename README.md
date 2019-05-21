# TODO API

RESTful API to create, read, update and delete TODO items.

## Getting Started

I used the GO programming language in this project, although I haven't used it before. Because, in terms of API development, GO offers many advantages over other programming languages. For example, concurrency model in Golang ensures faster performance than other programming languages.

### Prerequisites

* [Go] (https://golang.org/dl/) - Downloading Go programming language is needed
* [mux] (https://github.com/gorilla/mux/blob/master/README.md) - mux is used for router. It can be downloaded with this command

```
go get -u github.com/gorilla/mux
```

### Run

```
go build
./main
```

## Requests

### Create TODO

```
/todos POST action

Example of JSON

 {
    "title": "Buying books",
    "description": "Getting document for exam",
    "tags": ["School", "Exam"],
    "date": "2019-07-15T00:00:00Z"
}
```

### Get All TODOs

```
/todos GET action

Example of result

 {
    "id": "9566c74d-1003-7c4d-7bbb-0407d1e2c649",
    "title": "Buying books",
    "description": "Getting document for exam",
    "tags": [
        "School",
        "Exam"
    ],
    "date": "2019-07-15T00:00:00Z",
    "completed": false
}
```

### Get TODO According to Title

```
/todos/title/{title} GET action
```

### Get TODO According to Description

```
/todos/description/{description} GET action
```

### Get TODO According to Tags

```
/todos/tag/{tag} GET action
```



