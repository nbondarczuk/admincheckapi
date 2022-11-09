package config

import (
	//"flag"
	"admincheckapi/api/aws/awssm"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	v "admincheckapi/api/version"
)

// SetupValueSet is a ready to use, parsed structure, contrary to raw config from
// yaml or env variables
type SetupValueSet struct {
	ConfigFileName               string
	ServerIPAddress              string
	ServerPort                   string
	UsedBackend                  string
	TenantId                     string
	ClientId                     string
	Authority                    string
	Scopes                       []string
	ClientSecret                 string
	AdminGroupName               string
	UseGroupNamePattern          bool
	LogLogrusLevel               log.Level
	LogGORM                      bool
	LogHTTP                      bool
	Backends                     []string
	SQLMaxIdleConns              int
	SQLMaxOpenConns              int
	SQLMaxLifetime               time.Duration
	SecretNamePrefix             string
	AWSUseSecretStore            bool
}

//
// String makes SetupValueSet a Stringer (implicitely)
//
func (s *SetupValueSet) String() string {
	return fmt.Sprintf("%T%+v", *s, *s)
}

//
// Log makes it a Logger
//
func (s *SetupValueSet) Log() {
	log.Infoln("Starting admincheckapi")
	log.Infoln(" Version: " + v.Version)
	log.Infoln("   Build: " + v.Build)
	log.Infoln("Revision: " + v.Revision)
	log.Infoln("Config setup")
	log.Infoln("        Config file name: " + s.ConfigFileName)
	log.Infoln("    HTTP ServerIPAddress: " + s.ServerIPAddress)
	log.Infoln("         HTTP ServerPort: " + s.ServerPort)
	
	log.Infoln("           MSAD TenantId: " + s.hideSecretIfReq(s.TenantId))
	log.Infoln("           MSAD ClientId: " + s.hideSecretIfReq(s.ClientId))
	log.Infoln("          MSAD Authority: " + s.hideSecretIfReq(s.Authority))
	log.Infoln("             MSAD Scopes: " + s.hideSecretIfReq(fmt.Sprintf("%v", s.Scopes)))
	log.Infoln("       MSAD ClientSecret: " + s.hideSecretIfReq(s.ClientSecret))
	log.Infof( "MSAD UseGroupNamePattern: " + os.Getenv("MSAD_USE_GROUP_NAME_PATTERN"))

	log.Infoln("    AWS Use Secret Store: " + os.Getenv("AWS_USE_SECRET_STORE"))
	log.Infoln("              AWS Region: " + os.Getenv("AWS_REGION"))
	log.Infoln("       AWS Access Key ID: " + s.hideSecretIfReq(os.Getenv("AWS_ACCESS_KEY_ID")))
	log.Infoln("   AWS Secret Access Key: " + s.hideSecretIfReq(os.Getenv("AWS_SECRET_ACCESS_KEY")))
	log.Infof( "  AWS Secret Name Prefix: " + os.Getenv("AWS_SECRET_NAME_PREFIX"))	
	log.Infoln("          AdminGroupName: " + s.AdminGroupName)
	
	log.Infoln("          LogLogrusLevel: " + os.Getenv("LOG_LOGRUS"))
	log.Infoln("                 LogGORM: " + os.Getenv("LOG_GORM"))
	log.Infoln("                 LogHTTP: " + os.Getenv("LOG_HTTPLOG"))
	log.Infoln("             UsedBackend: " + s.UsedBackend)
	
	// Postgres credentials
	if s.UsedBackend == "postgres" {
		log.Infoln("           POSTGRES_USER: " + os.Getenv("POSTGRES_USER"))
		log.Infoln("           POSTGRES_PASS: " + s.hideSecretIfReq(os.Getenv("POSTGRES_PASS")))
		log.Infoln("         POSTGRES_DBNAME: " + os.Getenv("POSTGRES_DBNAME"))
		log.Infoln("           POSTGRES_HOST: " + os.Getenv("POSTGRES_HOST"))
		log.Infoln("           POSTGRES_PORT: " + os.Getenv("POSTGRES_PORT"))
	}

	// Postgres credentials
	if s.UsedBackend == "mysql" {
		log.Infoln("              MYSQL_USER: " + os.Getenv("MYSQL_USER"))
		log.Infoln("              MYSQL_PASS: " + s.hideSecretIfReq(os.Getenv("MYSQL_PASS")))
		log.Infoln("            MYSQL_DBNAME: " + os.Getenv("MYSQL_DBNAME"))
		log.Infoln("              MYSQL_HOST: " + os.Getenv("MYSQL_HOST"))
		log.Infoln("              MYSQL_PORT: " + os.Getenv("MYSQL_PORT"))
	}

	// SQL connection options
	log.Infoln("         SQLMaxIdleConns: " + fmt.Sprintf("%d", s.SQLMaxIdleConns))
	log.Infoln("         SQLMaxOpenConns: " + fmt.Sprintf("%d", s.SQLMaxOpenConns))
	log.Infoln("          SQLMaxLifetime: " + fmt.Sprintf("%d hours", s.SQLMaxLifetime))
}

//
// NewSetup create a new configuration
//
func NewSetupValueSet(input []byte) (*SetupValueSet, error) {
	s := &SetupValueSet{}
	err := s.load(input)
	if err != nil {
		return nil, err
	}

	return s, err
}

//
// load gets config from file, env vars (if found) and command line (if used)
//
func (s *SetupValueSet) load(input []byte) error {
	s.initDefaultValues()
	err := s.loadFromYaml(input)
	if err != nil {
		return err
	}

	err = s.setEnvValues()
	if err != nil {
		return fmt.Errorf("Invalid config: %s", err)
	}

	if s.AWSUseSecretStore {
		region := os.Getenv("AWS_REGION")
		var ss awssm.AWSSecretStorage
		if region != "" {
			ss = awssm.AWSSecretStorage{Region: region}
		}
		err = s.loadGraphAuthValuesFromSecretStorage(&ss)
		if err != nil {
			return fmt.Errorf("Invalid config: Error while loading config from secret storage: %s", err)
		}
	} else {
		err = s.loadGraphAuthValuesFromEnv()
		if err != nil {
			return fmt.Errorf("Invalid config: Error while loading config from environment: %s", err)
		}		
	}

	if s.UsedBackend == "" {
		return fmt.Errorf("No DB backend configured")
	}

	return nil
}

func (s *SetupValueSet) loadGraphAuthValuesFromEnv() error {
	var val string
	
	val = os.Getenv("MSAD_TENANT_ID")
	if val != "" {
		s.TenantId = val
	}

	val = os.Getenv("MSAD_CLIENT_ID")
	if val != "" {
		s.ClientId = val
	}

	val = os.Getenv("MSAD_CLIENT_SECRET")
	if val != "" {
		s.ClientSecret = val
	}

	val = os.Getenv("MSAD_AUTHORITY")
	if val != "" {
		s.Authority = val
	}	

	val = os.Getenv("MSAD_SCOPES")
	if val != "" {
		s.Scopes = []string{val}
	}

	return nil
}

func (s *SetupValueSet) loadGraphAuthValuesFromSecretStorage(ss *awssm.AWSSecretStorage) error {
	var tenantId, clientId, clientSecret, adminGroupName string
	
	if val, err := ss.GetSecret("MSAD_TENANT_ID_SEC"); err != nil {
		return err
	} else {
		// get value from returned {"key":"value"} format
		var data map[string]string
		err := json.Unmarshal([]byte(val), &data)
		if err != nil {
			return err
		}
		tenantId = data["MSAD_TENANT_ID_SEC"]
	}

	if val, err := ss.GetSecret("MSAD_CLIENT_ID"); err != nil {
		return err
	} else {
		var data map[string]string
		err := json.Unmarshal([]byte(val), &data)
		if err != nil {
			return err
		}
		clientId = data["MSAD_CLIENT_ID"]
	}

	if val, err := ss.GetSecret("MSAD_CLIENT_SECRET"); err != nil {
		return err
	} else {
		var data map[string]string
		err := json.Unmarshal([]byte(val), &data)
		if err != nil {
			return err
		}
		clientSecret = data["MSAD_CLIENT_SECRET"]
	}

	if val, err := ss.GetSecret("MSAD_ADMIN_GROUP_NAME"); err != nil {
		return err
	} else {
		var data map[string]string
		err := json.Unmarshal([]byte(val), &data)
		if err != nil {
			return err
		}
		adminGroupName = data["MSAD_ADMIN_GROUP_NAME"]
	}
	
	s.TenantId = tenantId
	s.ClientSecret = clientSecret
	s.ClientId = clientId
	s.AdminGroupName = adminGroupName
	
	return nil
}

//
// setInitConfig initializes the config with initial profile which must be not empty
func (s *SetupValueSet) initDefaultValues() {
	s.LogLogrusLevel = log.InfoLevel
	s.ServerIPAddress = DEFAULT_IP_ADDRESS
	s.ServerPort = DEFAULT_PORT
	s.AdminGroupName = DEFAULT_ADMIN_GROUP_NAME
	s.UseGroupNamePattern = DEFAULT_USE_GROUP_NAME_PATTERN
	s.SQLMaxIdleConns = DEFAULT_SQL_MAX_IDLE_CONNS
	s.SQLMaxOpenConns = DEFAULT_SQL_MAX_OPEN_CONNS
	s.SQLMaxLifetime = time.Hour * DEFAULT_SQL_MAX_LIFETIME
}

//
// initWithEnvValues uses env vars to set specific confi flags
//
func (s *SetupValueSet) setEnvValues() error {
	var val string

	val = os.Getenv("HTTP_PORT")
	if val != "" {
		s.ServerPort = val
	}

	val = os.Getenv("HTTP_ADDRESS")
	if val != "" {
		s.ServerIPAddress = val
	}

	val = os.Getenv("MSAD_ADMIN_GROUP_NAME")
	if val != "" {
		s.AdminGroupName = val
	}

	val = os.Getenv("MSAD_USE_GROUP_NAME_PATTERN")
	if val != "" {
		if val == "True" {
			s.UseGroupNamePattern = true
		} else if val == "False" {
			s.UseGroupNamePattern = false
		} else {
			return fmt.Errorf("Invalid value MSAD_USE_GROUP_NAME_PATTERN: %s, must be: False, True", val)
		}
	}

	val = os.Getenv("LOG_LOGRUS")
	if val != "" {
		switch val {
		case "Trace":
			s.LogLogrusLevel = log.TraceLevel
		case "Debug":
			s.LogLogrusLevel = log.DebugLevel
		case "Info":
			s.LogLogrusLevel = log.InfoLevel
		case "Warn":
			s.LogLogrusLevel = log.WarnLevel
		case "Error":
			s.LogLogrusLevel = log.ErrorLevel
		case "Fatal":
			s.LogLogrusLevel = log.FatalLevel
		case "Panic":
			s.LogLogrusLevel = log.PanicLevel
		default:
			return fmt.Errorf("Invalid logrus log level: %s, must be: Trace, Debug, Info, Warn, Error, Fatal, Panic", val)
		}
	}

	val = os.Getenv("LOG_GORM")
	if val != "" {
		switch val {
		case "True":
			s.LogGORM = true
		case "False":
			s.LogGORM = false
		default:
			return fmt.Errorf("Invalid gorm log flag: %s, must be: False, True", val)
		}
	}

	val = os.Getenv("LOG_HTTPLOG")
	if val != "" {
		switch val {
		case "True":
			s.LogHTTP = true
		case "False":
			s.LogHTTP = false
		default:
			return fmt.Errorf("Invalid httplog log flag: %s, must be: False, True", val)
		}
	}

	var err error

	val = os.Getenv("SQL_MAX_IDLE_CONNS")
	if val != "" {
		s.SQLMaxIdleConns, err = strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("Invalid env variable %s value: %s", "SQL_MAX_IDLE_CONNS", val)
		}
	}

	val = os.Getenv("SQL_MAX_OPEN_CONNS")
	if val != "" {
		s.SQLMaxOpenConns, err = strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("Invalid env variable %s value: %s", "SQL_MAX_OPEN_CONNS", val)
		}
	}

	val = os.Getenv("SQL_MAX_LIFETIME")
	if val != "" {
		var valint int
		valint, err = strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("Invalid env variable %s value: %s", "SQL_MAX_LIFETIME", val)
		}
		s.SQLMaxLifetime = time.Hour * time.Duration(valint)
	}

	val = os.Getenv("AWS_USE_SECRET_STORE")
	if val != "" {
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("Invalid env variable %s value: %s", "AWS_USE_SECRET_STORE", val)
		}
		s.AWSUseSecretStore = boolVal
	}

	val = os.Getenv("AWS_SECRET_NAME_PREFIX")
	if val != "" {
		s.SecretNamePrefix = val
	}
	
	return nil
}

//
// loadFromYamlFile loads the config.yaml file overriding default config
//
func (s *SetupValueSet) loadFromYaml(input []byte) error {
	var doc Document
	err := yaml.Unmarshal(input, &doc)
	if err != nil {
		return err
	}

	for _, logger := range doc.Loggers {
		s.setEnvVars(logger.Kind, logger.Env)
	}

	for _, provider := range doc.Providers {
		s.setEnvVars(provider.Kind, provider.Env)
	}

	for _, server := range doc.Servers {
		s.setEnvVars(server.Kind, server.Env)
	}

	for _, sqloption := range doc.SQLOptions {
		s.setEnvVars(sqloption.Kind, sqloption.Env)
	}

	for _, backend := range doc.Backends {
		s.setEnvVars(backend.Kind, backend.Env)
		s.UsedBackend = backend.Kind
	}

	return nil
}

//
// setEnvVars overrides the default values from config with the env
//
func (s *SetupValueSet) setEnvVars(kind string, env map[string]string) {
	s.Backends = append(s.Backends, kind)
	for key, envval := range env {
		envvar := fmt.Sprintf("%s_%s", strings.ToUpper(kind), strings.ToUpper(key))
		if os.Getenv(envvar) == "" {
			flgval := checkCmdLineUsage(strings.ToLower(key))
			if flgval == "" {
				os.Setenv(envvar, envval)
			} else {
				os.Setenv(envvar, flgval)
			}
		}
	}
}

//
// checkFlagsUsage looks for possible use of an option in the command line
//
func checkCmdLineUsage(flg string) (val string) {
	//flag.StringVar(&val, flg, "", "")
	//flag.Parse()
	return
}

//
// hideSecretIfReq reveals secrets in debug or trace mode
//
func (s *SetupValueSet) hideSecretIfReq(str string) string {
	if s.LogLogrusLevel == log.TraceLevel ||
		s.LogLogrusLevel == log.DebugLevel {
		return str
	}

	return "..."
}
