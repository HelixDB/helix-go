# helix-go

The official Go SDK for HelixDB 

## Table of Contents

-   [Prerequisites](#prerequisites)
-   [Installation](#installation)
-   [Quick Start](#quick-start)
-   [Client Configuration](#client-configuration)
-   [Making Queries](#making-queries)
-   [Handling Responses](#handling-responses)
-   [Error Handling](#error-handling)
-   [Complete Example](#complete-example)
-   [Best Practices](#best-practices)
-   [Requirements](#requirements)

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
    client := helix.NewClient("http://localhost:6969", 
        helix.WithTimeout(30*time.Second))
}

```

### Basic Query Pattern

All queries follow this simple pattern:

```go
response := client.Query("<endpoint>").ResponseMethod()

```

Where:

-   `<endpoint>` is your HelixQL query name
-   `ResponseMethod()` determines how you handle the response (e.g `.Scan(&pointerToStruct)`)

## Client Configuration

### WithTimeout

Configure how long the client waits for responses:

```go
client := helix.NewClient("http://localhost:6969", 
    helix.WithTimeout(60*time.Second))

```

## Making Queries

### Passing Data with WithData

The `WithData` option lets you pass input data to your queries. It accepts multiple data types:

#### Using Maps (Recommended for flexibility)

```go
userData := map[string]any{
    "name": "John",
    "age":  25,
}
client.Query("create_user", helix.WithData(userData))

```

#### Using Structs (Recommended for type safety)

```go
type UserInput struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
input := UserInput{Name: "John", Age: 25}
client.Query("create_user", helix.WithData(input))

```

#### Using JSON Strings

```go
jsonData := `{"name": "John", "age": 25}`
client.Query("create_user", helix.WithData(jsonData))

```

#### Using JSON Bytes

```go
jsonBytes := []byte(`{"name": "John", "age": 25}`)
client.Query("create_user", helix.WithData(jsonBytes))

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

var response CreateUserResponse
err := client.Query("create_user", helix.WithData(userData)).Scan(&response)
if err != nil {
    log.Fatal(err)
}
// Access: response.User

```

#### Scan Specific Fields with WithDest

Extract only the fields you need from the response:

```go
// Single field extraction
var users []User
err := client.Query("get_users").Scan(helix.WithDest("users", &users))

// Multiple field extraction
var users []User
var totalCount int
err := client.Query("get_users_with_count").Scan(
    helix.WithDest("users", &users),
    helix.WithDest("total_count", &totalCount),
)

```

**When to use WithDest:**

-   You only need specific fields from a large response
-   The response contains multiple top-level fields
-   You want to avoid creating response wrapper structs

### 2. AsMap() - Dynamic Access

Get the response as a Go map for flexible access:

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

-   Response structure is unknown or varies
-   For debugging and exploration
-   When you need flexible access to response data

### 3. Raw() - Maximum Control

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

-   You need maximum control over response processing
-   For custom JSON unmarshaling logic
-   You just need to know if the operation succeeded or not

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

-   **Client initialization** with default settings
-   **Creating multiple users** with `WithData` and `Scan`
-   **Querying data** with field-specific extraction using `WithDest`
-   **Creating relationships** between users using `Raw()` for operations
-   **Fetching related data** (followers/following) with different response methods
-   **Using AsMap** for flexible response handling

## Best Practices

### Choosing the Right Response Method

-   **Use `Scan()`** when you know the response structure and want type safety
-   **Use `Scan()` with `WithDest()`** when you only need specific fields from large responses
-   **Use `AsMap()`** for exploration, debugging, or when response structure varies
-   **Use `Raw()`** when you need custom processing or maximum control

### Error Handling

Always handle errors appropriately:

```go
if err := client.Query("endpoint").Scan(&result); err != nil {
    // Log the error with context
    log.Printf("Failed to execute query 'endpoint': %v", err)
    // Handle the error based on your application's needs
    return err
}

```

### Input Data Types

-   **Prefer structs** for type safety and clearer code
-   **Use maps** for flexible input scenarios