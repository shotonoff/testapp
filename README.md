# World of Wisdom

World of Wisdom is a CLI tool that allows you to get a random quote from a variety of sources.

## Explanation of the solution

This solution is used *Challenge–response* approach to verify the client.

This approach has chosen because it is simple, secure and stateless in the implementation.

### Challenge–response works as follows:

1. The client connects to the server. 
2. The server generates an arbitrary random data and sends it back to the client with difficulty.
3. The client looks for the *NONCE* number with the *DIFFICULTY* and sends it back to the server.
4. The server verifies the *NONCE* and *HASH* and sends the quote to the client.

### Protocol implementation

The data is transmitted in [gob](https://pkg.go.dev/encoding/gob) format.

The protocol consists of two sections:
1. 
2. Header - contains the length of the data section.
3. Payload - contains the data itself.

#### Protocol messages

##### Challenge

```go
type Challenge struct {
    Difficulty int
    Data       []byte
}
```

This structure is used to send the challenge to the client.

#### Solution

```go
type Solution struct {
    Nonce int
    Hash  []byte
}
```

This structure is used to send the solution to the server.

#### Quote

```go
type Quote struct {
    Text string
}
```

This structure is used to send the quote to the client.

## Build the docker image

To build the docker image, you will need to run the following command:

```bash
make docker/build
```

## How to run the project

1. To run the project, you will need to have Docker installed.
2. Run the following command to run the project:

```bash
make run
```

The output should look like this:

```text
docker-compose up -d
[+] Running 2/2
 ✔ Network app_testapp     Created                                                                                                             0.1s
 ✔ Container app-server-1  Started                                                                                                             0.5s
docker-compose run client
[+] Running 1/0
 ✔ Container app-server-1  Running                                                                                                             0.0s

[*] Received quote: "Beware of false knowledge; it is more dangerous than ignorance"

docker-compose down
[+] Running 2/2
 ✔ Container app-server-1  Removed                                                                                                             0.1s
 ✔ Network app_testapp     Removed
```

## How to run the tests

To run the tests, you will need to run the following command:

```bash
make test
```

__Note__: To run the tests you need installed GO locally.
