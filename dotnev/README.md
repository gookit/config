# Dotenv

Package `dotenv` that supports importing data from files (eg `.env`) to ENV

## Usage

```go
err := dotenv.Load("./", ".env")
// err := dotenv.LoadExists("./", ".env")

val := dotenv.Get("ENV_KEY")
// Or use 
// val := os.Getenv("ENV_KEY")

// with default value
val := dotenv.Get("ENV_KEY", "default value")
```
