TelegramBots is a project which is designed to run as a set of Google Cloud Functions.
All the functions communicate with a single SQL instance.

Infrastracture setup may be done by:
1. Creating a bot and a set of functions https://medium.com/swlh/build-a-telegram-bot-in-go-in-9-minutes-e06ad38acef1
2. Creating a SQL instance and establishing communication with it https://cloud.google.com/sql/docs/mysql/connect-functions
3. Setting proper environmental variables for all the functions
4. Creating push to master triggers https://cloud.google.com/build/docs/automating-builds/create-manage-triggers 
5. Emailing build results https://stackoverflow.com/questions/62449482/email-notification-to-a-user-once-the-build-is-successfully-deployed-using-googl 

Many thanks to the posts that helped a lot:
- On deploying GKF https://medium.com/google-cloud/google-cloud-functions-for-go-57e4af9b10da
