![logowhatsappWebAPI](https://github.com/Piyushhbhutoria/whatsappWebAPI/assets/4961282/211c86c0-a10d-4d84-ba27-eac2d1ce6bba)

# whatsappWebAPI

![Go-Build](https://github.com/Piyushhbhutoria/whatsappWebAPI/workflows/Go-Build/badge.svg)
![GitHub](https://img.shields.io/github/license/Piyushhbhutoria/whatsappWebAPI)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/Piyushhbhutoria/whatsappWebAPI)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Piyushhbhutoria/whatsappWebAPI)
[![Go Report Card](https://goreportcard.com/badge/github.com/Piyushhbhutoria/whatsappWebAPI)](https://goreportcard.com/report/github.com/Piyushhbhutoria/whatsappWebAPI)
![GitHub All Releases](https://img.shields.io/github/downloads/Piyushhbhutoria/whatsappWebAPI/total)
![GitHub repo size](https://img.shields.io/github/repo-size/Piyushhbhutoria/whatsappWebAPI)

This project is made with the help of [Rhymen/go-whatsapp - WhatsApp Web API](https://github.com/Rhymen/go-whatsapp).
Works only for numbers in India for now.

## Setup

1. Download the latest binary release.
2. run on your system.

or create your binary locally as below

Note : (Requirements) Golang setup on your local

1. Download the repo
2. run commands below

```go
go get .
go build .
```

3. Run the binary

## Steps to use -

1. Scan the QR code with whatsapp web
2. Put your lists for bulk message and pictures to send in same folder

## Usage

Press the following number to send message on whatsapp  

```
Send Text -> send <jid> <text>
Send Image -> sendimg <jid> <image path> [caption]
Send Bulk Text -> sendbulk <csv file>
Send Bulk Image -> sendbulkimg <csv file>
Exit -> Crtl+C
```

Demo bulk file examples stored in [files](/files) folder

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FPiyushhbhutoria%2FwhatsappWebAPI.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FPiyushhbhutoria%2FwhatsappWebAPI?ref=badge_large)
