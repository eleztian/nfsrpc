package nfsrpc

import (
	"sync"
	"errors"
	"strings"
	"unicode/utf8"
	"unicode"
)

/*
	RFC 5531
*/

// ProcedureID uniquely identifies a remote procedure
type ProcedureID struct {
	ProgramNumber   uint32
	ProgramVersion  uint32
	ProcedureNumber uint32
}

// Procedure represents a ProcedureID and name pair.
type Procedure struct {
	ID   ProcedureID
	Name string
}

// pMap is looked up in ServerCodec to map ProcedureID to method name.
// rMap is looked up in ClientCodec to map method name to ProcedureID.
var procedureRegistry = struct {
	sync.RWMutex
	pMap map[ProcedureID]string
	rMap map[string]ProcedureID
}{
	pMap: make(map[ProcedureID]string),
	rMap: make(map[string]ProcedureID),
}

func isValidProcedureName(procedureName string) bool {
	procedureTypeName := strings.Split(procedureName, ".")
	if len(procedureTypeName) != 2 {
		return false
	}

	for _, name := range procedureTypeName {
		if !isExported(name) {
			return false
		}
	}

	return true
}

// 首字母大写
func isExported(name string) bool {
	firstRune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(firstRune)
}

func RegisterProcedure(procedure Procedure) error {
	if !isValidProcedureName(procedure.Name) {
		return errors.New("invalid procedure name")
	}
	procedureRegistry.Lock()
	defer procedureRegistry.Unlock()
	procedureRegistry.pMap[procedure.ID] = procedure.Name
	procedureRegistry.rMap[procedure.Name] = procedure.ID
	return nil
}

func GetProcedureName(procedureID ProcedureID) (string, bool) {
	procedureRegistry.RLock()
	defer procedureRegistry.RUnlock()
	name, ok := procedureRegistry.pMap[procedureID]
	return name, ok
}

func GetProcedureID(name string) (ProcedureID, bool) {
	procedureRegistry.RLock()
	defer procedureRegistry.RWMutex.RUnlock()
	id, ok := procedureRegistry.rMap[name]
	return id,ok
}

func RemoveProcedureID(procedure interface{}) {
	procedureRegistry.Lock()
	defer procedureRegistry.Unlock()
	switch p := procedure.(type) {
	case string:
		procedureID, ok := procedureRegistry.rMap[p]
		if ok {
			delete(procedureRegistry.pMap, procedureID)
			delete(procedureRegistry.rMap, p)
		}
	case ProcedureID:
		procedureName, ok := procedureRegistry.pMap[p]
		if ok {
			delete(procedureRegistry.pMap, p)
			delete(procedureRegistry.rMap, procedureName)
		}
	}
}