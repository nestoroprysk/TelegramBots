TelegramBots is a project which is designed to run as a set of Google Cloud Functions.
All the functions communicate with a single SQL instance.

All the deploy may be done by
1. Creating a bot and a set of functions https://medium.com/swlh/build-a-telegram-bot-in-go-in-9-minutes-e06ad38acef1
2. Creating a SQL instance and establishing communication with it https://cloud.google.com/sql/docs/mysql/connect-functions
3. In all the functions, setting proper environmental variables, uploading the source and setting entry function names (e.g., upload the result of `cd admin && bots.zip *` and indicate the entry name `HandleAdminSQL`) 

Many thanks to the posts which helped a lot
- On deploying GKF https://medium.com/google-cloud/google-cloud-functions-for-go-57e4af9b10da
