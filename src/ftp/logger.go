package ftp

import log "github.com/sirupsen/logrus"

type Logger struct{}

func (Logger) Print(sessionId string, message interface{}) {
	log.WithField("sessionId", sessionId).Debug(message)
}

func (Logger) Printf(sessionId string, format string, v ...interface{}) {
	log.WithField("sessionId", sessionId).Debugf(format, v...)
}

func (Logger) PrintCommand(sessionId string, command string, params string) {
	log.WithFields(log.Fields{"sessionId": sessionId, "command": command, "params": params}).Debug("Command Executed")
}

func (Logger) PrintResponse(sessionId string, code int, message string) {
	log.WithField("sessionId", sessionId).WithField("code", code).Debug(message)
}

