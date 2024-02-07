# Weavy Chat GO Lang SDK

This library provides a Go client for interacting with the Weavy Chat API. It allows you to create applications, manage users, issue access tokens, and perform various operations within the Weavy Chat ecosystem.

## Installation

To install the library, use `go get`:

```bash
go get github.com/CeoFred/weavychat
```

## Usage

```go
import "github.com/CeoFred/weavychat"
```

## Documentation 
This library currently supports the following Weavy Chat API methods:

`NewWeavyServer`: Creates a new WeavyServer instance.
`NewApp`: Creates a new app.
`GetApp`: Retrieves an existing app.
`NewUser`: Creates a new user.
`AddUserToApp`: Adds users to an app.
`RemoveUserFromApp`: Removes users from an app.
`AppInit`: Initializes an app.
`GetAccessToken`: Issues an access token for a user.

## Authentication

You need to provide the server URL and API key to authenticate with the Weavy Chat API.

```go
weavyServer := weavychat.NewWeavyServer("your-weavy-server-url", "your-api-key")
```

## Creating Applications

You can create applications using the `NewApp` method:

```go
appRequest := &weavychat.AppRequest{
    ID:          1,
    Type:        weavychat.AppType("your-app-type"),
    UID:         "your-uid",
    DisplayName: "Your App",
    Metadata:    weavychat.Metadata{},
    Tags:        []string{"tag1", "tag2"},
}

app, err := weavyServer.NewApp(context.Background(), appRequest)
if err != nil {
    // Handle error
}
```

## Managing Users

You can create new users using the `NewUser` method:

```go
user := &weavychat.User{
    UID:         "user-uid",
    Email:       "user@example.com",
    GivenName:   "John",
    MiddleName:  "Doe",
    Name:        "John Doe",
    FamilyName:  "Doe",
    Nickname:    "JD",
    PhoneNumber: "+1234567890",
    Comment:     "A new user",
    Picture:     "user-avatar-url",
    Directory:   "directory-id",
    Metadata:    weavychat.Metadata{},
    Tags:        []string{"tag1", "tag2"},
    IsSuspended: false,
}

newUser, err := weavyServer.NewUser(context.Background(), user)
if err != nil {
    // Handle error
}
```

## Access Tokens

You can issue access tokens for users:

```go
accessToken, err := weavyServer.GetAccessToken(context.Background(), "user-uid", 3600)
if err != nil {
    // Handle error
}
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, please create an issue or a pull request on GitHub.

## License

This library is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.