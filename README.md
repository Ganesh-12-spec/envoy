# Envoy

A small CLI tool for managing `.env` files and secrets, built with Go.

The goal is to make environment variables easier to manage across development, staging, and production.

### Example

```bash
envoy init
envoy set API_KEY "secret"
envoy get API_KEY
envoy list
envoy env switch staging
```

### Tech

* Go
* Cobra
* AES-GCM
* Local filesystem

🚧 **Currently under development.**

Built as a learning project while exploring Go, backend engineering, CLI tools, and security.
