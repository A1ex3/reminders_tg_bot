FROM golang:latest

WORKDIR /opt/reminders_tg_bot
COPY . /opt/reminders_tg_bot

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get install -y sqlite3 jq

# Changing the path to the database file
RUN jq '.pathToDataBase = "/var/lib/remindersTgBot.db"' config.json > tmp.json && mv tmp.json config.json

# Create the database and build the project
RUN make makedb && make build

# Clearing out unnecessary files
RUN find . ! \( -name 'reminders_tg_bot' -o -name 'remindersTgBot.db' -o -name 'config.json' -o -name 'README.md' -o -name 'LICENSE' \) -delete

RUN mv remindersTgBot.db /var/lib/remindersTgBot.db && \
    mkdir /etc/reminders_tg_bot && \
    mv config.json /etc/reminders_tg_bot/config.json

CMD ["sh", "-c", "./reminders_tg_bot -config_path='/etc/reminders_tg_bot/config.json'"]