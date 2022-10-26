package main

import (
	"os"
	
	log "github.com/sirupsen/logrus"

	"admincheckapi/api/token"
)

func initlog(){
	log.SetLevel(log.TraceLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00",
	})
}	

func main() {
	initlog()
	
	for _, str := range os.Args {
		if t, err := token.NewToken([]byte(str)); err != nil {
			log.Errorf("Invalid token: %s", str)			
		} else {
			if groups, err := t.AdminGroups(); err != nil {
				log.Errorf("Error getting admin groups: %s", err)			
			} else {
				log.Printf("Found admin groups: %+v", groups)
			}
		}
	}
}
