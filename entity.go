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

type Candidates struct {
	Success    bool   `json:"success"`
	Interval   string `json:"interval"`
	StartedAt  int    `json:"started_at"`
	FinishedAt int    `json:"finished_at"`
	VoteTotal  int    `json:"vote_total"`
	Data       []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Youtube    string `json:"youtube"`
		VotesCount int    `json:"votes_count"`
	} `json:"data"`
	Winner []interface{} `json:"winner"`
}