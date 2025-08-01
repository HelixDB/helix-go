# helix-go

The official Go SDK for HelixDB - a graph database optimized for real-time queries and complex relationships.

## Installation

```bash
go get github.com/HelixDB/helix-go
```

## Quick Start

### 1. Initialize the Client

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
    client := helix.NewClient("http://localhost:6969", 
        helix.WithTimeout(30*time.Second))
}
```

### 2. Basic Query Structure

All queries follow this pattern:

```go
response := client.Query("endpoint_name", options...).ResponseMethod()
```

Where:
- `endpoint_name` is your HelixQL query name
- `options` configure the request (data, target types, etc.)
- `ResponseMethod()` determines how you handle the response

## Client Configuration

### WithTimeout

Configure how long the client waits for responses:

```go
client := helix.NewClient("http://localhost:6969", 
    helix.WithTimeout(60*time.Second))
```

## Query Options

### WithData

Pass input data to your query. Accepts multiple data types:

```go
// Using a map
userData := map[string]any{
    "name": "John",
    "age":  25,
}
client.Query("create_user", helix.WithData(userData))

// Using a struct
type UserInput struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
input := UserInput{Name: "John", Age: 25}
client.Query("create_user", helix.WithData(input))

// Using JSON string
jsonData := `{"name": "John", "age": 25}`
client.Query("create_user", helix.WithData(jsonData))

// Using JSON bytes
jsonBytes := []byte(`{"name": "John", "age": 25}`)
client.Query("create_user", helix.WithData(jsonBytes))
```

### WithTarget

Specify the expected response type (for future SDK enhancements):

```go
client.Query("get_users", helix.WithTarget([]User{}))
```

## Response Methods

### Scan

The most flexible method for handling responses. Scan the entire response into a struct or use field-specific scanning.

#### Scan Entire Response

```go
type CreateUserResponse struct {
    User User `json:"user"`
}

var response CreateUserResponse
err := client.Query("create_user", helix.WithData(userData)).Scan(&response)
if err != nil {
    log.Fatal(err)
}
// Access: response.User
```

#### Scan with WithDest (Field-Specific)

Extract specific fields from the response by name:

```go
// Single field extraction
var users []User
err := client.Query("get_users").Scan(helix.WithDest("users", &users))
if err != nil {
    log.Fatal(err)
}

// Multiple field extraction
var users []User
var totalCount int
err := client.Query("get_users_with_count").Scan(
    helix.WithDest("users", &users),
    helix.WithDest("total_count", &totalCount),
)
```

**When to use WithDest:**
- When you only need specific fields from a large response
- When the response contains multiple top-level fields
- When you want to avoid creating response wrapper structs

### AsMap

Get the response as a Go map for dynamic access:

```go
responseMap, err := client.Query("get_users").AsMap()
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
- When response structure is unknown or varies
- For debugging and exploration
- When you need flexible access to response data

### Raw

Get the raw byte response from HelixDB:

```go
rawBytes, err := client.Query("get_users").Raw()
if err != nil {
    log.Fatal(err)
}

// Process raw JSON
fmt.Println(string(rawBytes))

// Manual unmarshaling
var customResult MyCustomStruct
err = json.Unmarshal(rawBytes, &customResult)
```

**When to use Raw:**
- When you need maximum control over response processing
- For custom JSON unmarshaling logic
- When working with streaming or large responses
- For debugging raw server responses

## Error Handling

The SDK provides detailed error information:

```go
err := client.Query("invalid_endpoint").Scan(&result)
if err != nil {
    // Errors include HTTP status codes and response body
    log.Printf("Query failed: %v", err)
    // Example error: "404: endpoint not found"
}
```

Common error scenarios:
- **HTTP errors**: Status codes with server response
- **JSON parsing errors**: Invalid response format
- **Type errors**: Incompatible scan destinations
- **Network errors**: Connection timeouts or failures

## Data Type Requirements

### Input Data Constraints

```go
// ✅ Supported input types
map[string]any{"key": "value"}          // Maps
MyStruct{Field: "value"}                // Structs
`{"key": "value"}`                      // JSON strings
[]byte(`{"key": "value"}`)              // JSON bytes

// ❌ Unsupported input types
[]string{"item1", "item2"}              // Slices/Arrays
"plain string"                          // Non-JSON strings
42                                      // Primitive types
```

### Scan Destination Requirements

```go
// ✅ Valid scan destinations (must be pointers)
var user User
err := query.Scan(&user)

var users []User
err := query.Scan(&users)

var result map[string]any
err := query.Scan(&result)

// ❌ Invalid destinations
var user User
err := query.Scan(user)  // Not a pointer

err := query.Scan(nil)   // Nil pointer
```

## Prerequisites

Before using this SDK, ensure you have:

1. **HelixDB running**: The database should be accessible at your specified host
2. **HelixQL schema and queries defined**: Your database schema and query endpoints should be deployed

For HelixDB setup, visit the [official documentation](https://docs.helix-db.com).

## Complete Example

Here's a comprehensive example demonstrating user management and relationships:

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/HelixDB/helix-go"
)

type User struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Age       int32  `json:"age"`
    Email     string `json:"email"`
    CreatedAt int32  `json:"created_at"`
    UpdatedAt int32  `json:"updated_at"`
}

type CreateUserResponse struct {
    User User `json:"user"`
}

type FollowUserInput struct {
    FollowerId string `json:"followerId"`
    FollowedId string `json:"followedId"`
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
    
    var createResponse1 CreateUserResponse
    err := client.Query("create_user", helix.WithData(userData1)).Scan(&createResponse1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created user 1: %+v\n", createResponse1.User)
    
    // Create second user
    userData2 := map[string]any{
        "name":  "Bob Smith",
        "age":   32,
        "email": "bob@example.com",
        "now":   now,
    }
    
    var createResponse2 CreateUserResponse
    err = client.Query("create_user", helix.WithData(userData2)).Scan(&createResponse2)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created user 2: %+v\n", createResponse2.User)
    
    // Get all users using WithDest
    var users []User
    err = client.Query("get_users").Scan(helix.WithDest("users", &users))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Total users: %d\n", len(users))
    
    // Create follow relationship: Alice follows Bob
    followData := &FollowUserInput{
        FollowerId: createResponse1.User.ID,
        FollowedId: createResponse2.User.ID,
    }
    
    // Use Raw() for operations that don't return structured data
    _, err = client.Query("follow", helix.WithData(followData)).Raw()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s now follows %s\n", 
        createResponse1.User.Name, createResponse2.User.Name)
    
    // Get Bob's followers using WithDest
    var followers []User
    err = client.Query("followers", 
        helix.WithData(map[string]any{"id": createResponse2.User.ID})).
        Scan(helix.WithDest("followers", &followers))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s has %d followers: ", createResponse2.User.Name, len(followers))
    for _, follower := range followers {
        fmt.Printf("%s ", follower.Name)
    }
    fmt.Println()
    
    // Get Alice's following using AsMap for demonstration
    followingMap, err := client.Query("following",
        helix.WithData(map[string]any{"id": createResponse1.User.ID})).AsMap()
    if err != nil {
        log.Fatal(err)
    }
    
    if followingList, ok := followingMap["following"].([]interface{}); ok {
        fmt.Printf("%s is following %d users\n", 
            createResponse1.User.Name, len(followingList))
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

## Requirements

- Go 1.24.3 or later
- HelixDB instance running and accessible

## License

This SDK is part of the HelixDB ecosystem. Check the repository for license details.