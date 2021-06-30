TelegramBots is a GCP project that consists of a family of Telegram bots.

Run the following steps to do the initial setup. After the setup is done, `git push` to `master` deploys all the functions automatically. Magic!

Here are the first time only deploy steps:
- Telegram
  - Ask BotFather for the `ADMIN_BOT_TOKEN`
    ```
    python -m webbrowser https://t.me/botfather
    export ADMIN_BOT_TOKEN=<new-token>
    curl https://api.telegram.org/bot$ADMIN_BOT_TOKEN/getMe
    ```
  - Ask BotFather for the `EXPENSES_BOT_TOKEN`
    ```
    python -m webbrowser https://t.me/botfather
    export EXPENSES_BOT_TOKEN=<new-token>
    curl https://api.telegram.org/bot$EXPENSES_BOT_TOKEN/getMe
    ```
- Setup
  - Install `gcloud` (this works for MACOS)
    ```bash
    brew install google-cloud-sdk
    ```
  - Initialize the project with id like `telegram-bots-new` (this will prompt all the necessary information)
    ```bash
    gcloud init
    ```
- SQL
  - Create SQL instance `bot` (the operation may take a few minutes)
    ```bash
    gcloud sql instances create bot --region=europe-west3 --tier=db-f1-micro
    ```
  - Generate password for the SQL `root` user to `BOT_SQL_ROOT_PASS`
    ```bash
    export BOT_SQL_ROOT_PASS=$(openssl rand -base64 14)
    ```
  - Set the SQL `root` password
    ```bash
    gcloud sql users set-password root --host='%' --instance=bot --password=$BOT_SQL_ROOT_PASS
    ```
  - Verify that `BOT_SQL_ROOT_PASS` is fine by using it to connect
    ```bash
    gcloud sql connect bot --user=root
    ```
  - Ask for the `BOT_SQL_CONNECTION_NAME`
    ```bash
    export BOT_SQL_CONNECTION_NAME=$(gcloud sql instances describe bot --format=json | jq -r .connectionName)
    ```
- Secret
  - Create all the secrets (they will be passed as environmental variables to functions)
    ```bash
    ./create-secrets
    ```
- Cloud Build
  - Create the push to `master` trigger (click on the link to connect the repo and fill-in the desired `--repo-name` and `--repo-owner`)
    ```bash
    gcloud beta builds triggers create github --name=deploy --repo-name=TelegramBots --branch-pattern="^master$" --repo-owner=nestoroprysk --build-config=cloudbuild.yaml
    ```
  - Grant access to develop functions and access secrets with the project id you set up as input (replace `telegram-bots-new` with your id)
    ```bash
    python -m webbrowser https://console.cloud.google.com/cloud-build/settings/service-account?folder=&organizationId=&project=telegram-bots-new 
    ```
- Functions
  - Deploy functions 
    ```bash
    git checkout master && touch drop-me && git add -A && git commit -m "Triggering the deploy of functions" && git push
    ```
  - Set Telegram hooks
    ```bash
    curl --data "url=$(gcloud functions describe Admin --region=europe-west3 --format=json | jq -r .httpsTrigger.url)" https://api.telegram.org/bot$ADMIN_BOT_TOKEN/SetWebhook
    curl --data "url=$(gcloud functions describe Expenses --region=europe-west3 --format=json | jq -r .httpsTrigger.url)" https://api.telegram.org/bot$EXPENSES_BOT_TOKEN/SetWebhook
    ```
  - Add `allUsers` `Cloud Function Invoker` permissions to functions using the UI
- Install Hooks
  - ```bash
    make install-hooks
    ```

Sources:
- Installing `gcloud` https://cloud.google.com/sdk/docs/install
- Creating bots https://medium.com/swlh/build-a-telegram-bot-in-go-in-9-minutes-e06ad38acef1
- Creating SQL instances https://cloud.google.com/sql/docs/mysql/create-instance
- Creating secrets https://cloud.google.com/secret-manager/docs/creating-and-accessing-secrets
- Setting builds
  - https://cloud.google.com/build/docs/automating-builds/create-manage-triggers#gcloud 
  - https://cloud.google.com/sdk/gcloud/reference/beta/builds/triggers/create/github 
- Using secrets https://cloud.google.com/build/docs/securing-builds/use-secrets
- Deploying functions
  - https://cloud.google.com/sdk/gcloud/reference/functions/deploy
  - https://cloud.google.com/functions/docs/env-var
