# Contributing

Thank you for your interest in CheckBag!

## Environment Setup

Setting up your environment follows many of the same installation steps found on the [README](https://github.com/benjaminRoberts01375/CheckBag/blob/main/README.md#installation). Development and Production both use Docker.

Production uses the usual `docker-compose.yml` file, which makes it easy when running `docker compose up -d --build`. As expected, this gives a very paired down development experience.

The development dockerfile `docker-compose.dev.yml` needs to be specified in your docker commands: `docker compose -f docker-compose.dev.yml up -d --build` and `docker compose -f docker-compose.dev.yml down`, and provides a significantly more robust development environment. This environment is a bit slower, but provides live code reloading for both the frontend and backends. This is meant to be launched on `localhost`.

For a nice simple way to reload the docker setup, you can simply run either:

- `docker compose -f docker-compose.dev.yml down && docker compose -f docker-compose.dev.yml up -d --build`
- `docker compose -f docker-compose.dev.yml down && docker compose up -d --build`

depending if you need to launch either the development environment or production one. You can the `docker-compose.dev.yml` for tearing down either.

### Hidden Files

Visual Studio Code is configured to hide some files that are irrelevant like the `node_modules` folder, various GitHub related files, and the `art projects` folder. They're still accessible via your file manager if you really need to edit them.

### Recommended Packages

While the development environment can launch on its own, it's recommended to install:

- Backend language: [Go](https://go.dev/dl/)
- Frontend package manager: [Bun](https://bun.com/)

The rest should be handled automatically by the recommended Visual Studio Code extensions.
