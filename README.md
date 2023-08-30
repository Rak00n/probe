# probe
 Small infrastructure monitoring solution with Telegram integration
 
### Features

- Windows and Linux agents
- Free Space on device tracking (see example configs)
- Linux RAID health monitoring (real reason for the app)
- Telegram Notifications

#### Telegram Notifications

Telegram Bot Token is required. It must be set in configuration file. You can get your token with [Telegram Help Page](https://core.telegram.org/bots/tutorial "Telegram Help Page")

#### Configuration file

Configuration file is in JSON format.
- UUID - unique ID. Any alphanumeric characters. It is not used in current version of the app but will be part of high-level logic in future releases.
- Host - you machine name. This value will be displayed in Telegram notifications.
- Secret - Any alphanumeric characters. It is used to encrypt data transfered between several probes.
- RelayTo - 
