package main

type sendText struct {
	Receiver string `json:"to"`
	Message  string `json:"text"`
}

type sendImage struct {
	Receiver string `json:"to"`
	Message  string `json:"text"`
	Image    string `json:"image"`
}
