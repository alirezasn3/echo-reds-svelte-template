# Echo Redis Svelte Template

This repository is a template for a full-stack application using the following technologies:

- **Go**: A statically typed, compiled programming language designed for simplicity and efficiency.
- **Echo Framework**: A high-performance, extensible, minimalist web framework for Go.
- **Redis**: An in-memory data structure store, used as a database, cache, and message broker.
- **SvelteKit**: A framework for building web applications using Svelte.
- **Tailwind CSS**: A utility-first CSS framework for rapidly building custom user interfaces.

### TODO
- Add listen port and address for Echo and Node servers to the `config.json` file.

### Prerequisites
Make sure you have the following installed:

- [Go](https://golang.org/doc/install)
- [Node.js](https://nodejs.org/) and [npm](https://www.npmjs.com/get-npm)
- [Redis](https://redis.io/download)

### Setup

After cloning the repository, run the following command in frontend directory to install the necessary npm packages:

```sh
npm install
```

### Configuration

Create a `config.json` file in the root directory of the project with the following content:

```json
{
    "redisURL": "redis://localhost",
    "jwtSecret": "jwtsecret",
    "domain": "localhost",
    "allowedOrigins": [
        "http://localhost:8080"
    ]
}
```

### Development

To start the development server, use the `--dev` argument:

```sh
./echo-redis-svelte-template --dev
```

This will start the vite server in development mode with hot reloading for developing the frontend.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.