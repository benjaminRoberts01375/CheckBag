# Contributing

Thank you for your interest in CheckBag!

## Environment Setup

Setting up your environment is pretty simple thanks to Docker. Assuming Docker is already installed, simply clone the repository, `cd` into it, and run either of the following:

### Local Development

```bash
docker compose -f docker-compose.dev.yml up -d --build
```

### Production-alike

```bash
docker compose up -d --build
```

## Hidden Files

Visual Studio Code is configured to hide some files that are irrelevant like the `node_modules` folder, various GitHub related files, and the `art projects` folder. They're still accessible via your file manager if you really need to edit them.

## Recommended Packages

While the development environment doesn't _require_ anything other than Docker, it may be helpful to install:

- Backend language: [Go](https://go.dev/dl/)
- Frontend package manager: [Bun](https://bun.com/)

## Sending in Patches

CheckBag uses the typical flow for contributing:

1. Fork the repo
2. Clone: `git clone https://github.com/<your-username>/CheckBag`
3. Add this CheckBag repo as a remote to your version: `git remote add upstream https://github.com/benjaminRoberts01375/CheckBag`
4. Create a new branch on your copy: `git checkout -b feature/my-new-feature`
5. Make your changes
6. Add your changes: `git add .` and `git commit -m "Fixed this bug"`
7. Push your changes: `git push origin feature/my-new-feature`
8. Then finally open a [pull request](https://github.com/benjaminRoberts01375/CheckBag/pulls).
