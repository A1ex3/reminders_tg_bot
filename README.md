# Reminders Telegram Bot
## A simple telegram bot for notifications about upcoming events.
![https://github.com/A1ex3/reminders_tg_bot/blob/main/.images/main.png](https://github.com/A1ex3/reminders_tg_bot/blob/main/.images/main.png?raw=true)

## File `config.json`
| Key | Type | Description |
| - | - | - |
| tgBotApiToken | String | Token of the Telegram bot. ||
| tgBotDebug | Boolean | Will output debugging information. ||
| dateTimeFormats | Array[string] | Contains templates for formatting the date time. ||
| registrationAccess | Boolean | Allows new users to register, via the /start command. ||
| pathToDataBase | String | Contains the path to the sqlite database. ||
| maxCountEventsPerUser | Integer | Determines how many notifications a user can have recorded. ||

## Build
### Program requirements:
- OS Linux/Windows x64
- Make
- Go
- SQLite3
- Git

### Repository cloning
```bash
git clone https://github.com/A1ex3/reminders_tg_bot.git
```

### Creating a database from `schema.sql` file
```bash
make makedb
```

### Program assembly
```bash
make build
```

## How to use
### You need to download the archive with the build for your OS
- linux
```bash
wget https://github.com/A1ex3/reminders_tg_bot/releases/download/latest/reminders_tg_bot_linux_x64.tar.gz
```
- windows powershell
```powershell
Invoke-WebRequest -Uri "https://github.com/A1ex3/reminders_tg_bot/releases/download/latest/reminders_tg_bot_windows_x64.zip" -OutFile "C:\RemindersTgBot"
```

### It is necessary to unpack the archive
- linux
```bash
tar xvzf reminders_tg_bot_linux_x64.tar.gz
```
- windows powershell
```powershell
Expand-Archive -Path C:\RemindersTgBot\reminders_tg_bot_windows_x64.zip -DestinationPath C:\RemindersTgBot
```

### You need to customize the `config.json` configuration file. you need to insert the api-token for the telegram bot into the `tgBotApiToken` field. [More about config.json](#file-configjson).
### or create an environment variable
```bash
export TGBOTAPITOKEN=token_value
```
### Then you need to add commands to the bot, `menu - <description>` and `get - <description>`
### Run App
- linux
```bash
chmod +x reminders_tg_bot
```
```bash
./reminders_tg_bot -config_path="config.json"
```
- windows powershell
```
.\reminders_tg_bot.exe -config_path="config.json"
```

## Docker Container Launch.
### First you need to create a `config.json` file, where it is in the place where it will be stored is usually `/etc/reminders_tg_bot/config.json`
### Next, you need to paste the data into `config.json` from this [file](https://github.com/A1ex3/reminders_tg_bot/blob/main/config.json) file and fill it in.
- If the token is written to `config.json`
```bash
docker run -p 443:443 -p 80:80 --name reminders_tg_bot -v /etc/reminders_tg_bot/config.json:/etc/reminders_tg_bot/config.json -d ghcr.io/a1ex3/reminders_tg_bot:latest
```
- If the token is written to the environment variable `TGBOTAPITOKEN`
```bash
docker run -p 443:443 -p 80:80 --name reminders_tg_bot -e TGBOTAPITOKEN=tgBotApiToken -v /etc/reminders_tg_bot/config.json:/etc/reminders_tg_bot/config.json -d ghcr.io/a1ex3/reminders_tg_bot:latest
```
