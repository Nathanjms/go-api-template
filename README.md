# Go API Template

A template to quickly start a new API with Go and the Echo framework, featuring:

- Sentry Integration
- AWS Integration (or can be CloudFlare R2)
- MySQL/MariaDB Integration
- JWT Authentication, using cookies for authentication and authorization
- A Bruno collection for API documentation

## Development

- Copy `.env.example` to `.env` and fill in the details.
  - _Note: Currently there is a dependency even on dev to use CloudFlare R2. I may change this to allow minIO, but R2 has a good free tier so will leave it for now._
- For the JWT Key:
  1. Generate a rsa public/private key pair with:
     1. `openssl genrsa -out private.pem 2048`
     2. `openssl rsa -in private.pem -pubout > public-key.pem`
  2. Copy the contents of the private.pem into your `.env` file under the `RSA_PRIVATE_KEY` key, with the following changes:
     1. REMOVE THE TOP AND BOTTOM LINES (ie. `-----BEGIN RSA PRIVATE KEY-----` and `-----END RSA PRIVATE KEY-----`)
     2. Remove the newline characters so all is on a single line
     3. Wrap in double quotes for safety
  3. Copy the `public-key.pem` into the `cmd/api/internal/jwtHelper` directory.
  4. Delete the original files if you ran the above steps inside your project directory.
- Ensure MySQL/MariaDB database is setup with credentials matching those of the `.env`
- To Run;
  - If using `air`, can run `air ./cmd/api` for hot reloading
  - Else can use `go run ./cmd/api` and then rerun every time a change occurs
- The API collection is saved in this repo as a Bruno collection. Download Bruno and import the collection.
- Visit `http://localhost:3001` to test that it is working!

**Note: This is a WIP and so is likely to change and may have some issues. If you encounter any issues, please raise an issue report on GitHub.**
