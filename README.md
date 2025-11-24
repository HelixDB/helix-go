```markdown
# helix-go

The official Go SDK for HelixDB 

## Table of Contents

-   [Prerequisites](#prerequisites)
-   [Installation](#installation)
-   [Quick Start](#quick-start)
-   [Client Configuration](#client-configuration)
-   [Making Queries](#making-queries)
-   [Handling Responses](#handling-responses)
-   [Complete Example](#complete-example)
-   [Best Practices](#best-practices)

## Prerequisites

Before using this SDK, ensure you have:

1.  **HelixDB running**: The database should be accessible at your specified host
2.  **HelixQL schema and queries defined**: Your database schema and query endpoints should be deployed

For HelixDB setup, visit the [official documentation](https://docs.helix-db.com).

## Installation

```bash
go get github.com/HelixDB/helix-go
```

## Quick Start

### Basic Setup

```go
package main

import (
    "time"
    "github.com/HelixDB/helix-go"
)

func main() {
    // Create client with default timeout (10 seconds)
    client := helix.NewClient("http://localhost:6969")
    
    // Or with custom timeout
    client = helix.NewClient(
        "http://localhost:6969",
        helix.WithTimeout(30*time.Second),
    )
}
```

### Basic Query Pattern

All queries follow this simple pattern:

```go
res, err := client.Query("<endpoint>", /* optional options... */)
if err != nil {
    // handle error
}

// Choose how to handle the response:
err = res.Scan(&destStruct)              // structured
m, err := res.AsMap()                    // dynamic map
raw := res.Raw()                         // raw bytes
```

Where:

- `<endpoint>` is your HelixQL query name
- The response handling is done via methods on `*Response`:
  - `Scan(...)` for typed decoding
  - `AsMap()` for dynamic access
  - `Raw()` for raw bytes

## Client Configuration

### WithTimeout

Configure how long the client waits for responses:

```go
client := helix.NewClient(
    "http://localhost:6969",
    helix.WithTimeout(5*time.Second),
)
```

## Making Queries

### Creating an HQL Schema

```hql
// schema.hx
N::User {
    name: String,
    age: U32,
    email: String,
    created_at: I32,
}

E::Follows {
    From: User,
    To: User,
    Properties: {
        since: I32,
    }
}
```

### Creating HQL Queries

```hql
// queries.hx
QUERY create_user(name: String, age: U32, email: String, now: I32) =>
    user <- AddN<User>({name: name, age: age, email: email, created_at: now})
    RETURN user 

QUERY get_users() =>
    users <- N<User>
    RETURN users 

QUERY follow(follower_id: ID, followed_id: ID) =>
    follower <- N<User>(follower_id)
    followed <- N<User>(followed_id)
    AddE<Follows>::From(follower)::To(followed)
    RETURN "Success" 

QUERY followers(id: ID) =>
    followers <- N<User>(id)::In<Follows>
    RETURN followers 

QUERY following(id: ID) =>
    following <- N<User>(id)::Out<Follows>
    RETURN following 
```

### Passing Data with WithData

The `WithData` option lets you pass input data to your queries. It accepts multiple data types.

#### Using Maps (Recommended for flexibility)

```go
userData := map[string]any{
    "name": "John",
    "age":  25,
}

res, err := client.Query("create_user", helix.WithData(userData))
// handle err and use res...
```

#### Using Structs (Recommended for type safety)

```go
type UserInput struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

input := UserInput{Name: "John", Age: 25}

res, err := client.Query("create_user", helix.WithData(input))
// handle err and use res...
```

#### Using JSON Strings

```go
jsonData := `{"name": "John", "age": 25}`

res, err := client.Query("create_user", helix.WithData(jsonData))
// handle err and use res...
```

#### Using JSON Bytes

```go
jsonBytes := []byte(`{"name": "John", "age": 25}`)

res, err := client.Query("create_user", helix.WithData(jsonBytes))
// handle err and use res...
```

## Handling Responses

Choose the response method that best fits your needs:

### 1. Scan() - Most Flexible

The most powerful method for handling structured responses.

#### Scan Entire Response

```go
type CreateUserResponse struct {
    User User `json:"user"`
}

res, err := client.Query("create_user", helix.WithData(userData))
if err != nil {
    log.Fatal(err)
}

var response CreateUserResponse
if err := res.Scan(&response); err != nil {
    log.Fatal(err)
}

// Access: response.User
```

#### Scan Specific Fields with WithDest

Extract only the fields you need from the response:

```go
// Single field extraction
res, err := client.Query("get_users")
if err != nil {
    log.Fatal(err)
}

var users []User
if err := res.Scan(helix.WithDest("users", &users)); err != nil {
    log.Fatal(err)
}

// Multiple field extraction
res, err = client.Query("get_users_with_count")
if err != nil {
    log.Fatal(err)
}

var totalCount int
if err := res.Scan(
    helix.WithDest("users", &users),
    helix.WithDest("total_count", &totalCount),
); err != nil {
    log.Fatal(err)
}
```

**When to use WithDest:**

- You only need specific fields from a large response
- The response contains multiple top-level fields
- You want to avoid creating response wrapper structs

### 2. AsMap() - Dynamic Access

Get the response as a Go map for flexible access:

```go
res, err := client.Query("get_users")
if err != nil {
    log.Fatal(err)
}

responseMap, err := res.AsMap()
if err != nil {
    log.Fatal(err)
}

// Access nested data
users := responseMap["users"]
fmt.Println(users)

// Type assertion for further processing
if usersList, ok := responseMap["users"].([]interface{}); ok {
    fmt.Printf("Found %d users\n", len(usersList))
}
```

**When to use AsMap:**

- Response structure is unknown or varies
- For debugging and exploration
- When you need flexible access to response data

### 3. Raw() - Maximum Control

Get the raw byte response from HelixDB:

```go
res, err := client.Query("get_users")
if err != nil {
    log.Fatal(err)
}

rawBytes := res.Raw()

// Process raw JSON
fmt.Println(string(rawBytes))

// Manual unmarshaling
var customResult MyCustomStruct
if err := json.Unmarshal(rawBytes, &customResult); err != nil {
    log.Fatal(err)
}
```

**When to use Raw:**

- You need maximum control over response processing
- For custom JSON unmarshaling logic

## Complete Example

Here's a comprehensive example demonstrating user management and relationships:

```go
// main.go
package main

import (
	"fmt"
	"github.com/HelixDB/helix-go"
	"log"
	"time"
)

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Age       int32  `json:"age"`
	Email     string `json:"email"`
	CreatedAt int32  `json:"created_at"`
}

type CreateUserResponse struct {
	User User `json:"user"`
}

type FollowUserInput struct {
	FollowerId string `json:"follower_id"`
	FollowedId string `json:"followed_id"`
}

func main() {
	// Initialize client
	client := helix.NewClient("http://localhost:6969")

	now := int32(time.Now().Unix())

	// Create first user
	userData1 := map[string]any{
		"name":  "Alice Johnson",
		"age":   28,
		"email": "alice@example.com",
		"now":   now,
	}

	res, err := client.Query("create_user", helix.WithData(userData1))
	if err != nil {
		log.Fatal(err)
	}

	var createResponse1 CreateUserResponse
	if err := res.Scan(&createResponse1); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nCreated user 1: %+v\n", createResponse1.User)

	// Create second user
	userData2 := map[string]any{
		"name":  "Bob Smith",
		"age":   32,
		"email": "bob@example.com",
		"now":   now,
	}

	res, err = client.Query("create_user", helix.WithData(userData2))
	if err != nil {
		log.Fatal(err)
	}

	var createResponse2 CreateUserResponse
	if err := res.Scan(&createResponse2); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nCreated user 2: %+v\n", createResponse2.User)

	// Get all users using WithDest
	res, err = client.Query("get_users")
	if err != nil {
		log.Fatal(err)
	}

	var users []User
	if err := res.Scan(helix.WithDest("users", &users)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nTotal users: %d\n", len(users))

	// Create follow relationship: Alice follows Bob
	followData := &FollowUserInput{
		FollowerId: createResponse1.User.ID,
		FollowedId: createResponse2.User.ID,
	}

	res, err = client.Query("follow", helix.WithData(followData))
	if err != nil {
		log.Fatal(err)
	}

	// Use Raw() for operations that don't return structured data
	_ = res.Raw()

	fmt.Printf("\n%s now follows %s\n", createResponse1.User.Name, createResponse2.User.Name)

	// Get Bob's followers using WithDest
	res, err = client.Query("followers",
		helix.WithData(map[string]any{"id": createResponse2.User.ID}),
	)
	if err != nil {
		log.Fatal(err)
	}

	var followers []User
	if err := res.Scan(helix.WithDest("followers", &followers)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%s has %d followers:\n", createResponse2.User.Name, len(followers))
	for _, follower := range followers {
		fmt.Printf("\t%s\n", follower.Name)
	}

	// Get Alice's following using AsMap for demonstration
	res, err = client.Query("following",
		helix.WithData(map[string]any{"id": createResponse1.User.ID}),
	)
	if err != nil {
		log.Fatal(err)
	}

	followingMap, err := res.AsMap()
	if err != nil {
		log.Fatal(err)
	}

	if followingList, ok := followingMap["following"].([]any); ok {
		fmt.Printf("\n%s is following %d users\n", createResponse1.User.Name, len(followingList))

		for _, userFollowing := range followingList {
			if m, ok := userFollowing.(map[string]any); ok {
				fmt.Printf("\t%v\n", m["name"])
			}
		}
	}

	fmt.Println("Example completed successfully!")
}
```

This example demonstrates:

- **Client initialization** with default settings
- **Creating multiple users** with `WithData` and `Scan`
- **Querying data** with field-specific extraction using `WithDest`
- **Creating relationships** between users using `Raw()` for operations
- **Fetching related data** (followers/following) with different response methods
- **Using AsMap** for flexible response handling

## Best Practices

### Choosing the Right Response Method

- **Use `Scan()`** when you know the response structure and want type safety
- **Use `Scan()` with `WithDest()`** when you only need specific fields from large responses
- **Use `AsMap()`** for exploration, debugging, or when response structure varies
- **Use `Raw()`** when you need custom processing or maximum control

### Input Data Types

- **Prefer structs** for type safety and clearer code
- **Use maps** for flexible input scenarios
- **Avoid slices/arrays** as top-level input — HelixDB expects key–value objects
- **Use JSON strings/bytes** only when you're manually preparing JSON

---

If you encounter issues or want to contribute, feel free to open an issue or submit a PR on the GitHub repository.

Happy querying with HelixDB 🚀
