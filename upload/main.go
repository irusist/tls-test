package main

import (
	"bytes"
	"fmt"
	"github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"path/filepath"
)

func main() {
	wrapper.SetCtx(
		// 插件名称
		"upload-plugin",
		// 为解析插件配置，设置自定义函数
		wrapper.ParseConfigBy(parseConfig),
		// 为处理请求头，设置自定义函数
		wrapper.ProcessRequestHeadersBy(onHttpRequestHeaders),
		wrapper.ProcessRequestBodyBy(onHttpRequestBody),
	)
}

// 自定义插件配置
type MyConfig struct {
	// 用于发起HTTP调用client
	client wrapper.HttpClient
}

// 在控制台插件配置中填写的yaml配置会自动转换为json，此处直接从json这个参数里解析配置即可
func parseConfig(json gjson.Result, config *MyConfig, log wrapper.Log) error {
	log.Debugf("upload-plugin parseConfig start...")

	return nil
}

func onHttpRequestHeaders(ctx wrapper.HttpContext, config MyConfig, log wrapper.Log) types.Action {
	log.Debugf("upload-plugin onHttpRequestHeaders start...")

	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
	}

	for _, h := range hs {
		proxywasm.LogWarnf("request header <-- %s: %s", h[0], h[1])
	}

	contentType, _ := proxywasm.GetHttpRequestHeader("content-type")
	ctx.SetContext("contentType", contentType)

	return types.ActionContinue
}

func onHttpRequestBody(context wrapper.HttpContext, config MyConfig, body []byte, log wrapper.Log) types.Action {
	log.Debugf("upload-plugin onHttpRequestBody start...")

	contextType := context.GetContext("contentType").(string)
	parse(contextType, body)
	return types.ActionContinue

}

func parseMultipartFormData(data []byte, boundary string) ([]*multipart.Part, error) {
	reader := multipart.NewReader(bytes.NewReader(data), boundary)
	var parts []*multipart.Part

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			// 所有部分已读取完毕
			break
		}
		if err != nil {
			fmt.Println("Error reading part:", err)
			return nil, err
		}
		parts = append(parts, part)
	}

	return parts, nil
}

func parse(contentType string, data []byte) {
	// 示例请求体数据和边界（通常你会从HTTP请求的Content-Type头部获取这个边界）
	// 此处替换为实际的multipart/form-data编码的请求体
	//data := []byte("...")
	//contentType := "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW"

	// 解析Content-Type头部，获取boundary参数值
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		fmt.Println("Error parsing media type:", err)
		panic(err)
	}
	boundary, ok := params["boundary"]
	if !ok {
		panic("未找到boundary参数")
	}

	parts, err := parseMultipartFormData(data, boundary)
	if err != nil {
		panic(err)
	}

	for _, part := range parts {
		defer part.Close()

		fmt.Printf("部分字段名称: %s\n", part.FormName())
		// 如果有文件，通常可以找到文件名
		if part.FileName() != "" {
			// 此时可以通过part.Header中获取文件相关的其它信息
			fmt.Printf("File name: %s\n", part.FileName())
			fmt.Printf("File name ext: %s\n", filepath.Ext(part.FileName()))
		}

		header := make(textproto.MIMEHeader)
		for key, values := range part.Header {
			header[textproto.CanonicalMIMEHeaderKey(key)] = values
			fmt.Printf("%s---%s: %s\n", key, textproto.CanonicalMIMEHeaderKey(key), values)
		}

		//partBytes, err := ioutil.ReadAll(part)
		//if err != nil {
		//	panic(err)
		//}
		//
		//fmt.Printf("部分内容: %s\n\n", string(partBytes))

		// 此处可以处理文件保存或其他逻辑， 文件内容需要根据 part.Header 的 Content-Type 来处理？
		// 处理文件扩展名安全控制
		// 处理文件大小控制   暂时使用 请求头的 Content-Length 来确定
	}
}
