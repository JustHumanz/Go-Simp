package main

import "github.com/JustHumanz/Go-Simp/pkg/database"

type MembersPayload struct {
	ID       int64
	NickName string
	EnName   string
	JpName   string
	Region   string
	Fanbase  string
	Status   string
	BiliBili interface{}
	Youtube  interface{}
	Twitter  interface{}
	Twitch   interface{}
	Group    database.Group
	IsLive   []string
}

type GroupPayload struct {
	ID        int64
	GroupName string
	GroupIcon string
	Youtube   interface{}
}
