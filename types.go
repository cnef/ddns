package main

type Record struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type RecordList struct {
	Records []Record `json:"records"`
}
