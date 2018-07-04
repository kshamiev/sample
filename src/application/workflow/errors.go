package workflow // import "application/workflow"

const (
	ErrNone int = iota
	ErrInitConfiguration
	ErrCantChangeWorkDirectory
	ErrInitLogging
	ErrInitJWT
	ErrCantCreateWorkers
	ErrCantStartWorkers
	ErrAllWorkersStopWithError
	ErrCatchPanic
	ErrPidFile
	ErrFailedDeletePIDFile
	ErrCantOpenSocket
)

var errText = map[int]string{
	ErrNone:                    "",
	ErrInitConfiguration:       "Error loading configuration: %s\n",
	ErrCantChangeWorkDirectory: "Can't change work directory to %s\n",
	ErrInitLogging:             "Error initialize logging: %s\n",
	ErrInitJWT:                 "Error initialize JWT keys: %s\n",
	ErrCantCreateWorkers:       "Can't create worker: %s\n",
	ErrCantStartWorkers:        "Can't start worker: %s\n",
	ErrAllWorkersStopWithError: "All workers has stopped with error: %s\n",
	ErrCatchPanic:              "Catch panic: %s\n",
	ErrPidFile:                 "Can't create PID file, incorrectly specified a file name, no access to folder, or another process already running. Error: %s",
	ErrFailedDeletePIDFile:     "Failed to delete pid file: %s",
	ErrCantOpenSocket:          "Can not connect to server: %s",
}
