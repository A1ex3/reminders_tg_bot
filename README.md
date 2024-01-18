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

### You need to customize the `config.json` configuration file. you need to insert the api-token for the telegram bot into the `tgBotApiToken` field. [More about config.json](#File config.json)
### Then you need to add commands to the bot, `menu - <description>` and `get - <description>`
