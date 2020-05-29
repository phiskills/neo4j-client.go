package neo4j_test

import (
	"fmt"
	"github.com/phiskills/neo4j-client.go"
	"strings"
	"testing"
)

type example struct {
	operation string
	params    neo4j.Records
}

func (e example) String() string {
	return fmt.Sprintf("[Query] %s\n[Params] %v\n", e.operation, e.params)
}

const errorFormat = "## Invalid Query:\n- received:\n%v\n- expected:\n%v"

var client = &neo4j.Client{}

func TestQuery_Custom(t *testing.T) {
	example := example{
		operation: "MATCH (n)--()--(m) RETURN n, m",
		params:    neo4j.Records{},
	}
	query := client.NewRequest()
	query = query.Custom(example.operation, example.params)
	validate(t, query, example)
}

func TestQuery_Create_Node(t *testing.T) {
	queries := []string{
		"CREATE (user:User:Customer{age: $user_age, name: $user_name})",
		"RETURN user.name, user.age",
	}
	example := example{
		operation: strings.Join(queries, " "),
		params:    neo4j.Records{"user_name": "John", "user_age": 20},
	}
	user := &neo4j.Node{
		Id:     "user",
		Labels: []string{"User", "Customer"},
		Props:  neo4j.Records{"name": "John", "age": 20},
	}
	query := client.NewRequest()
	query = query.Create(user).Return(user.Properties("name", "age")...)
	validate(t, query, example)
}

func TestQuery_Create_Relationship(t *testing.T) {
	queries := []string{
		"MATCH (user:User{id: $user_id})",
		"MATCH (product:Product{id: $product_id})",
		"CREATE (user)-[owns:OWNS{id: $owns_id}]->(product)",
		"RETURN user.id, product.id, owns.id",
	}
	example := example{
		operation: strings.Join(queries, " "),
		params:    neo4j.Records{"user_id": "000", "product_id": "111", "owns_id": "222"},
	}
	user := &neo4j.Node{
		Id:     "user",
		Labels: []string{"User"},
		Props:  neo4j.Records{"id": "000"},
	}
	product := &neo4j.Node{
		Id:     "product",
		Labels: []string{"Product"},
		Props:  neo4j.Records{"id": "111"},
	}
	owns := &neo4j.Relationship{
		Id:        "owns",
		Type:      "OWNS",
		Props:     neo4j.Records{"id": "222"},
		Direction: neo4j.FromOriginToDestination,
	}
	path := &neo4j.Path{
		Origin:       &neo4j.Node{Id: "user"},
		Relationship: owns,
		Destination:  &neo4j.Node{Id: "product"},
	}
	query := client.NewRequest()
	query = query.Match(user)
	query = query.Match(product)
	query = query.Create(path)
	query = query.Return(
		user.Property("id"),
		product.Property("id"),
		owns.Property("id"),
	)
	validate(t, query, example)
}

func TestQuery_Create_Path(t *testing.T) {
	queries := []string{
		"MATCH (user:User{id: $user_id})",
		"MATCH (product:Product{id: $product_id})",
		"CREATE ()<--(user)-[owns:OWNS{id: $owns_id}]->()--(product)",
		"RETURN user.id, product.id, owns.id",
	}
	example := example{
		operation: strings.Join(queries, " "),
		params:    neo4j.Records{"user_id": "000", "product_id": "111", "owns_id": "222"},
	}
	user := &neo4j.Node{
		Id:     "user",
		Labels: []string{"User"},
		Props:  neo4j.Records{"id": "000"},
	}
	product := &neo4j.Node{
		Id:     "product",
		Labels: []string{"Product"},
		Props:  neo4j.Records{"id": "111"},
	}
	owns := &neo4j.Relationship{
		Id:        "owns",
		Type:      "OWNS",
		Props:     neo4j.Records{"id": "222"},
		Direction: neo4j.FromOriginToDestination,
	}
	path := &neo4j.Path{
		Relationship: &neo4j.Relationship{Direction: neo4j.FromDestinationToOrigin},
		Destination: &neo4j.Path{
			Origin:       &neo4j.Node{Id: "user"},
			Relationship: owns,
			Destination: &neo4j.Path{
				Destination: &neo4j.Node{Id: "product"},
			},
		},
	}
	query := client.NewRequest()
	query = query.Match(user)
	query = query.Match(product)
	query = query.Create(path)
	query = query.Return(
		user.Property("id"),
		product.Property("id"),
		owns.Property("id"),
	)
	validate(t, query, example)
}

func TestQuery_Delete_Relationship(t *testing.T) {
	queries := []string{
		"MATCH (user:User{id: $user_id})",
		"MATCH (user)-[owns:OWNS{id: $owns_id}]->(product:Product{id: $product_id})",
		"DELETE owns, product",
	}
	example := example{
		operation: strings.Join(queries, " "),
		params:    neo4j.Records{"user_id": "000", "product_id": "111", "owns_id": "222"},
	}
	user := &neo4j.Node{
		Id:     "user",
		Labels: []string{"User"},
		Props:  neo4j.Records{"id": "000"},
	}
	product := &neo4j.Node{
		Id:     "product",
		Labels: []string{"Product"},
		Props:  neo4j.Records{"id": "111"},
	}
	owns := &neo4j.Relationship{
		Id:        "owns",
		Type:      "OWNS",
		Props:     neo4j.Records{"id": "222"},
		Direction: neo4j.FromOriginToDestination,
	}
	path := &neo4j.Path{
		Origin:       &neo4j.Node{Id: "user"},
		Relationship: owns,
		Destination:  product,
	}
	query := client.NewRequest()
	query = query.Match(user)
	query = query.Match(path)
	query = query.Delete("owns", "product")
	validate(t, query, example)
}

func TestQuery_Merge_Relationship(t *testing.T) {
	queries := []string{
		"MATCH (user:User{id: $user_id})",
		"MERGE (product:Product{id: $product_id})",
		"MERGE (user)-[owns:OWNS{id: $owns_id}]->(product)",
		"RETURN user.id, product.id, owns.id",
	}
	example := example{
		operation: strings.Join(queries, " "),
		params:    neo4j.Records{"user_id": "000", "product_id": "111", "owns_id": "222"},
	}
	user := &neo4j.Node{
		Id:     "user",
		Labels: []string{"User"},
		Props:  neo4j.Records{"id": "000"},
	}
	product := &neo4j.Node{
		Id:     "product",
		Labels: []string{"Product"},
		Props:  neo4j.Records{"id": "111"},
	}
	owns := &neo4j.Relationship{
		Id:        "owns",
		Type:      "OWNS",
		Props:     neo4j.Records{"id": "222"},
		Direction: neo4j.FromOriginToDestination,
	}
	path := &neo4j.Path{
		Origin:       &neo4j.Node{Id: "user"},
		Relationship: owns,
		Destination:  &neo4j.Node{Id: "product"},
	}
	query := client.NewRequest()
	query = query.Match(user)
	query = query.Merge(product)
	query = query.Merge(path)
	query = query.Return(
		user.Property("id"),
		product.Property("id"),
		owns.Property("id"),
	)
	validate(t, query, example)
}

func TestQuery_Optional_Match_Relationship(t *testing.T) {
	queries := []string{
		"MATCH (user:User{id: $user_id})",
		"OPTIONAL MATCH (product:Product{id: $product_id})",
		"OPTIONAL MATCH (user)-[owns:OWNS{id: $owns_id}]->(product)",
		"RETURN user.id, product.id, owns.id",
	}
	example := example{
		operation: strings.Join(queries, " "),
		params:    neo4j.Records{"user_id": "000", "product_id": "111", "owns_id": "222"},
	}
	user := &neo4j.Node{
		Id:     "user",
		Labels: []string{"User"},
		Props:  neo4j.Records{"id": "000"},
	}
	product := &neo4j.Node{
		Id:     "product",
		Labels: []string{"Product"},
		Props:  neo4j.Records{"id": "111"},
	}
	owns := &neo4j.Relationship{
		Id:        "owns",
		Type:      "OWNS",
		Props:     neo4j.Records{"id": "222"},
		Direction: neo4j.FromOriginToDestination,
	}
	path := &neo4j.Path{
		Origin:       &neo4j.Node{Id: "user"},
		Relationship: owns,
		Destination:  &neo4j.Node{Id: "product"},
	}
	query := client.NewRequest()
	query = query.Match(user)
	query = query.Optional().Match(product)
	query = query.Optional().Match(path)
	query = query.Return(
		user.Property("id"),
		product.Property("id"),
		owns.Property("id"),
	)
	validate(t, query, example)
}

func TestQuery_Set_Relationship(t *testing.T) {
	queries := []string{
		"MATCH (user:User{id: $user_id})--(product:Product{id: $product_id})",
		"SET user.age = $user_age, user.name = $user_name",
		"SET product.price = $product_price",
		"RETURN user.name, product.price",
	}
	example := example{
		operation: strings.Join(queries, " "),
		params: neo4j.Records{
			"user_id": "000", "product_id": "111",
			"user_name": "John", "user_age": 21,
			"product_price": 100,
		},
	}
	user := &neo4j.Node{
		Id:     "user",
		Labels: []string{"User"},
		Props:  neo4j.Records{"id": "000"},
	}
	product := &neo4j.Node{
		Id:     "product",
		Labels: []string{"Product"},
		Props:  neo4j.Records{"id": "111"},
	}
	path := &neo4j.Path{
		Origin:      user,
		Destination: product,
	}
	query := client.NewRequest()
	query = query.Match(path)
	query = query.Set(user, neo4j.Records{"name": "John", "age": 21})
	query = query.Set(product, neo4j.Records{"price": 100})
	query = query.Return(
		user.Property("name"),
		product.Property("price"),
	)
	validate(t, query, example)
}

func TestQuery_Merge_OnCreate_OnMatch(t *testing.T) {
	queries := []string{
		"MATCH (user:User{id: $user_id})",
		"MERGE (product:Product{id: $product_id})",
		"MERGE (user)-[owns:OWNS{id: $owns_id}]->(product)",
		"ON CREATE SET owns.created = $owns_created",
		"ON MATCH SET owns.updated = $owns_updated",
		"RETURN user.id, product.id, owns.id",
	}
	example := example{
		operation: strings.Join(queries, " "),
		params: neo4j.Records{
			"user_id": "000", "product_id": "111", "owns_id": "222",
			"owns_created": "TODAY", "owns_updated": "TODAY",
		},
	}
	user := &neo4j.Node{
		Id:     "user",
		Labels: []string{"User"},
		Props:  neo4j.Records{"id": "000"},
	}
	product := &neo4j.Node{
		Id:     "product",
		Labels: []string{"Product"},
		Props:  neo4j.Records{"id": "111"},
	}
	owns := &neo4j.Relationship{
		Id:        "owns",
		Type:      "OWNS",
		Props:     neo4j.Records{"id": "222"},
		Direction: neo4j.FromOriginToDestination,
	}
	path := &neo4j.Path{
		Origin:       &neo4j.Node{Id: "user"},
		Relationship: owns,
		Destination:  &neo4j.Node{Id: "product"},
	}
	query := client.NewRequest()
	query = query.Match(user)
	query = query.Merge(product)
	query = query.Merge(path)
	query = query.OnCreate().Set(owns, neo4j.Records{"created": "TODAY"})
	query = query.OnMatch().Set(owns, neo4j.Records{"updated": "TODAY"})
	query = query.Return(
		user.Property("id"),
		product.Property("id"),
		owns.Property("id"),
	)
	validate(t, query, example)
}

func TestQuery_Where_Conditions(t *testing.T) {
	queries := []string{
		"MATCH (user)",
		"WHERE NOT user.banned AND user.name = $user_name XOR (user.age < $user_age AND user.grade >= $user_grade)",
		"OR (NOT (user.city STARTS WITH $user_city XOR user.country ENDS WITH $user_country))",
		"AND (NOT (NOT user.description CONTAINS $user_description OR user.email =~ $user_email))",
		"AND (user.uuid IS NOT NULL XOR user.admin IS NULL)",
		"RETURN user.name",
	}
	example := example{
		operation: strings.Join(queries, " "),
		params: neo4j.Records{
			"user_name": "John", "user_age": 30, "user_grade": 4,
			"user_city": "Los", "user_country": "ka",
			"user_description": "good", "user_email": ".*@.*",
		},
	}
	user := &neo4j.Node{Id: "user"}
	query := client.NewRequest()
	subCondition1 := user.Property("age").LessThan(30).And(user.Property("grade").GreaterEqual(4))
	subCondition2 := user.Property("city").StartsWith("Los").XOr(user.Property("country").EndsWith("ka"))
	subCondition3 := neo4j.Not(user.Property("description").Contains("good")).Or(user.Property("email").Matches(".*@.*"))
	subCondition4 := user.Property("uuid").IsNotNull().XOr(user.Property("admin").IsNull())
	query = query.Match(user)
	query = query.Where(
		neo4j.Not(user.Property("banned")).And(
			user.Property("name").IsEqual("John"),
		).XOr(
			subCondition1,
		).Or(
			neo4j.Not(subCondition2),
		).And(
			neo4j.Not(subCondition3),
		).And(
			subCondition4,
		),
	)
	query = query.Return(user.Property("name"))
	validate(t, query, example)
}

func TestQuery_With_OrderBy_Skip_Limit(t *testing.T) {
	queries := []string{
		"MATCH (user:User{name: $user_name})",
		"WITH user",
		"ORDER BY user.age DESC",
		"RETURN user.name, user.age",
		"SKIP 10 LIMIT 5",
	}
	example := example{
		operation: strings.Join(queries, " "),
		params:    neo4j.Records{"user_name": "John"},
	}
	user := &neo4j.Node{
		Id:     "user",
		Labels: []string{"User"},
		Props:  neo4j.Records{"name": "John"},
	}
	query := client.NewRequest()
	query = query.Match(user).With("user")
	query = query.OrderBy(user.Property("age")).Desc()
	query = query.Return(user.Properties("name", "age")...)
	query = query.Skip(10).Limit(5)
	validate(t, query, example)
}

func validate(t *testing.T, query neo4j.Query, example example) {
	received := fmt.Sprintf("%s", query)
	fmt.Printf("# Test:\n%s\n", received)
	expected := fmt.Sprintf("%s", example)
	if received != expected {
		t.Errorf(errorFormat, received, expected)
	}
}
