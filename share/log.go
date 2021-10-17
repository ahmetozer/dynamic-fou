package share

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

// Start the logger
func InitLogger(logfile string, level uint8) {
	err := CheckFolder(filepath.Dir(logfile))
	Err(err)

	encoder := getEncoder()
	atom := zap.NewAtomicLevel()
	if level >= 1 && level <= 3 {
		switch level {
		case 1:
			atom.SetLevel(zap.DebugLevel)
		case 2:
			atom.SetLevel(zap.InfoLevel)
		case 3:
			atom.SetLevel(zap.ErrorLevel)
		default:
			atom.SetLevel(zap.InfoLevel)
		}

	}

	var core zapcore.Core
	if logfile == "-" {
		writeSyncer := os.Stdout
		core = zapcore.NewCore(encoder, writeSyncer, atom)
	} else {
		writeSyncer, err := getLogWriter(logfile)
		Err(err)
		core = zapcore.NewCore(encoder, writeSyncer, atom)
	}

	Logger = zap.New(core)
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter(filePath string) (zapcore.WriteSyncer, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(file), nil
}

// Check folder is exist or not
// if not, create the folder
func CheckFolder(folder string) error {
	_, err := os.Stat(folder)

	if os.IsNotExist(err) {
		err2 := os.Mkdir(folder, os.ModePerm)
		if err2 != nil {
			return fmt.Errorf("\"%s\" folder does not exist and cannot be created", folder)
		}
	}
	return nil
}

// err!=nil
func Err(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func LogDefer() {
	Logger.Sync()
}
