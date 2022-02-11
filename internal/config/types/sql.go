package types

import (
	"database/sql"
	"strings"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
)

const (
	sslDisabledMode = "disable"
)

type SQL struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	Schema          string
	SSL             *SSL
	MaxOpenConns    uint32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration

	//Additional options to be appended as options=<Options>cd
	//The value will be taken as is. So be sure to separate multiple options by a space
	Options string
}

type SQLBase struct {
	Host            string
	Port            string
	Database        string
	Schema          string
	SSL             SSLBase
	MaxOpenConns    uint32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration

	//Additional options to be appended as options=<Options>
	//The value will be taken as is. So be sure to separate multiple options by a space
	Options string
}

type SQLBase2 struct {
	Host            string
	Port            string
	Database        string
	Schema          string
	MaxOpenConns    uint32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	SSL             SSLBase

	//Additional options to be appended as options=<Options>
	//The value will be taken as is. So be sure to separate multiple options by a space
	Options string
}

type SQLUser struct {
	User            string
	Password        string
	SSL             SSLUser
	MaxOpenConns    uint32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type SSL struct {
	// type of connection security
	Mode string
	// RootCert Path to the CA certificate
	RootCert string
	// Cert Path to the client certificate
	Cert string
	// Key Path to the client private key
	Key string
}

type SSLBase struct {
	// type of connection security
	Mode string
	// RootCert Path to the CA certificate
	RootCert string
}

type SSLUser struct {
	// Cert Path to the client certificate
	Cert string
	// Key Path to the client private key
	Key string
}

func (s *SQL) connectionString() string {
	fields := []string{
		"host=" + s.Host,
		"port=" + s.Port,
		"user=" + s.User,
		"dbname=" + s.Database,
		"application_name=zitadel",
		"sslmode=" + s.SSL.Mode,
	}
	if s.Options != "" {
		fields = append(fields, "options="+s.Options)
	}
	if s.Password != "" {
		fields = append(fields, "password="+s.Password)
	}
	s.checkSSL()
	if s.SSL.Mode != sslDisabledMode {
		fields = append(fields, "sslrootcert="+s.SSL.RootCert)
		if s.SSL.Cert != "" {
			fields = append(fields, "sslcert="+s.SSL.Cert)
		}
		if s.SSL.Cert != "" {
			fields = append(fields, "sslkey="+s.SSL.Key)
		}
	}

	return strings.Join(fields, " ")
}

func (s *SQL) Start() (*sql.DB, error) {
	client, err := sql.Open("postgres", s.connectionString())
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "TYPES-9qBtr", "unable to open database connection")
	}
	// as we open many sql clients we set the max
	// open cons deep. now 3(maxconn) * 8(clients) = max 24 conns per pod
	client.SetMaxOpenConns(int(s.MaxOpenConns))
	client.SetConnMaxLifetime(s.MaxConnLifetime)
	client.SetConnMaxIdleTime(s.MaxConnIdleTime)

	return client, nil
}

func (s *SQL) checkSSL() {
	if s.SSL == nil || s.SSL.Mode == sslDisabledMode || s.SSL.Mode == "" {
		s.SSL = &SSL{Mode: sslDisabledMode}
		return
	}
	if s.SSL.RootCert == "" {
		logging.LogWithFields("TYPES-LFdzP",
			"cert set", s.SSL.Cert != "",
			"key set", s.SSL.Key != "",
			"rootCert set", s.SSL.RootCert != "",
		).Fatal("fields for secure connection missing")
	}
}

func (u SQLUser) Start(base SQLBase) (*sql.DB, error) {
	return (&SQL{
		Host:     base.Host,
		Port:     base.Port,
		User:     u.User,
		Password: u.Password,
		Database: base.Database,
		Options:  base.Options,
		SSL: &SSL{
			Mode:     base.SSL.Mode,
			RootCert: base.SSL.RootCert,
			Cert:     u.SSL.Cert,
			Key:      u.SSL.Key,
		},
	}).Start()
}
