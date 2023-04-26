package main

import (
	"fmt"
	"os"
	"time"
)

var msgId = 0

type SyslogMessage struct {
	Header         *SyslogHeader
	StructuredData *SyslogStructuredData
	Message        *SyslogInternalMessage
}

type SyslogHeader struct {
	PRI       string
	VERSION   string
	TIMESTAMP string
	HOSTNAME  string
	APPNAME   string
	PROCID    string
	MSGID     string
}

type SyslogStructuredData struct {
	SdElement *SdElement
}

type SyslogInternalMessage struct {
	Msg *[]byte
}

type SdElement struct {
	SdId    string
	SdParam string
}

func InitializeSyslogMessage() *SyslogMessage {
	sysm := &SyslogMessage{}
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Could not get hostname for syslog message, using LOCALHOST")
		hostname = "localhost"
	}

	msgId++

	sysm.Header = &SyslogHeader{
		PRI:       "<0>",
		VERSION:   "0",
		TIMESTAMP: time.Now().UTC().String(),
		HOSTNAME:  hostname,
		APPNAME:   "teragrep",
		PROCID:    "0",
		MSGID:     fmt.Sprintf("%v", msgId),
	}

	sysm.StructuredData = nil
	sysm.Message = nil

	return sysm
}

func (sysm *SyslogMessage) AddMessagePart(msg []byte) {
	sysm.Message = &SyslogInternalMessage{Msg: &msg}
}

func (sysm *SyslogMessage) Bytes() *[]byte {
	msgStr := ""

	// header
	msgStr += sysm.Header.PRI + " " + sysm.Header.VERSION + " " + sysm.Header.TIMESTAMP + " " + sysm.Header.HOSTNAME + " " + sysm.Header.APPNAME + " " + sysm.Header.PROCID + " " + sysm.Header.MSGID
	// structuredData
	if sysm.StructuredData != nil {
		msgStr += " " + sysm.StructuredData.SdElement.SdId + "=" + sysm.StructuredData.SdElement.SdParam
	}
	// msg
	if sysm.Message != nil {
		msgStr += " " + string(*sysm.Message.Msg)
	}

	ret := []byte(msgStr)

	return &ret
}
