# probe
 Small infrastructure monitoring solution with Telegram integration
 
### Features

- Windows and Linux agents
- Free Space on device tracking (see example configs)
- Linux RAID health monitoring (real reason for the app)
- Telegram Notifications
- Chaining probes for isolated infrastructure. ( isolatedProbe1 -> isolatedProbe2:49854 -> isolatedProbe3:49854 -> on-lineProbe:49854 -> Telegram)
- AES-encryption for notification transfer

#### Telegram Notifications

Telegram Bot Token is required. It must be set in configuration file. You can get your token with [Telegram Help Page](https://core.telegram.org/bots/tutorial "Telegram Help Page")

#### Configuration file

Configuration file is in JSON format.
- UUID - Unique ID. Any alphanumeric characters. It is not used in current version of the app but will be part of high-level logic in future releases.
- Host - Your machine name. This value will be displayed in Telegram notifications.
- Secret - Any alphanumeric characters. It is used to encrypt data transfered between several probes.
- RelayTo - "telegram" - send notifications via Telegram (requires internet accessibility). "probe" - send notifications to another **probe**
- ListenAddress - host:port to listen for incomming notifications from other probes and relaying them according to **RelayTo**. Complex chains are supported.
- TelegramBotKey - BOT secret key. It is required to send notifications via Telegram. You can get your token with [Telegram Help Page](https://core.telegram.org/bots/tutorial "Telegram Help Page")
- TelegramUserID - Your Telegram User ID. You can get it via third-party services on the Internet or by starting the probe with "TelegramUserID":"000000000" and sending "/getmyid" command into the bot.
- Jobs - Scheduled tasks. Each task must be described as embedded JSON. See example configuration files.
 