package core

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qor/validations"
	"github.com/sirupsen/logrus"
)

type (
	// Server stores gin.Engine, Database and Log for side-wide usage
	Server struct {
		Engine *gin.Engine
		db     *gorm.DB
		Log    *logrus.Logger

		jwtSigningKey string

		// set this to true if you want to call (*gorm.DB).Debug() on
		// every request
		databaseDebug bool
	}
)

const (
	defaultCookieSecret = "cookie*secret-nWS37AzEYActW4X"
	defaultSessionName  = "ldb/core-session"
)

var server *Server

// ServerFromConfigFile create a runable server from config file path.
func ServerFromConfigFile(path string) (*Server, error) {
	config, err := NewConfigFromJSONFile("config.json")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"module": "example",
			"action": "main",
		}).Panicf("Create server failed. Could not load config file, %s", err.Error())
	}
	return ServerFromConfig(config)
}

// ServerFromConfig create a runable server from config object.
func ServerFromConfig(conf Config) (*Server, error) {
	server = &Server{
		Log:           logrus.New(),
		jwtSigningKey: "jwt*signing&key+nD5gUktrSQnSyxq#",
		databaseDebug: false,
	}
	logFields := logFieldsForMethod("EngineFromConfig")
	mode, err := conf.StringValueForKey("application.mode")
	if err != nil || mode != "release" {
		server.Log.WithFields(logFields).Warn("[1/3: Init Engine] application.mode is not 'release' or missing, Debug Mode is use by default.")
		server.Engine = gin.Default()
	} else {
		server.Log.WithFields(logFields).Info("[1/3: Init Engine] Server will run in release mode.")
		server.Engine = gin.New()
	}

	dbProfile, _ := conf.StringValueForKey("application.db-profile")
	dbConfig, err := conf.ConfigValueForKey("database." + dbProfile)
	err = server.loadDatabase(dbConfig)
	if err != nil {
		server.Log.WithFields(logFields).Panicf("[2/3: Connect database] Failed, %", err.Error())
		return &Server{}, err
	}
	server.Log.WithFields(logFields).Info("[2/3: Connect database] Success")

	databaseMigration(server.db)
	server.Log.WithFields(logFields).Info("[3/3: Core Models migaration] Success")

	// hook gorm callback for validation
	validations.RegisterCallbacks(server.db)
	// init default cookie store
	server.UseCookieStore(defaultCookieSecret, defaultSessionName)
	// serve static in ./static with /static path by default
	server.Use(static.Serve("/static", static.LocalFile("./static", true)))

	initRoutes(server.Engine)
	return server, nil
}

// Run is warper of (*gin.Engine).Run(port)
func (server *Server) Run(port string) {
	server.Engine.Run(port)
}

// Use is wrapper for (*gin.Engine).Use(...)
func (server *Server) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return server.Engine.Use(middleware...)
}

// Group is shortcut for (*gin.Engine).Group(...)
func (server *Server) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return server.Engine.Group(relativePath, handlers...)
}

// GET is shortcut for (*gin.Engine).GET(string, ...gin.HandlerFunc)
func (server *Server) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return server.Engine.GET(relativePath, handlers...)
}

// POST is shortcut for (*gin.Engine).POST(string, ...gin.HandlerFunc)
func (server *Server) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return server.Engine.POST(relativePath, handlers...)
}

// UseCookieStore if you want to name your own cookie & session
func (server *Server) UseCookieStore(secret string, name string) {
	store := sessions.NewCookieStore([]byte(secret))
	server.Use(sessions.Sessions(name, store))
}

// ServeStatic serving static resources. A single variadic bool is accepted
// for file's indexes.
func (server *Server) ServeStatic(path string, filePath string, indexes ...bool) {
	if len(indexes) > 0 {
		server.Use(static.Serve(path, static.LocalFile(filePath, indexes[0])))
		return
	}
	server.Use(static.Serve(path, static.LocalFile(filePath, true)))
}

func (server *Server) loadDatabase(conf Config) error {
	var conn DatabaseConnection
	var err error
	// Any error here will lead to error on opening connection,
	// so just check it at one place.
	dialect, _ := conf.StringValueForKey("dialect")
	databaseName, _ := conf.StringValueForKey("database")
	username, _ := conf.StringValueForKey("username")
	password, _ := conf.StringValueForKey("password")
	host, _ := conf.StringValueForKey("host")
	port, _ := conf.StringValueForKey("port")

	conn, err = NewDatabaseConnection(dialect, databaseName, username, password, host, port)
	if err != nil {
		return err
	}
	server.db, err = gorm.Open(dialect, conn.ConnectionString())
	if err != nil {
		return err
	}
	return nil
}

// SetJWTSigningKey if you want to use a custom signing key
func (server *Server) SetJWTSigningKey(key string) {
	server.jwtSigningKey = key
}

// SetDBDebug set active to true if you want to call (*gorm.DB).Debug()
// on every query.
func (server *Server) SetDBDebug(active bool) {
	server.databaseDebug = active
}

// DB return (*gorm.DB) object with Debug() method,
// set by SetDatabaseDebug(bool)
func (server *Server) DB() *gorm.DB {
	if server.databaseDebug {
		return server.db.Debug()
	}
	return server.db
}
