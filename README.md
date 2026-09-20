# Envoy

**Envoy** is a secure command-line secret manager written in Go.

It stores secrets locally in an encrypted vault and provides commands for setting, retrieving, listing, deleting, importing, and exporting secrets.

The project was built as a practical exploration of **Go, CLI development, cryptography, filesystem operations, concurrency, testing, and Docker**.

---

## Features

* 🔐 AES-GCM encrypted secret storage
* 🔑 Password-derived encryption keys using `scrypt`
* 🧂 Random salt generation
* 🗂️ Namespaced secrets such as `database/password`
* 📋 List stored secrets
* 🗑️ Delete secrets with confirmation
* 📦 Encrypted vault export
* 📥 Encrypted vault import and merge
* 🔒 File locking for concurrent CLI processes
* 📝 Audit logging without recording secret values
* ⚙️ `.envoyrc` configurable data directory
* 🧪 Integration testing
* 🏁 Race detector testing
* 🐳 Docker support

---

## Installation

Clone the repository:

```bash
git clone https://github.com/Ganesh-12-spec/envoy.git
cd envoy
```

Run Envoy directly with Go:

```bash
go run ./cmd/envoy --help
```

---

## Basic Usage

### Initialize

```bash
go run ./cmd/envoy init
```

This creates the Envoy configuration and initializes the encryption setup.

### Store a secret

```bash
go run ./cmd/envoy set API_KEY "my-secret"
```

### Retrieve a secret

```bash
go run ./cmd/envoy get API_KEY
```

### List secrets

```bash
go run ./cmd/envoy list
```

### Delete a secret

```bash
go run ./cmd/envoy delete API_KEY
```

---

## Namespaces

Secrets can be organized using:

```text
namespace/key
```

For example:

```bash
go run ./cmd/envoy set database/password "mypassword"
go run ./cmd/envoy set database/host "localhost"
go run ./cmd/envoy set api/key "secret"
```

List them with:

```bash
go run ./cmd/envoy list
```

---

## Import and Export

Export the encrypted vault:

```bash
go run ./cmd/envoy export backup.enc
```

Import and merge a backup:

```bash
go run ./cmd/envoy import backup.enc
```

The vault backup remains encrypted rather than exposing plaintext secret values.

---

## `.envoyrc`

Envoy supports a `.envoyrc` file for configuring where its data is stored.

By default, Envoy uses:

```text
.envoy/
```

A `.envoyrc` file can point Envoy to another data directory.

This allows the same Envoy commands to work with a different vault location without changing the application code.

---

## Security

Envoy uses several security mechanisms:

### Password-derived key

A random salt is generated during initialization.

The master password and salt are used to derive the encryption key using `scrypt`.

### AES-GCM

Secrets are encrypted using AES-GCM before being written to the vault.

The vault stores encrypted ciphertext and the corresponding nonce rather than plaintext secrets.

### Password verification

Envoy stores a salted password hash for master-password verification.

### File permissions

Sensitive local files are created with restrictive permissions where appropriate.

### File locking

Envoy uses an operating-system file lock around vault operations that read, modify, or write the vault.

This protects the vault when multiple Envoy processes operate on it at the same time.

### Audit log

Envoy records actions such as:

```text
2026-09-19T13:22:04Z SET database/password
```

The audit log records the action and target, but **does not record the actual secret value or master password**.

---

## Architecture

```text
                    Envoy CLI
                       │
                       ▼
                    Cobra
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
       Commands      Config        Paths
          │            │            │
          └────────────┼────────────┘
                       ▼
                     Vault
                       │
              ┌────────┴────────┐
              ▼                 ▼
          Encryption          Locking
              │                 │
              └────────┬────────┘
                       ▼
                 Local Filesystem
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
       config.json  vault.json  audit.log
```

---

## Project Structure

```text
envoy/
├── cmd/
│   └── envoy/
│       └── main.go
│
├── internal/
│   ├── audit/
│   │   └── audit.go
│   │
│   ├── commands/
│   │   ├── delete.go
│   │   ├── export.go
│   │   ├── get.go
│   │   ├── import.go
│   │   ├── init.go
│   │   ├── list.go
│   │   ├── root.go
│   │   └── set.go
│   │
│   ├── config/
│   │   ├── config.go
│   │   └── integration_test.go
│   │
│   ├── crypto/
│   │   ├── crypto.go
│   │   ├── secret.go
│   │   └── crypto_test.go
│   │
│   ├── lock/
│   │   ├── lock.go
│   │   └── lock_test.go
│   │
│   └── paths/
│       └── paths.go
│
├── Dockerfile
├── .dockerignore
├── go.mod
├── go.sum
└── README.md
```

---

## Testing

Run the complete test suite:

```bash
go test ./...
```

Run the tests with Go's race detector:

```bash
go test -race ./...
```

The race-detector test suite includes verification of Envoy's file-locking behavior across separate processes.

---

## Docker

Build the Docker image:

```bash
docker build -t envoy .
```

Run Envoy:

```bash
docker run --rm envoy
```

Show the CLI help:

```bash
docker run --rm envoy --help
```

The Docker image uses a multi-stage build:

```text
Go source
   │
   ▼
Go builder image
   │
   ▼
Compiled Envoy binary
   │
   ▼
Small runtime image
```

Local vault files and development artifacts are excluded from the Docker build context using `.dockerignore`.

---

## Engineering Concepts Practiced

This project was developed to practice concepts beyond basic CRUD programming:

* Go project structure
* Cobra CLI development
* Error handling
* File I/O
* JSON persistence
* Password handling
* `scrypt`
* AES-GCM
* Random salt and nonce generation
* Namespacing
* File locking
* Multi-process concurrency
* Audit logging
* Configuration management
* Integration testing
* Race detection
* Docker multi-stage builds

---

## Development History

The project was developed incrementally through small commits.

Major milestones included:

```text
CLI foundation
      ↓
Configuration
      ↓
Encryption
      ↓
Encrypted vault
      ↓
Namespaces
      ↓
Delete confirmation
      ↓
Encrypted export/import
      ↓
CLI polish
      ↓
File locking
      ↓
Concurrency testing
      ↓
Audit logging
      ↓
.envoyrc
      ↓
Integration testing
      ↓
Docker
```

The repository history contains the implementation steps behind these features.

---

## Disclaimer

Envoy is a learning project and should not be treated as a production-grade replacement for established secret-management systems without additional security review, hardening, and threat-model validation.

---

## License

This project is currently a personal learning project.
