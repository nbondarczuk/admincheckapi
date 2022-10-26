package config

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"admincheckapi/api/token/jwk"
	v "admincheckapi/api/version"	
)

var Setup *SetupValueSet

//
// checkCmdLineArgs detects usage of cmd line args like -v (version info)
// or -h (default flags). They cause immediate exit but there is a printout
// on the screen with some key info.
//
func checkCmdLineArgs() {
	v := flag.Bool("v", false, "version info")	
	flag.Parse()

	if *v {
		printVersionInfoAndExit()
	}
}

//
// Init gets the contents of file and uses it to make a config
// It may panic. Handling of it is not a required way.
//
func Init(version, build, revision string) {
	v.Set(version, build, revision)
	checkCmdLineArgs()
	
	input, err := LoadConfigYamlFromFile(DEFAULT_CONFIG_FILE_NAME)
	if err != nil {
		panic(err)
	}

	Setup, err = NewSetupValueSet(input)
	if err != nil {
		panic(err)
	}

	//
	// Logging config based on config file or env variable
	//
	
	logrus.SetLevel(Setup.LogLogrusLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00",
	})

	jwk.InitJWKCache()
	
	Setup.Log()
}

//
// LoadConfigYamlFromFile gets the contents of the config file
//
func LoadConfigYamlFromFile(flnm string) ([]byte, error) {
	input, err := ioutil.ReadFile(flnm)
	if err != nil {
		return nil, err
	}

	return input, nil
}
