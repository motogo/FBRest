package models

import (
	"database/sql"

	"encoding/json"
	_struct "softend.de/fbrest/FBxRESTCore/struct"
	"fmt"
	"null"
	"runtime"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type ModelGetData struct {
	DB *sql.DB
}

type Db_init struct {
	DB *sql.DB
}

//func bToKb(b uint64) uint64 {
//return b / 1024
//}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func RString(b string, l int) string {
	var ft string = "%-" + strconv.Itoa(l) + "s"
	return fmt.Sprintf(ft, b)
}

func outMem(inf1, funcstr string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Debug(" ")
	//log.WithFields(log.Fields{inf1: m.Sys,	})        .Debug(funcstr+"->used total bytes system memory")
	log.WithFields(log.Fields{fmt.Sprintf("%-10s", inf1): bToMb(m.Sys)}).Debug(RString(funcstr+"->allocated total system memory", 60))
	log.WithFields(log.Fields{fmt.Sprintf("%-10s", inf1): bToMb(m.TotalAlloc)}).Debug(RString(funcstr+"->used total heap memory", 60))
	log.WithFields(log.Fields{fmt.Sprintf("%-10s", inf1): bToMb(m.Alloc)}).Debug(RString(funcstr+"->allocated local heap memory", 60))
	log.WithFields(log.Fields{fmt.Sprintf("%-10s", inf1): m.NumGC}).Debug(RString(funcstr+"->allocated GC memory", 60))
}

func (model ModelGetData) GetSQLData(cmd string) (getStruct string, err error) {
	const funcstr = "func models.GetSQLData"
	log.Debug(" ")
	log.WithFields(log.Fields{"SQL command": cmd}).Debug(funcstr)
	log.Debug(" ")
	row, err := model.DB.Query(cmd)
	if err != nil {
		log.WithFields(log.Fields{"Error": err.Error()}).Error(funcstr + "->query command")
		//return  nil,err
		return "", err
	} else {

		start := time.Now()
		runtime.GC()
		outMem("Memory before Query (Mb)", funcstr)

		colNames, _ := row.Columns()
		ct, err := row.ColumnTypes()

		if err != nil {
			log.WithFields(log.Fields{"Error": err.Error()}).Error(funcstr + "->getting colnames")
		}

		readCols := make([]interface{}, len(colNames))
		writeCols := make([]null.String, len(colNames))
		//for i, _ := range writeCols {
		for i := range writeCols {
			readCols[i] = &writeCols[i]
		}

		var ergstr string = "["
		var n int = 0
		for row.Next() {
			_isiStruct2 := make([]_struct.SqlResponseData, len(colNames))
			err := row.Scan(readCols...)
			if err != nil {
				log.WithFields(log.Fields{"Error": err.Error()}).Error(funcstr + "->scan next row")
			} else {
				if n == 0 {
					ergstr = ergstr + "{"
				} else {
					ergstr = ergstr + ",{"
				}
				n++
				var sep string
				for i := range ct {
					_isiStruct2[i].Colvalue = writeCols[i].String
					_isiStruct2[i].Colname = ct[i].Name()
					_isiStruct2[i].Datatype = ct[i].DatabaseTypeName()
					_isiStruct2[i].IsNull = !writeCols[i].Valid

					if _isiStruct2[i].Datatype == "VARYING" {
						sep = "\""
					} else if _isiStruct2[i].Datatype == "DATE" {
						sep = "\""
					} else if _isiStruct2[i].Datatype == "TIMESTAMP" {
						sep = "\""
					} else if _isiStruct2[i].Datatype == "TIME" {
						sep = "\""
					} else {
						sep = ""
					}
					if i == 0 {
						ergstr = ergstr + "\"" + _isiStruct2[i].Colname + "\":" + sep + _isiStruct2[i].Colvalue + sep
					} else {
						ergstr = ergstr + ",\"" + _isiStruct2[i].Colname + "\":" + sep + _isiStruct2[i].Colvalue + sep
					}
				}
				ergstr = ergstr + "}"
			}
		}
		ergstr = ergstr + "]"
		log.Debug(" ")
		log.WithFields(log.Fields{"Numbers of rows": strconv.Itoa(n)}).Debug(RString(funcstr+"->got data complete", 60))

		outMem("Memory after Query (Mb)", funcstr)
		runtime.GC()
		outMem("Memory after GC    (Mb)", funcstr)
		log.Debug(" ")
		elapsed := time.Since(start)
		log.WithFields(log.Fields{"Time used": elapsed}).Debug(RString(funcstr+"->used time", 60))
		return ergstr, nil
	}
}

func (model ModelGetData) GetSQLDataOld(cmd string) (getStruct []string, err error) {
	const funcstr = "func models.GetSQLData"
	log.WithFields(log.Fields{"Getting data": cmd}).Debug(funcstr)

	row, err := model.DB.Query(cmd)
	if err != nil {
		log.WithFields(log.Fields{"Error": err.Error()}).Error(funcstr + "->query command")
		return nil, err
	} else {

		colNames, err := row.Columns()

		if err != nil {
			log.WithFields(log.Fields{"Error": err.Error()}).Error(funcstr + "->get colnames")
		}

		readCols := make([]interface{}, len(colNames))
		writeCols := make([]null.String, len(colNames))
		var _isiStruct []string

		for i := range writeCols {
			readCols[i] = &writeCols[i]
		}

		for row.Next() {
			err := row.Scan(readCols...)
			if err != nil {
				log.WithFields(log.Fields{"Error": err.Error()}).Error(funcstr + "->scan next row")
			} else {
				pagesJson, err := json.Marshal(writeCols)
				if err != nil {
					log.WithFields(log.Fields{"Error": err.Error()}).Error(funcstr + "->marshal to JSON")
				}
				_isiStruct = append(_isiStruct, string(pagesJson))
			}
		}
		log.WithFields(log.Fields{"Getting data done": cmd}).Debug(funcstr)
		return _isiStruct, nil
	}
}

func (model ModelGetData) RunSQLData(cmd string) (getStruct []string, err error) {
	const funcstr = "func models.RunSQLData"
	log.WithFields(log.Fields{"Getting data": cmd}).Debug(funcstr)

	row, err := model.DB.Query(cmd)
	if err != nil {
		log.WithFields(log.Fields{"Error": err.Error()}).Error(funcstr + "->query command")
		return nil, err
	} else {

		colNames, err := row.Columns()

		if err != nil {
			log.WithFields(log.Fields{"Error": err.Error()}).Error(funcstr + "->get colnames")
		}

		readCols := make([]interface{}, len(colNames))
		writeCols := make([]null.String, len(colNames))
		var _isiStruct []string

		for i := range writeCols {
			readCols[i] = &writeCols[i]
		}

		for row.Next() {
			err := row.Scan(readCols...)
			if err != nil {
				log.WithFields(log.Fields{"Error": err.Error()}).Error(funcstr + "->scan next row")
			} else {
				pagesJson, err := json.Marshal(writeCols)
				if err != nil {
					log.WithFields(log.Fields{"Error": err.Error()}).Error(funcstr + "->marshal to JSON")
				}
				_isiStruct = append(_isiStruct, string(pagesJson))
			}
		}
		log.WithFields(log.Fields{"Getting data done": cmd}).Debug(funcstr)
		return _isiStruct, nil
	}
}
