package workflow // import "application/workflow"

// Коды ошибок приложения
const (
	ErrNone                    uint8 = iota // 000
	ErrInitConfiguration                    // 001
	ErrCantChangeWorkDirectory              // 002
	ErrInitLogging                          // 003
	ErrInitJWT                              // 004
	ErrCantCreateWorkers                    // 005
	ErrCantStartWorkers                     // 006
	ErrAllWorkersStopWithError              // 007
	ErrCatchPanic                           // 008
	ErrPidFile                              // 009
	ErrPidFileDelete                        // 010
	ErrFailedDeletePIDFile                  // 011
	ErrCantOpenSocket                       // 012
	ErrApplicationError                     // 014
)

var errText = map[uint8]string{
	ErrNone:                    "%s",
	ErrInitConfiguration:       "Error loading configuration: %s",
	ErrCantChangeWorkDirectory: "Can't change work directory to %s",
	ErrInitLogging:             "Error initialize logging: %s",
	ErrInitJWT:                 "Error initialize JWT keys: %s",
	ErrCantCreateWorkers:       "Can't create worker: %s",
	ErrCantStartWorkers:        "Can't start worker: %s",
	ErrAllWorkersStopWithError: "All workers has stopped with error: %s",
	ErrCatchPanic:              "Catch panic: %s",
	ErrPidFile:                 "Can't create PID file, incorrectly specified a file name, no access to folder, or another process already running. Error: %s",
	ErrPidFileDelete:           "Delete PID file error: %s",
	ErrFailedDeletePIDFile:     "Failed to delete pid file: %s",
	ErrCantOpenSocket:          "Can not connect to socket: %s",
	ErrApplicationError:        "%s",
}
