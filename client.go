package neo4j

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"sort"
)

type Records map[string]interface{}
type Result interface{}
type Transaction func(Job) (Result, error)

type Client struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Job interface {
	Execute(Query) ([]Records, error)
}

type job struct {
	transaction neo4j.Transaction
}

func (c *Client) NewRequest() Query {
	return query{}
}

func (c *Client) Read(transaction Transaction) ([]Records, error) {
	return c.run(transaction, neo4j.AccessModeRead)
}

func (c *Client) Write(transaction Transaction) ([]Records, error) {
	return c.run(transaction, neo4j.AccessModeWrite)
}

func (c *Client) run(transaction Transaction, accessMode neo4j.AccessMode) (result []Records, err error) {
	auth := neo4j.BasicAuth(c.Username, c.Password, "")
	addr := fmt.Sprintf("bolt://%s:%d", c.Host, c.Port)
	driver, err := neo4j.NewDriver(addr, auth)
	if err != nil {
		return
	}
	defer driver.Close()
	session, err := driver.Session(accessMode)
	if err != nil {
		return
	}
	defer session.Close()
	workTransaction := func(tx neo4j.Transaction) (interface{}, error) {
		return transaction(&job{transaction: tx})
	}
	var raw interface{}
	switch accessMode {
	case neo4j.AccessModeRead:
		raw, err = session.ReadTransaction(workTransaction)
	case neo4j.AccessModeWrite:
		raw, err = session.WriteTransaction(workTransaction)
	default:
		panic(fmt.Errorf("invalid AccessMode %v", accessMode))
	}
	if err != nil {
		return
	}
	result = raw.([]Records)
	return
}

func (j *job) Execute(query Query) (records []Records, err error) {
	operation, params := query.eval()
	log.Printf("[Query] %s", operation)
	log.Printf("[Params] %s", params)
	result, err := j.transaction.Run(operation, params)
	if err != nil {
		return
	}
	for result.Next() {
		err = result.Err()
		if err != nil {
			return
		}
		record := Records{}
		for _, key := range result.Record().Keys() {
			record[key], _ = result.Record().Get(key)
		}
		records = append(records, record)
	}
	err = result.Err()
	if err != nil {
		return
	}
	log.Printf("[Result] %s", records)
	return
}

func (r Records) Merge(records Records) Records {
	if r == nil {
		return records
	}
	for key, value := range records {
		r[key] = value
	}
	return r
}

func (r Records) Equals(o Records) bool {
	if len(r) != len(o) {
		return false
	}
	if len(r) == 0 {
		return true
	}
	for key, value := range r {
		v, ok := o[key]
		if !ok || v != value {
			return false
		}
	}
	return true
}

func (r Records) Keys() []string {
	var keys []string
	for key := range r {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (r Records) GetOrElse(key string, defaultValue interface{}) interface{} {
	if r == nil {
		return defaultValue
	}
	value, ok := r[key]
	if !ok {
		return defaultValue
	}
	return value
}
