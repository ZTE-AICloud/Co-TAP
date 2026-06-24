package main

import (
	"log"
	"os"
	"syscall"
)

var (
	DirToCheck = []string{"../uapregistry-works/logs"}
	ErrLogFile = "../uapregistry-works/logs/stderr.log"
)

func checkDirs(d string) {
	if _, err := os.Stat(d); err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(d, os.FileMode(int(0750))); err != nil {
			log.Fatalf("Failed to mkdir:%s", d)
		}
	}
}

// @title           AI 注册中心&知识图谱
// @version         1.0
// @description     提供AI注册中心功能，以及知识图谱相关功能.
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @tag.name node
// @tag.description 知识图谱节点(CRUD操作, 批量创建操作)
// @tag.name relationship
// @tag.description 知识图谱关系(CRUD操作, 批量创建操作)
// @tag.name graph
// @tag.description 知识图谱整体，节点与关系。支持导出与导入全量图数据
func main() {
	for _, d := range DirToCheck {
		checkDirs(d)
	}
	f, err := os.OpenFile(ErrLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(int(0640)))
	if err != nil {
		log.Fatalf("Failed to open file to log stderr:%v", err)
	}
	err = syscall.Dup3(int(f.Fd()), int(os.Stderr.Fd()), 0)
	if err != nil {
		log.Fatalf("Failed to redirect stderr to regular file:%v", err)
	}
	cli := NewCLI(os.Stdout, os.Stderr)
	os.Exit(cli.Run(os.Args))
}
