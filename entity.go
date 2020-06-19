package main

type FBInfos []struct {
	FbID    string `json:"fb_id"`
	FbName  string `json:"fb_name"`
	FbEmail string `json:"fb_email"`
}

type Vote struct {
	Success    bool   `json:"success"`
	CreatedAt  string `json:"created_at"`
	Left       int    `json:"left"`
	Msg        string `json:"msg"`
	Datetime   string `json:"datetime"`
	Timestamps int    `json:"timestamps"`
}