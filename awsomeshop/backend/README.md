# AWSomeShop Backend

## Dependencies Installed

All required Go dependencies have been successfully added to the project:

- ✅ **Gin** (v1.9.1) - Web framework
- ✅ **GORM** (v1.31.1) - ORM library
- ✅ **GORM MySQL Driver** (v1.6.0) - MySQL database driver
- ✅ **JWT** (v4.5.2) - JSON Web Token authentication
- ✅ **bcrypt** (golang.org/x/crypto v0.47.0) - Password hashing
- ✅ **CORS** (v1.3.1) - Cross-Origin Resource Sharing middleware

## Go Version Requirement

**Important:** This project requires **Go 1.21 or higher** to build and run properly.

The current system has Go 1.18.2, which is insufficient for some transitive dependencies (particularly `crypto/sha3`, `slices`, and other standard library packages introduced in Go 1.21+).

### To Upgrade Go:

1. Download Go 1.21+ from https://golang.org/dl/
2. Install the new version
3. Verify installation: `go version`

## Configuration

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

Then edit `.env` with your actual configuration values.

## Next Steps

After upgrading Go to version 1.21+, you can:

1. Build the project: `go build`
2. Run the project: `go run main.go`
3. Run tests: `go test ./...`
