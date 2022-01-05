## Store Management System

using Go, gRPC

System uses:

- Client streaming
- Server streaming
- Bi-directional streaming
- JWT Authentication with gRPC interceptors

### Development

```shell

# install dependencies
make install

# generate protocol buffer files (types)
make gen

# remove protocol buffer files (types)
make clean

# run in Development mode
make run

# run all tests
make test
```

### Debugging with [Evans](https://github.com/ktr0731/evans)

```shell
# run server
make server

# run Evans (2nd terminal)
evans -r repl -p 8080
```
