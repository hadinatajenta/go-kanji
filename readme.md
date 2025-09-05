# gobackend

Go backend boilerplate with clean architecture vibes. Structured, scalable, gak ribet.

## 🧱 Project Struc
```
├── app # Wiring & provider injection
├── core # Core logic: config binding, contracts, env
├── env # External config (YAML, etc)
├── infra # Infra layer (logging, DB, etc)
├── shared # Common utils & response handling
├── src # Feature modules / domain logic
├── main.go # Entry point
```
## 🧪 Tech Stack

- Golang
- Clean Architecture style
- Wire for DI
- Modular file structure
- No magic, just straight code

## 🚀 Getting Started

go mod tidy
go run main.go

## 📁 Config

Edit your env in env/app-config-map.yml.
