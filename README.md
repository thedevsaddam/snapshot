# Snapshot


[![Build Status](https://travis-ci.org/thedevsaddam/snapshot.svg?branch=master)](https://travis-ci.org/thedevsaddam/snapshot)
![Project status](https://img.shields.io/badge/version-0.1-green.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/thedevsaddam/snapshot)](https://goreportcard.com/report/github.com/thedevsaddam/snapshot)
[![GoDoc](https://godoc.org/github.com/thedevsaddam/govalidator?status.svg)](https://godoc.org/github.com/thedevsaddam/snapshot)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Robust, Persistent, Key-Value (KV) store purely written in Golang

### Installation
```bash
go get github.com/thedevsaddam/snapshot
```

### Usage

```go
package main

import (
	"fmt"
	"github.com/thedevsaddam/snapshot"
	"time"
)

//make a type to use
type User struct {
	Name, Occupation string
	CreatedAt        time.Time
}

func main() {

	//create a snapshot collection
	userCollection, err := snapshot.New("users")
	if err != nil {
		fmt.Println(err)
	}

	//add item to collection
	userCollection.Put("john", User{Name: "John Doe", Occupation: "Software Engineer", CreatedAt: time.Now()})
	userCollection.Put("jane", User{Name: "Jane Doe", Occupation: "UI/UX Designer", CreatedAt: time.Now()})

	//get an item from collection
	john := User{}
	userCollection.Get("john", &john)
	fmt.Printf("%s is a %s\n", john.Name, john.Occupation) //John Doe is a Software Engineer

	//check an item is exist in a collection
	fmt.Println(userCollection.Has("john")) //true
	fmt.Println(userCollection.Has("tom"))  //false

	//get all the item keys list
	fmt.Println(userCollection.List())

	//get total item count
	fmt.Println(userCollection.TotalItem())

	//remove a key from collection
	userCollection.Remove("john")

	//remove all the keys with collection
	userCollection.Flush()

}

```

### Roadmap
- [ ] Caching
- [ ] Logger
- [ ] Code review

### License
The **snapshot** is a open-source software licensed under the [MIT License](LICENSE.md).
