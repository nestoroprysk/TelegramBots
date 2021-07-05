# Next

- Make sure that if tests fail in CI, the pipeline fails
- Run DB tests in the CI using docker-compose
- Add capturing issues
- Add unit tests for the cmd package
- Add integration tests (with the mock telegram and mysql DB) for the cmd package

# Technical

- Add prettier and formatter as a hook
- Collect statistics on at least, success, fail, and error
- Resolve in-code TODOs
- Stabilize
    - Create a test bot and run integration tests against it with an in-memory DB
- Add linter
- Local dev
    - Create the local dev environment that starts SQL and function listeners in dockers
    - Makefile
        - `make admin` reloads admin
        - `make expenses` reloads expenses
- Cloudbuild
    - Lint
    - Integration test
    - System test
    - Deploy only after that
    - Trigger deployed functions and rollback (if possible) on fail
- Email failed builds
- Secure admin
    - Research https://nordicapis.com/developing-secure-bots-using-the-telegram-apis/
- Replace security-related `README.md` steps with the `gcloud` commands
- Add Github badges

# Functional

- Add limits per category
- Add statistics
- Document cloud functions
    - Research tools like Swagger
- Try terraform for the infra initial setup

# Tests

- Admin queries select 1
- Admin executes insert
- Expenses creates a new user
- Expenses uses an already created user
