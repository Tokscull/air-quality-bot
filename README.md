# Air Quality Bot
Telegram bot to determine air quality near you

Usecases:
- Find out the air quality near you
- Set up daily air quality notifications

# Run locally
To run application locally you need:
- Create your own telegram bot, using [BotFather](https://t.me/botfather)
- Get an access token for api.waqi.info, using [Air Quality Open Data Platform](https://aqicn.org/data-platform/token/)
- Install [ngrok](https://ngrok.com/download) or any other tool, that allow you expose a local server to the Internet

Then follow these steps:
1. Create a `.env` file in `air-quality-bot/` directory
2. Copy all from [.env.examle](https://github.com/Tokscull/air-quality-bot/blob/main/.env.examle) to `.env` file
3. Set up a tunnel from `:8080` port using the following command in terminal:
    ````
    ngrok http 8080
    ````
    
    You will see something similar to this message:
    ````
    Session Status                online
    Account                       your_account (Plan: Free)
    Version                       2.3.40
    Region                        United States (us)
    Web Interface                 http://127.0.0.1:4040
    Forwarding                    http://6118-217-97-78-152.ngrok.io -> http://localhost:8080
    Forwarding                    https://6118-217-97-78-152.ngrok.io -> http://localhost:8080

    Connections                   ttl     opn     rt1     rt5     p50     p90
    ````
4. Copy a forwarding URL from the ngrok result screen. Choose the one, that starts with `https://`. Paste it into `.env` file as a `TELEGRAM_BOT_WEBHOOK_URL`
5. Fill in the remaining fields in the `.env` file

# Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change
