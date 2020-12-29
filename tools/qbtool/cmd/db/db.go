package db

import (
	"fmt"
	"github.com/hero1s/golib/tools/qbtool/cmd/base"
	"github.com/hero1s/golib/tools/qbtool/cmd/db/diff"
	"github.com/hero1s/golib/tools/qbtool/cmd/db/reverse"
)

var Cmd = &base.Command{
	CustomFlags: true,
	Run:         runCmd,
	UsageLine:   base.ToolName + " db command [arguments...]",
	Short:       "根据数据库表结构生成GO Struct代码, 比对两个数据库，生成差异SQL语句",
	Long: `
COMMANDS:
  reverse 由数据库表结构生成Go语言结构体
  diff    对比两个数据库，生成差异SQL语句

ARGUMENTS:
  reverse:
   -source   数据库,例子 username:password@tcp(ip:port)/db_name?charset=ut8
   -path     输出路径
   -single   是否单文件输出
  diff:
   -source  源数据库,例子 username:password@tcp(ip:port)/db_name?charset=ut8
   -target  目标数据库,例子 username:password@tcp(ip:port)/db_name?charset=ut8
   -path    对比文件输出路径
`,
}

var source, target, path *string
var singleFile *bool

func runCmd(cmd *base.Command, args []string) {
	if len(args) >= 2 {
		source = cmd.Flag.String("source", "", "源数据库地址")
		target = cmd.Flag.String("target", "", "源数据库地址")
		path = cmd.Flag.String("path", ".", "输出路径")
		singleFile = cmd.Flag.Bool("single", false, "是否单文件输出")
		cmd.Flag.Parse(args[2:])
	} else {
		fmt.Printf("please see %s help db\n", base.ToolName)
		return
	}
	action := args[1]
	if action == "reverse" {
		if err := reverse.RunReverse(*source, *path,*singleFile); err != nil {
			fmt.Println(err)
		}
	} else if action == "diff" {
		if err := diff.RunDiff(*source, *target, *path); err != nil {
			fmt.Println(err)
		}
	}
}
