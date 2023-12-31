package main

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	DatabaseDriver           string `env:"DATABASE_DRIVER" required:"true"`
	DatabaseConnectionString string `env:"DATABASE_CONNECTIONSTRING" required:"true"`

	OIDCDiscovery    url.URL `env:"OIDC_DISCOVERY" required:"true"`
	OIDCClientID     string  `env:"OIDC_CLIENTID" required:"true"`
	OIDCClientSecret string  `env:"OIDC_CLIENTSECRET" required:"true"`
	OIDCRedirectURL  url.URL `env:"OIDC_REDIRECTURL" required:"true"`

	CORSOrigin string `env:"CORS_ORIGIN"`

	CSRFKey string `env:"CSRF_KEY" required:"true" default:"changeme"`

	SessionStoreKey string `env:"SESSION_STORE_KEY" required:"true" default:"changeme"`

	ListenPort int `env:"PORT" required:"true" default:"8080"`
}

func ParseConfig() (conf Config, rerr error) {
	conf = Config{}
	rconf := reflect.ValueOf(&conf).Elem()
	tconf := rconf.Type()

	for i := 0; i < rconf.NumField(); i++ {
		f := rconf.Field(i)
		t := tconf.Field(i).Tag
		if e := t.Get("env"); e != "" {
			reqtag := t.Get("required")
			var req bool = false
			var err error
			if reqtag != "" {
				req, err = strconv.ParseBool(reqtag)
				if err != nil {
					rerr = fmt.Errorf("Parsing 'required' struct tag for field %s: %w", tconf.Field(i).Name, err)
					return
				}
			}
			val := os.Getenv(e)
			if val == "" && req {
				val = t.Get("default")
				if val == "" {
					rerr = fmt.Errorf("Missing configuration value %s", tconf.Field(i).Name)
					return
				}
			}
			switch f.Type().Name() {
			case "string":
				f.SetString(val)
				break
			case "int":
				iv, err := strconv.ParseInt(val, 0, 0)
				if err != nil {
					rerr = fmt.Errorf("Can't parse envar %s value '%s' as int: %w", tconf.Field(i).Name, val, err)
					return
				}
				f.SetInt(iv)
				break
			case "URL":
				uv, err := url.Parse(val)
				if err != nil {
					rerr = fmt.Errorf("Can't parse envar %s value '%s' as URL: %w", tconf.Field(i).Name, val, err)
					return
				}
				f.Set(reflect.ValueOf(*uv))
				break
			default:
				rerr = fmt.Errorf("Unexpected config struct type %s on field %s", f.Type().Name(), tconf.Field(i).Name)
				return
			}
		}
	}
	return
}
