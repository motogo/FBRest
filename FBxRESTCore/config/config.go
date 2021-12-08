package config

import ("database/sql"
		"strconv"
		_sessions "softend.de/fbrest/FBxRESTCore/sessions"
		_"github.com/nakagami/firebirdsql"
		log "github.com/sirupsen/logrus"		
		_dbscheme "softend.de/fbrest/FBxRESTCore/dbscheme"
)

const (
	DefaultUser string = "SYSDBA"
	DefaultPassword string = "masterkey"
	DefaultPort int = 3050
	DefaultLocation = "localhost"
)

func MakeConnectionStringFromStruct(datas _dbscheme.DatabaseAttributes) (cmd string) {
	
	cmd = MakeConnectionString(datas.Port,datas.Location, datas.Database , datas.User , datas.Password)	
	return 
}

func MakeConnectionString(port int,location string, filename string, user string, password string) (cmd string) {
	if(port < 1){
	    port = DefaultPort
	}
	if(len(location)) < 1{
	    location = DefaultLocation
	}

	if(len(user) < 1) {
		user = DefaultUser
	}

	if(len(password) < 1) {
		password = DefaultPassword
	}

	cmd = string(user+":"+password+"@"+location+":"+strconv.Itoa(port)+"/"+filename)	
	return
}

func ConnLocation(port int,location string, filename string, user string, password string) (db *sql.DB, err error) {
   
    var connstr = MakeConnectionString(port,location, filename , user , password)	
	return ConnLocationWithString(connstr)
}

func ConnLocationWithString(connectionstring string) (db *sql.DB, err error) {
   
	log.WithFields(log.Fields{"Open database:": connectionstring,	}).Info("func ConnLocationWithString")    
	db, err = sql.Open("firebirdsql", connectionstring) 
	if err != nil {		
		log.WithFields(log.Fields{"Open database, Error": err.Error(),	}).Error("func ConnLocationWithString")
	}
	return
}
func ConnLocationWithSession(kv _sessions.Item) (db *sql.DB, err error) {
   
	log.WithFields(log.Fields{"Open database:": kv.Token,	}).Info("func ConnLocationWithSession")    
	db, err = sql.Open("firebirdsql", kv.Value) 
	if err != nil {		
		log.WithFields(log.Fields{"Open database, Error": err.Error(),	}).Error("func ConnLocationWithSession")
	}
	return
}

func TestConnLocation(connstr string) (err error) {
    
	log.WithFields(log.Fields{"Open database, connection:": connstr,	}).Info("func TestConnLocation")    
	
	var db *sql.DB
	db, err = sql.Open("firebirdsql", connstr) 

	if err != nil {
		return
	}

	err = db.Ping(); 
	return
}




