# whatsappWebAPI

This project is made with the help of [Rhymen/go-whatsapp - WhatsApp Web API](https://github.com/Rhymen/go-whatsapp).
Works only for numbers in India for now.

## Setup

1. Download the latest binary release.  
2. run on your system.  
3. Create a files folder in the directory where binary is running

or create your binary locally as below

Note : (Requirements) Golang setup on your local 

1. Download the repo
2. run commands
```shell
go get .
go build main.go
````  
3. Run the binary

## Steps to use - 

1. Scan the QR code with whatsapp web
2. Put your lists for bulk message and pictures to send in files folder
3. Use the Following links to send text/image

## Usage

Press the following number to send message on whatsapp  
Test --> 0  
Send Text --> 1  
Send Image --> 2  
Send bulk text --> 3  
Send bulk image --> 4  
Exit --> 5  

Demo bulk file examples stored in [files](/files) folder

## Extra

Sample with API hosting, Sentry and bot integrated in hosting folder

## License

![GitHub](https://img.shields.io/github/license/Piyushhbhutoria/whatsappWebAPI?style=for-the-badge)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/Piyushhbhutoria/whatsappWebAPI?style=for-the-badge)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/Piyushhbhutoria/whatsappWebAPI?style=for-the-badge)
