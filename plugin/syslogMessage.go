package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

var msgId = 0

type SyslogPriority string

const (
	PRIO_EMERGENCY SyslogPriority = "0"
	PRIO_ALERT     SyslogPriority = "1"
	PRIO_CRITICAL  SyslogPriority = "2"
	PRIO_ERROR     SyslogPriority = "3"
	PRIO_WARNING   SyslogPriority = "4"
	PRIO_NOTICE    SyslogPriority = "5"
	PRIO_INFO      SyslogPriority = "6"
	PRIO_DEBUG     SyslogPriority = "7"
)

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

	sysm.Header = &SyslogHeader{
		PRI:       "<0>",
		VERSION:   "0",
		TIMESTAMP: time.Now().Format(time.RFC3339),
		HOSTNAME:  hostname,
		APPNAME:   "teragrep",
		PROCID:    "0",
		MSGID:     fmt.Sprintf("%v", msgId),
	}

	sysm.StructuredData = nil
	sysm.Message = nil

	msgId++

	return sysm
}

func (sysm *SyslogMessage) AddMessagePart(msg []byte) {
	sysm.Message = &SyslogInternalMessage{Msg: &msg}
}

func (sysm *SyslogMessage) SetPriority(prio SyslogPriority) {
	if sysm.Header != nil {
		sysm.Header.PRI = fmt.Sprintf("<%s>", prio)
	}
}

func (sysm *SyslogMessage) SetVersion(ver string) {
	if sysm.Header != nil {
		sysm.Header.VERSION = ver
	}
}

func (sysm *SyslogMessage) SetProcId(id string) {
	if sysm.Header != nil {
		sysm.Header.PROCID = id
	}
}

func (sysm *SyslogMessage) SetTimeNano(nano int64) {
	if sysm.Header != nil {
		sysm.Header.TIMESTAMP = time.Unix(0, nano).Format(time.RFC3339)
	}
}

func (sysm *SyslogMessage) SetHostname(hn string) {
	if sysm.Header != nil {
		sysm.Header.HOSTNAME = hn
	}
}

func (sysm *SyslogMessage) SetAppName(an string) {
	if sysm.Header != nil {
		sysm.Header.APPNAME = an
	}
}

func (sysm *SyslogMessage) Bytes() *[]byte {
	buf := bytes.NewBuffer(make([]byte, 0, 512))

	// header
	buf.WriteString(sysm.Header.PRI)
	buf.WriteByte(' ')
	buf.WriteString(sysm.Header.VERSION)
	buf.WriteByte(' ')
	buf.WriteString(sysm.Header.TIMESTAMP)
	buf.WriteByte(' ')
	buf.WriteString(sysm.Header.HOSTNAME)
	buf.WriteByte(' ')
	buf.WriteString(sysm.Header.APPNAME)
	buf.WriteByte(' ')
	buf.WriteString(sysm.Header.PROCID)
	buf.WriteByte(' ')
	buf.WriteString(sysm.Header.MSGID)

	// structuredData
	if sysm.StructuredData != nil {
		buf.WriteByte(' ')
		buf.WriteString(sysm.StructuredData.SdElement.SdId)
		buf.WriteByte('=')
		buf.WriteString(sysm.StructuredData.SdElement.SdParam)
	}

	// msg
	if sysm.Message != nil {
		buf.WriteByte(' ')
		buf.Write(*sysm.Message.Msg)
	}

	ret := buf.Bytes()
	return &ret
}
