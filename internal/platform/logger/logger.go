package logger
import("os";"github.com/sirupsen/logrus")
func New(level,format string)*logrus.Logger{l:=logrus.New();l.SetOutput(os.Stdout);if p,e:=logrus.ParseLevel(level);e==nil{l.SetLevel(p)};if format=="json"{l.SetFormatter(&logrus.JSONFormatter{})}else{l.SetFormatter(&logrus.TextFormatter{FullTimestamp:true})};return l}