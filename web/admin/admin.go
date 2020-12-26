package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hero1s/golib/log"
	"github.com/hero1s/golib/task"
	"github.com/gin-gonic/gin"
	"net/http"
	"text/template"
	"time"
)

// M is Map shortcut
type M map[string]interface{}

var passWord = "123123"

func checkPasswd(c *gin.Context) bool {
	enterPassWord, _ := c.Cookie("admin-passwd")
	if passWord == enterPassWord {
		return true
	}
	enterPassWord = c.Query("passwd")
	if enterPassWord != passWord {
		c.Abort()
		log.Errorf("密码错误:%v", enterPassWord)
		return false
	}
	c.SetCookie("admin-passwd", enterPassWord, 3600, "/", "localhost", false, true)
	return true
}

func Init(g *gin.Engine, pass string) {
	g.Use(addQps())
	// keep in mind that all data should be html escaped to avoid XSS attack
	g.GET("/", adminIndex)
	g.GET("/qps", qpsIndex)
	g.GET("/prof", profIndex)
	g.GET("/task", taskStatus)
	passWord = pass
}

// 统计访问请求
func addQps() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		StatisticsMap.AddStatistics(c.Request.Method, c.Request.URL.Path, "admin", latency)
	}
}

// AdminIndex is the default http.Handler for admin module.
// it matches url pattern "/".
func adminIndex(c *gin.Context) {
	if !checkPasswd(c) {
		return
	}
	writeTemplate(c, map[interface{}]interface{}{}, indexTpl, defaultScriptsTpl)
}

// QpsIndex is the http.Handler for writing qps statistics map result info in http.ResponseWriter.
// it's registered with url pattern "/qps" in admin module.
func qpsIndex(c *gin.Context) {
	if !checkPasswd(c) {
		return
	}
	data := make(map[interface{}]interface{})
	data["Content"] = StatisticsMap.GetMap()

	// do html escape before display path, avoid xss
	if content, ok := (data["Content"]).(M); ok {
		if resultLists, ok := (content["Data"]).([][]string); ok {
			for i := range resultLists {
				if len(resultLists[i]) > 0 {
					resultLists[i][0] = template.HTMLEscapeString(resultLists[i][0])
				}
			}
		}
	}

	writeTemplate(c, data, qpsTpl, defaultScriptsTpl)
}

// ProfIndex is a http.Handler for showing profile command.
// it's in url pattern "/prof" in admin module.
func profIndex(c *gin.Context) {
	if !checkPasswd(c) {
		return
	}
	c.Request.ParseForm()
	command := c.Request.Form.Get("command")
	if command == "" {
		return
	}

	var (
		format = c.Request.Form.Get("format")
		data   = make(map[interface{}]interface{})
		result bytes.Buffer
	)
	ProcessInput(command, &result)
	data["Content"] = template.HTMLEscapeString(result.String())

	if format == "json" && command == "gc summary" {
		dataJSON, err := json.Marshal(data)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(c.Writer, dataJSON)
		return
	}

	data["Title"] = template.HTMLEscapeString(command)
	defaultTpl := defaultScriptsTpl
	if command == "gc summary" {
		defaultTpl = gcAjaxTpl
	}
	writeTemplate(c, data, profillingTpl, defaultTpl)
}

func writeJSON(rw http.ResponseWriter, jsonData []byte) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonData)
}

// TaskStatus is a http.Handler with running task status (task name, status and the last execution).
// it's in "/task" pattern in admin module.
func taskStatus(c *gin.Context) {
	if !checkPasswd(c) {
		return
	}
	data := make(map[interface{}]interface{})

	// Run Task
	c.Request.ParseForm()
	taskname := c.Request.Form.Get("taskname")
	if taskname != "" {
		if t, ok := task.AdminTaskList[taskname]; ok {
			if err := t.Run(); err != nil {
				data["Message"] = []string{"error", template.HTMLEscapeString(fmt.Sprintf("%s", err))}
			} else {
				t.SetPrev(t.GetNext())
				data["Message"] = []string{"success", template.HTMLEscapeString(fmt.Sprintf("%s run success,Now the Status is <br>%s", taskname, t.GetStatus()))}
			}
		} else {
			data["Message"] = []string{"warning", template.HTMLEscapeString(fmt.Sprintf("there's no task which named: %s", taskname))}
		}
	}

	// List Tasks
	content := make(M)
	resultList := new([][]string)
	var fields = []string{
		"Task Name",
		"Task Spec",
		"Task Status",
		"Last Time",
		"",
	}
	for tname, tk := range task.AdminTaskList {
		result := []string{
			template.HTMLEscapeString(tname),
			template.HTMLEscapeString(tk.GetSpec()),
			template.HTMLEscapeString(tk.GetStatus()),
			template.HTMLEscapeString(tk.GetPrev().String()),
		}
		*resultList = append(*resultList, result)
	}

	content["Fields"] = fields
	content["Data"] = resultList
	data["Content"] = content
	data["Title"] = "Tasks"
	writeTemplate(c, data, tasksTpl, defaultScriptsTpl)
}

func writeTemplate(c *gin.Context, data map[interface{}]interface{}, tpls ...string) {
	tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
	for _, tpl := range tpls {
		tmpl = template.Must(tmpl.Parse(tpl))
	}
	tmpl.Execute(c.Writer, data)
}
