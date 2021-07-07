# Next

- Refactor responding into a combined responder

# Technical

- Disable pushing if changes present
- Add integration tests (with the mock telegram and the local mysql DB) for the cmd package
- Add linter
- Add prettier and formatter
- Fix generating docs
- Resolve in-code TODOs
- Stabilize
    - Create a test bot and run integration tests against it with an in-memory DB
- Local dev
    - Create the local dev environment that starts SQL and function listeners in dockers
    - Makefile
        - `make admin` reloads admin
        - `make expenses` reloads expenses
- Cloudbuild
    - Lint
    - System test
    - Trigger deployed functions and rollback (if possible) on fail
- Email failed builds
- Secure admin and expenses
- Replace security-related `README.md` steps with the `gcloud` commands
- Add Github badges
- Try terraform for the infra initial setup


# Functional

- Add limits per category
- Add statistics

# Documenration

- Document cloud functions
    - Research tools like Swagger

# Tests

- Admin queries select 1
- Admin executes insert
- Expenses creates a new user
- Expenses uses an already created user
