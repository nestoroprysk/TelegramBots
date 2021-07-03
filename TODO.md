# Technical

- Implement close to a lot of interfaces to assert that all that should be executed actually if
- React to error responses
- Make prettry output actually pretty https://core.telegram.org/bots/api#sendmessage
- Resolve in-code TODOs
- Stabilize
    - Cover all with unit tests
    - Cover functions with injected Telegram and SQL 
    - Create a test bot and run integration tests against it with an in-memory DB
- Add linter
- Local dev
    - Create the local dev environment that starts SQL and function listeners in dockers
    - Makefile
        - `make` spins the local dev
        - `make sql` reloads SQL
        - `make admin` reloads admin
        - `make expenses` reloads expenses
- Commit hooks
    - Lint
    - Build
    - Unit test
- Cloudbuild
    - Lint
    - Build
    - Unit test
    - Integration test
    - System test
    - Deploy only after that
    - Trigger deployed functions and rollback (if possible) on fail
- Email failed builds
- Secure admin
    - Research https://nordicapis.com/developing-secure-bots-using-the-telegram-apis/
- Replace security-related `README.md` steps with the `gcloud` commands
- Support capturing
- Add Github badges

# Functional

- Implement the expenses bot
    - MVP is to run admin per a unique user DB
- Add limits per category
- Add statistics
- Document cloud functions
    - Research tools like Swagger

Ideas:
- Try terraform for the infra initial setup

# Bugs

- Encode response in a normal way, for it's base64 encoded string for now
