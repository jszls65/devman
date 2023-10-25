package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

// func main() {
func main22() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./src/persistence/models",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	gormdb, err := gorm.Open(sqlite.Open("../dev-utils.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("连接数据库失败: ", err)
		return
	}
	g.UseDB(gormdb)

	//g.GenerateModelAs("alert_log", "AlertLog")
	g.GenerateModelAs("alert_job", "AlertJob")
	// g.GenerateModelAs("interface_call_log", "InterfaceCallLog")
	// g.GenerateModelAs("request_log", "RequestLog")
	//g.GenerateModelAs("request_summary", "RequestSummary")
	g.Execute()
}
