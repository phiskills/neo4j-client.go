# Phi Skills Neo4j Client for Go

| **Homepage** | [https://phiskills.com][0]        |
| ------------ | --------------------------------- | 
| **GitHub**   | [https://github.com/phiskills][1] |

## Overview

This project contains the Go module to create a **Neo4j client**.  

## Installation

```bash
go get github.com/phiskills/neo4j-client.go
```

## Quick start

```go
package main
import "github.com/phiskills/neo4j-client.go"

client := &neo4j.Client{
	Host: "localhost",
	Port: 7687,
	Username: "neo4j",
	Password: "test",
}
result, err := client.Write(func(job neo4j.Job) (neo4j.Result, error) {
    user := &neo4j.Node{
        Id:     "user",
        Labels: []string{"User", "Customer"},
        Props:  neo4j.Records{"name": "John", "age": 20},
    }
    query := client.NewRequest()
    query = query.Create(user).Return(user.Property("id"))
    records, err := job.Execute(query)
    return records, err
})
for _, record := range result {
    fmt.Printf("user.id = %", record["user.id"])
}
```
For more details, see [Neo4j - CYPHER MANUAL: Chapter 3. Clauses][10].

[0]: https://phiskills.com
[1]: https://github.com/phiskills
[10]: https://neo4j.com/docs/cypher-manual/current/clauses/
