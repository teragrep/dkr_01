package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/containerd/fifo"
	"github.com/docker/docker/api/types/plugins/logdriver"
	"github.com/docker/docker/daemon/logger"
	protoio "github.com/gogo/protobuf/io"
	"github.com/pkg/errors"
	"github.com/teragrep/rlp_05/pkg/RelpBatch"
	"github.com/teragrep/rlp_05/pkg/RelpConnection"
	"github.com/teragrep/rlp_05/pkg/RelpDialer"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Driver struct {
	mu           sync.Mutex
	logs         map[string]*logPair
	idx          map[string]*logPair
	logger       logger.Logger
	relpHostname string
	relpPort     int
	maxRetries   int
}

type logPair struct {
	stream       io.ReadCloser
	info         logger.Info
	relpConn     *RelpConnection.RelpConnection
	relpHostname string
	relpPort     int
	connected    bool
	tlsMode      bool
	maxRetries   int
	hostname     string
	appName      string
}

func (lg *logPair) Close() {
	err := lg.stream.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not close log pair stream for container id: "+lg.info.ContainerID)
	}
}

func newDriver() *Driver {
	return &Driver{
		logs:         make(map[string]*logPair),
		idx:          make(map[string]*logPair),
		relpHostname: "127.0.0.1",
		relpPort:     1601,
		maxRetries:   5,
	}
}

func (d *Driver) StartLogging(file string, logCtx logger.Info) error {
	d.mu.Lock()
	if _, exists := d.logs[file]; exists {
		d.mu.Unlock()
		return fmt.Errorf("logger for %q already exists", file)
	} else {
		fmt.Println("logger created for " + file)
	}
	d.mu.Unlock()

	f, err := fifo.OpenFifo(context.Background(), file, syscall.O_RDONLY, 0700)
	if err != nil {
		return errors.Wrapf(err, "error opening logger fifo: %q", file)
	}

	d.mu.Lock()

	// get relp hostname and port, if given.
	logOptErr := ValidateLogOpts(logCtx.Config)
	if logOptErr != nil {
		return errors.Wrap(logOptErr, "could not validate log opts for driver")
	}
	v, ok := logCtx.Config[RELP_HOSTNAME_OPT]
	if ok {
		d.relpHostname = v
	}

	v, ok = logCtx.Config[RELP_PORT_OPT]
	if ok {
		convPort, convErr := strconv.ParseInt(v, 10, 64)
		if convErr != nil {
			panic("Could not parse port from string " + v)
		}
		d.relpPort = int(convPort)
	}

	// tls mode or not?
	tlsMode := false
	v, ok = logCtx.Config[RELP_TLS_OPT]
	if ok {
		if v == "true" || v == "True" || v == "TRUE" {
			tlsMode = true
		}
	}

	customHostname := ""
	customAppName := ""
	v, ok = logCtx.Config[SYSLOG_HOSTNAME]
	if ok {
		customHostname = v
	}

	v, ok = logCtx.Config[SYSLOG_APPNAME]
	if ok {
		customAppName = v
	}

	lf := &logPair{
		stream:       f,
		info:         logCtx,
		relpHostname: d.relpHostname,
		relpPort:     d.relpPort,
		tlsMode:      tlsMode,
		connected:    false,
		hostname:     customHostname,
		appName:      customAppName,
	}

	// relp connection for log pair
	// initialize relp connection
	if lf.tlsMode {
		lf.relpConn = &RelpConnection.RelpConnection{RelpDialer: &RelpDialer.RelpTLSDialer{}}
	} else {
		lf.relpConn = &RelpConnection.RelpConnection{RelpDialer: &RelpDialer.RelpPlainDialer{}}
	}
	lf.relpConn.Init()

	// try to connect every 500 ms
	fmt.Fprintln(os.Stdout, fmt.Sprintf("Trying to connect file %s to RELP server %v:%v", file, lf.relpHostname, lf.relpPort))
	for !lf.connected {
		_, err = lf.relpConn.Connect(d.relpHostname, d.relpPort)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not connect to relp server due to: '"+err.Error()+"'. Attempting reconnection in 500ms")
			time.Sleep(time.Millisecond * 500)
		} else {
			fmt.Fprintf(os.Stdout, "File %s Connected to RELP server\n", file)
			lf.connected = true
		}
	}

	d.logs[file] = lf
	d.idx[logCtx.ContainerID] = lf

	d.mu.Unlock()
	go consumeLog(lf)

	return nil
}

func (d *Driver) StopLogging(file string) error {
	d.mu.Lock()

	lg, ok := d.logs[file]
	if ok {
		if err := lg.stream.Close(); err != nil {
			return err
		}
		ok := lg.relpConn.Disconnect()
		if !ok {
			fmt.Fprintln(os.Stderr, "Could not disconnect gracefully, forcing")
			lg.relpConn.TearDown()
		}
		delete(d.logs, file)
	}

	d.mu.Unlock()
	return nil
}

func consumeLog(lg *logPair) {
	dec := protoio.NewUint32DelimitedReader(lg.stream, binary.BigEndian, 1e6)

	defer dec.Close()
	defer lg.Close()

	var buf logdriver.LogEntry
	currentRetries := 0
	for {
		if err := dec.ReadMsg(&buf); err != nil {
			if err == io.EOF || err == os.ErrClosed || strings.Contains(err.Error(), "file already closed") {
				// exit loop if EOF or FIFO closed
				fmt.Fprintf(os.Stderr, "FIFO Stream closed %s\n", err.Error())
				return
			}

			if lg.maxRetries != -1 && currentRetries > lg.maxRetries {
				fmt.Fprintf(os.Stderr, "Current amount of retries exceeded the max retries amount. Shutting down logger")
				return
			}

			currentRetries++
			fmt.Fprintf(os.Stderr, "Encountered error, retrying")
			time.Sleep(500 * time.Millisecond)
			dec = protoio.NewUint32DelimitedReader(lg.stream, binary.BigEndian, 1e6)
		}
		currentRetries = 0

		// Write message to RELP server
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%s: [%s] [%d] %s", lg.info.ContainerID, buf.Source, buf.TimeNano, buf.Line))

		batch := RelpBatch.RelpBatch{}
		batch.Init()

		// make syslog message format
		syslogMsg := InitializeSyslogMessage()
		syslogMsg.AddMessagePart(buf.Line)
		syslogMsg.SetProcId(lg.info.ContainerID)
		syslogMsg.SetTimeNano(buf.TimeNano)
		syslogMsg.SetPriority(PRIO_WARNING)

		if lg.hostname != "" {
			syslogMsg.SetHostname(lg.hostname)
		}

		if lg.appName != "" {
			syslogMsg.SetAppName(lg.appName)
		}

		batch.Insert(*syslogMsg.Bytes()) // was buf.line

		// send and verify batch
		for notDone := true; notDone; {
			commitErr := lg.relpConn.Commit(&batch)
			if commitErr != nil {
				fmt.Println("Error committing batch: " + commitErr.Error())
			}

			if !batch.VerifyTransactionAll() {
				batch.RetryAllFailed()
				retryRelpConnection(lg.relpConn, lg.relpHostname, lg.relpPort)
			} else {
				notDone = false
			}
		}

		buf.Reset()
	}
}

// ReadLogs functionality is not (yet) provided by this plugin
func (d *Driver) ReadLogs(info logger.Info, config logger.ReadConfig) (io.ReadCloser, error) {
	return nil, nil
}

// retryRelpConnection attempts to forcefully disconnect & reconnect every 1000ms until succeeds
func retryRelpConnection(relpSess *RelpConnection.RelpConnection, hostname string, port int) {
	relpSess.TearDown()
	var success bool
	var err error
	success, err = relpSess.Connect(hostname, port)
	for !success || err != nil {
		fmt.Println("Got error while retrying relp connection: " + err.Error())
		relpSess.TearDown()
		time.Sleep(1000 * time.Millisecond)
		success, err = relpSess.Connect(hostname, port)
	}
}
