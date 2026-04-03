![logowhatsappWebAPI](https://github.com/Piyushhbhutoria/whatsappWebAPI/assets/4961282/211c86c0-a10d-4d84-ba27-eac2d1ce6bba)

# whatsappWebAPI

![Go-Build](https://github.com/Piyushhbhutoria/whatsappWebAPI/workflows/Go-Build/badge.svg)
![GitHub](https://img.shields.io/github/license/Piyushhbhutoria/whatsappWebAPI)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/Piyushhbhutoria/whatsappWebAPI)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Piyushhbhutoria/whatsappWebAPI)
[![Go Report Card](https://goreportcard.com/badge/github.com/Piyushhbhutoria/whatsappWebAPI)](https://goreportcard.com/report/github.com/Piyushhbhutoria/whatsappWebAPI)
![GitHub All Releases](https://img.shields.io/github/downloads/Piyushhbhutoria/whatsappWebAPI/total)
![GitHub repo size](https://img.shields.io/github/repo-size/Piyushhbhutoria/whatsappWebAPI)

A WhatsApp Web API built with Go using [whatsmeow](https://github.com/tulir/whatsmeow). Send messages, images, and manage WhatsApp interactions programmatically.

**Important:** Phone numbers must include the country code (e.g., `91XXXXXXXXXX` for India).

## Setup

**Option 1:** Download the latest binary from [releases](https://github.com/Piyushhbhutoria/whatsappWebAPI/releases)

**Option 2:** Build from source (requires Go 1.26+)

```bash
git clone https://github.com/Piyushhbhutoria/whatsappWebAPI.git
cd whatsappWebAPI
make run
```

**Configuration flags:**

- `-debug`: Enable debug logs
- `-db-dialect`: `sqlite3` or `postgres` (default: `sqlite3`)
- `-db-address`: Database connection string

## Usage

1. Run the application and scan the QR code with WhatsApp Web
2. Use commands interactively:

### Commands

**Messaging:**

- `send <jid> <text>` - Send text message
- `sendimg <jid> <image path> [caption]` - Send image
- `sendbulk <csv file>` - Bulk text (CSV: `<jid>,<message>`)
- `sendbulkimg <csv file>` - Bulk images (CSV: `<jid>,<image path>,[caption]`)

**User Management:**

- `checkuser <phone numbers...>` - Check if users are on WhatsApp
- `getuser <jids...>` - Get user info
- `getavatar <jid> [preview]` - Get user avatar

**Groups:**

- `listgroups` - List all groups
- `getgroup <group_jid>` - Get group info
- `getinvitelink <group_jid> [--reset]` - Get/reset invite link
- `queryinvitelink <link>` - Query invite link info
- `joininvitelink <link>` - Join group via invite link

**Presence & Privacy:**

- `subscribepresence <jid>` - Subscribe to presence updates
- `presence <presence_type>` - Send presence
- `chatpresence <presence_type> <jid> [media_type]` - Send chat presence
- `privacysettings` - Get privacy settings

**Utility:**

- `reconnect` - Reconnect to WhatsApp
- `logout` - Logout from WhatsApp
- `appstate <types...> [resync]` - Sync app state
- `Ctrl+C` - Exit

### CSV Format Examples

**Bulk Text:**

```csv
919876543210,Hello from bulk message
919876543211,Another message
```

**Bulk Image:**

```csv
919876543210,/path/to/image1.jpg,Caption 1
919876543211,/path/to/image2.jpg,Caption 2
```

Demo files available in [files](/files) folder.

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FPiyushhbhutoria%2FwhatsappWebAPI.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FPiyushhbhutoria%2FwhatsappWebAPI?ref=badge_large)
