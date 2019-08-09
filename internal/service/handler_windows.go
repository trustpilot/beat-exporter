// +build windows

package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
)

type beatExporterService struct {
	stopCh chan<- bool
}

func (s *beatExporterService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s.stopCh <- true
				break loop
			default:
				log.Error(fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

// SetupServiceListener setups service handler for windows
func SetupServiceListener(stopCh chan<- bool, serviceName string, logger log.StdLogger) error {
	isInteractive, err := svc.IsAnInteractiveSession()
	if err != nil {
		return err
	}

	if !isInteractive {
		go func() {
			err = svc.Run(serviceName, &beatExporterService{stopCh: stopCh})
			if err != nil {
				logger.Printf("Failed to start service: %v", err)
			}
		}()
	}

	return nil
}
