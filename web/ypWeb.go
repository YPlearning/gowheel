package web

import (
    "net/http"
)

type Resp struct {
    Code    string `json:"code"`
    Msg     string `json:"msg"`
}

func start() {
    http.HandleFunc("/hello", handler)
	http.Handle("/", http.FileServer(http.Dir(".")))
    http.ListenAndServe(":8080", nil)
}

/*******************************************************************************
*	Internal Callback Function
*******************************************************************************/
func handler(writer http.ResponseWriter, request *http.Request) {
    writer.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
    writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
    writer.Header().Set("content-type", "application/json") //返回数据格式是json

	request.ParseForm()
	tablename, tError :=  request.Form["tablename"]

	var result  Resp
    if !tError {
        result.Code = "401"
        result.Msg = "查询失败"
    } else {
        result.Code = "200"
        result.Msg = tablename[0]
    }


    writer.Write([]byte(result.Msg))
	
}