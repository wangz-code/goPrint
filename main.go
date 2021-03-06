package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"net/http"

	"github.com/gin-gonic/gin"
)

type setting struct {
	Customerid string `json:"customerid"`
	Store      string `json:"store"`
	User       string `json:"user"`
	Pwd        string `json:"pwd"`
	Machineno  string `json:"machineno"`
}
type settingData struct {
	TemplateName string  `json:"templateName"`
	Value        setting `json:"value"`
}

func ReShutDownEXE() {
	fmt.Println("重启主机")
	arg := []string{"-r", "-t", "20"}
	cmd := exec.Command("shutdown", arg...)
	d, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(d))
	return
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, IndexHtml())
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/reboot", func(c *gin.Context) {
		ReShutDownEXE()
		c.JSON(200, gin.H{
			"state": "ok",
		})
	})
	r.GET("/config", func(c *gin.Context) {
		//settings := setting{Customerid: "test", User: "张三", Store: "118", Pwd: "123", Machineno: "A111"},
		settingData := []settingData{
			{
				TemplateName: "张三的配置",
				Value:        setting{Customerid: "test", User: "张三", Store: "118", Pwd: "123", Machineno: "A111"}},
			{
				TemplateName: "李四的配置",
				Value:        setting{Customerid: "test", User: "李四", Store: "M01", Pwd: "123", Machineno: "A111"}},
		}
		c.JSON(200, gin.H{
			"message": "pong",
			"state":   "ok",
			"data":    settingData,
		})
	})

	r.POST("/saveConfig", func(c *gin.Context) {
		// 执行保存
		fmt.Print(c)

		json := make(map[string]interface{}) //注意该结构接受的内容
		c.BindJSON(&json)
		fmt.Printf("%v", &json)
		c.JSON(http.StatusOK, gin.H{
			"data":  json["data"],
			"state": "ok",
		})

	})
	// 生成配置 bat 文件
	r.POST("/generateConfig", func(c *gin.Context) {
		fileName := "chromeRun.bat"
		json := make(map[string]string) //注意该结构接受的内容
		c.BindJSON(&json)
		fmt.Printf("%v", &json)
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("json====", json["data"])
		params := "start \"\" chrome.exe --kiosk https://www.baidu.com/#/print/list?p=" + json["data"]
		content := []byte(params)
		err = ioutil.WriteFile(fileName, content, 0777)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("write file successful")
		defer file.Close()
		c.JSON(http.StatusOK, gin.H{
			"data":  json["data"],
			"state": "ok",
		})
	})
	exec.Command(`open`, `http://localhost:8080`).Start()
	r.Run()
}

func IndexHtml() string {

	return `
   <!DOCTYPE html>
   <html lang="en">
	   <head>
		   <meta charset="UTF-8" />
		   <meta http-equiv="X-UA-Compatible" content="IE=edge" />
		   <meta name="viewport" content="width=device-width, initial-scale=1.0" />
   
		   <link rel="shortcut icon" href="https://s4.ax1x.com/2022/01/19/7rtNO1.png" />
		   <link href="https://cdn.bootcdn.net/ajax/libs/spectre.css/0.5.9/spectre.min.css" rel="stylesheet">
		   <title>安装配置</title>
		   <style>
			   body {
				   display: flex;
				   justify-content: center;
			   }
			   .toast {
				   position: absolute;
				   top: 20px;
				   right: calc(50% - 200px);
				   width: 400px;
			   }
			   .bg {
				   padding: 50px 80px;
				   height: 100vh;
				   align-items: center;
			   }
			   .card {
				   width: 800px;
				   border-radius: 8px;
				   box-shadow: 0 3px 6px -4px #0000001f, 0 6px 16px #00000014, 0 9px 28px 8px #0000000d;
				   transition: color 0.3s;
			   }
		   </style>
	   </head>
	   <body>
		   <div id="app">
			   <div class="toast toast-primary" v-if="toastShow">
				   <button class="btn btn-clear float-right"></button>
				   提示信息
			   </div>
			   <div class="bg">
				   <div class="modal modal-sm active" v-if="modalShow">
					   <a class="modal-overlay" aria-label="Close"></a>
					   <div class="modal-container">
						   <div class="modal-header">
							   <a href="#close" class="btn btn-clear float-right" aria-label="Close" @click="()=>modalShow=!modalShow"></a>
							   <div class="modal-title h5">配置名称</div>
						   </div>
						   <div class="modal-body">
							   <div class="content">
								   <div class="form-group">
									   <label class="form-label">配置名称</label>
									   <input class="form-input" v-model="templatName" type="text" placeholder="请输入配置名称" />
								   </div>
							   </div>
						   </div>
						   <div class="modal-footer">
							   <button class="btn btn-primary" @click="saveConfig">确认</button>
						   </div>
					   </div>
				   </div>
   
				   <div class="card">
					   <div class="card-header">
						   <div class="card-title h2">自助打印配置</div>
						   <div class="card-subtitle text-gray">自定义配置</div>
					   </div>
					   <div class="card-body">
						   <form action="">
							   <div class="form-group">
								   <label class="form-label">历史配置</label>
								   <select class="form-select" v-model="selected" @change="selectChange">
									   <option v-for="(item,idx) in settingData" :value="idx">{item.templateName}</option>
								   </select>
							   </div>
   
							   <div class="form-group">
								   <label class="form-label">客户编码</label>
								   <input class="form-input" v-model="setting.customerid" name="customerid" type="text" placeholder="客户识别码" />
							   </div>
							   <div class="form-group">
								   <label class="form-label">商场编码</label>
								   <input class="form-input" v-model="setting.store" name="store" type="text" placeholder="商场编码" />
							   </div>
							   <div class="form-group">
								   <label class="form-label">登录账户</label>
								   <input class="form-input" v-model="setting.user" name="user" type="text" placeholder="用户账户" />
							   </div>
							   <div class="form-group">
								   <label class="form-label">账户密码</label>
								   <input class="form-input" v-model="setting.pwd" name="pwd" type="text" placeholder="账户密码" />
							   </div>
							   <div class="form-group">
								   <label class="form-label">机器编码</label>
								   <input class="form-input" v-model="setting.machineno" name="machineno" type="text" placeholder="机器编码" />
							   </div>
						   </form>
					   </div>
					   <div class="card-footer">
						   <button class="btn btn-primary btn-lg" @click="()=>modalShow=!modalShow">保存当前配置</button>
						   <button class="btn btn-primary btn-lg" @click="generateConfig">生成脚本并设置开机启动</button>
						   <button class="btn btn-primary btn-lg" @click="rebootSys">重启系统</button>
					   </div>
				   </div>
			   </div>
		   </div>
	   </body>
	   <script src="https://cdn.jsdelivr.net/npm/vue@2.6.14/dist/vue.min.js"></script>
	   <script>
	   class XHTTP {
		get(url, callback) {
			var xhr = new XMLHttpRequest();
			xhr.open("GET", url, true);
			xhr.send();
			xhr.onreadystatechange = function () {
				if (xhr.readyState == 4 && xhr.status == 200) {
					var res = JSON.parse(xhr.responseText);
					callback(res);
				}
			};
		}
		post(url, data, callback) {
			var xhr = new XMLHttpRequest();
			xhr.open("POST", url, true);
			xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
			xhr.send(data);
			xhr.onreadystatechange = function () {
				if (xhr.readyState == 4 && xhr.status == 200) {
					var res = JSON.parse(xhr.responseText);
					callback(res);
				}
			};
		}
	}
		   var http = new XHTTP();
		   var app = new Vue({
			   delimiters: ["{", "}"],
			   el: "#app",
			   data: {
				   toastShow: false,
				   modalShow: false,
				   templatName: "",
				   selected: "",
				   settingData: [],
				   setting: {},
			   },
			   mounted() {
				   this.getConfig();
			   },
			   methods: {
				   getConfig() {
					   http.get("/config", (res) => {
						   this.settingData = res.data;
						   this.selected = 0;
						   this.setting = res.data[0].value;
					   });
				   },
				   // 重启
				   rebootSys() {
					   http.get("/reboot");
				   },
				   // 配置下拉变更
				   selectChange() {
					   this.setting = this.settingData[this.selected].value;
				   },
				   // 保存配置信息
				   saveConfig() {
					   http.post("/saveConfig", JSON.stringify({ data: this.setting }), (res) => {
						   console.log("res log==>", res);
						   if (res.state == "ok") {
							   this.toastShow = true;
							   this.modalShow = false;
							   setTimeout(() => {
								   this.toastShow = false;
							   }, 2000);
						   }
					   });
				   },
				   // 生成配置信息
				   generateConfig() {
					   http.post("/generateConfig", JSON.stringify({ data: encodeURI(JSON.stringify(this.setting)) }), (res) => {
						   console.log("res log==>", res);
						   if (res.state == "ok") {
							   this.toastShow = true;
							   setTimeout(() => {
								   this.toastShow = false;
							   }, 2000);
						   }
					   });
				   },
			   },
		   });
   
		   
	   </script>
   </html>
   
	
	
	
	`
}
