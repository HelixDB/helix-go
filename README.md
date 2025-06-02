# helix-go
The official Go SDK for HelixDB

## Getting Started

### Setting up HelixDB

#### Install HelixCLI
```bash
curl -sSL "https://install.helix-db.com" | bash
helix install
helix init
```

#### Create a HelixQL schema

```js
// ./helixdb-cfg/schema.hx
N::User {
    name: String,
    age: U32,
    email: String,
    created_at: I32,
    updated_at: I32,
}
```

#### Create HelixQL queries

```js
// ./helixdb-cfg/queries.hx
QUERY create_user(name: String, age: U32, email: String, now: I32) =>
    user <- AddN<User>({name: name, age: age, email: email, created_at: now, updated_at: now})
    RETURN user

QUERY get_users() =>
    users <- N<User>
    RETURN users
```

#### Check for queries (optional)

```bash
helix check
```

#### Deploy HelixQL queries

```bash
helix deploy
```

### Using Go with HelixDB

#### Install helix-go

```bash
go get github.com/HelixDB/helix-go
```

#### Send requests to HelixDB

```go
// ./main.go
var HelixClient *helix.Client

func main() {

	// Connect to client
	HelixClient = helix.NewClient("http://localhost:6969")

	// Create user data

	now := time.Now()

	timestamp := now.Unix()

	timestamp32 := int32(timestamp)

	newUser := map[string]any{
		"name":  "John",
		"age":   21,
		"email": "johndoe@email.com",
		"now":   timestamp32,
	}

	// Create user in Helix
	createUserRes, err := HelixClient.Query(
		"create_user",
		helix.WithData(newUser),
	)
	if err != nil {
		log.Fatalf("Error while creating user: %s", err)
	}

	fmt.Println(createUserRes)

	// Get all users
	users, err := HelixClient.Query("get_users")
	if err != nil {
		log.Fatalf("Error while creating user: %s", err)
	}

	fmt.Println(users)
}
```
