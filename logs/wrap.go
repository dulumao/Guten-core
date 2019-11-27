package logs

import (
	"github.com/dulumao/Guten-core/env"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Wrap(logger **log.Logger, name string) {
	logDir := fmt.Sprintf("%s_%s", env.Value.Server.LogDir+string(os.PathSeparator)+filepath.Base(os.Args[0]), name)

	if err := Mkdir(logDir); err == nil {
		if fd, err := os.OpenFile(logDir+string(os.PathSeparator)+fmt.Sprintf("%s.log", time.Now().Format("20060102")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err == nil {
			*logger = log.New(fd, "", log.LstdFlags)
		}
	}
}
