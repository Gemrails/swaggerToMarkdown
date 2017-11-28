package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/akkuman/parseConfig"
)

//SwaggerAction SwaggerAction
type SwaggerAction struct{}

//CreateSWManager CreateManager
func CreateSWManager() *SwaggerAction {
	return &SwaggerAction{}
}

//ShowConf ShowConf
func (s *SwaggerAction) ShowConf(sw string) error {
	if err := checkConf(sw); err != nil {
		logrus.Errorf("config check error, %v", err)
		return err
	}
	config := getConf(sw)
	ls := config.Get("paths")
	for k, v := range ls.(map[string]interface{}) {
		methods := []string{"post", "get", "put", "delete"}
		for _, method := range methods {
			rc := dealMethod(method, k, v, &config)
			fmt.Print(rc)
		}
	}
	return nil
}

type toolsSwaggerParams struct {
	Type        string                 `json:"type"`
	Format      string                 `json:"format"`
	Description string                 `json:"description"`
	Name        string                 `json:"name"`
	In          string                 `json:"in"`
	Required    string                 `json:"required"`
	Schema      map[string]interface{} `json:"schema"`
}

const (
	outparams = "|%s |%s |%s |%s |%s |\n"
	outEG     = "**请求URL:**\n\n- `%s`\n\n**简要描述**\n\n- %s\n\n**请求方式:**\n\n- %s\n\n**参数:**\n\n|参数名|参数位置|必选|类型|说明|\n|:----  |:------|:---|:----- |-----   |\n%s\n**返回示例:**\n\n- %s\n\n* * *\n"
)

func dealMethod(method, path string, msg interface{}, c *parseConfig.Config) string {
	logrus.Debugf("in path %s", path)
	methodKey := fmt.Sprintf("paths > %s > %s", path, method)
	methodV := c.Get(methodKey)
	logrus.Debugf("path 「 %s 」 has no method %s", path, method)
	if methodV == nil {
		//logrus.Warnf("path 「 %s 」 has no method %s", path, method)
		return ""
	}
	//describe, url, method, params\n, return
	desc := fmt.Sprintf("%s, %s", c.Get(methodKey+" > summary"), c.Get(methodKey+" > description"))
	paramsL := c.Get(methodKey + " > parameters")
	params := ""
	if paramsL == nil {
		goto IGNORE
	}
	for _, param := range paramsL.([]interface{}) {
		p := assignment(param.(map[string]interface{}))
		if p.In == "path" {
			paramStr := fmt.Sprintf(outparams, p.Name, "path", p.Required, p.Type, p.Description)
			if params == "" {
				params = paramStr
			} else {
				params += paramStr
			}
		}
		if p.In == "body" {
			pp := dealBodyParams(p.Schema)
			if pp != "" {
				params += pp
			}
		}
	}
IGNORE:
	outStr := fmt.Sprintf(outEG, path, desc, strings.ToUpper(method), params, "return")
	return outStr
}

func dealBodyParams(schema map[string]interface{}) string {
	requiredL, ok := schema["required"]
	if !ok {
		requiredL = []interface{}{}
	}
	if propertiesL, ok := schema["properties"]; ok {
		params := ""
		for k, v := range propertiesL.(map[string]interface{}) {
			required := "False"
			for _, vv := range requiredL.([]interface{}) {
				if vv == k {
					required = "True"
					break
				}
			}
			ty := v.(map[string]interface{})["type"]
			description := v.(map[string]interface{})["description"]
			if description != nil && ty != nil {
				//fmt.Printf("ty is %v\n", ty)
				desc := ""
				if strings.Contains(description.(string), "\nin:") {
					desc = strings.Split(description.(string), "\nin:")[0]
				}
				paramStr := fmt.Sprintf(outparams, k, "body", required, ty.(string), desc)
				if params == "" {
					params = paramStr
				} else {
					params += paramStr
				}
				continue
			}
			paramStr := fmt.Sprintf(outparams, k, "body", required, "Object", "N/A")
			if params == "" {
				params = paramStr
			} else {
				params += paramStr
			}
			continue
		}
		return params
	}
	return ""
}

func assignment(mm map[string]interface{}) *toolsSwaggerParams {
	var tsp toolsSwaggerParams
	if t, ok := mm["type"]; ok {
		tsp.Type = t.(string)
	}
	if desc, ok := mm["description"]; ok {
		tsp.Description = desc.(string)
	} else {
		tsp.Description = "-"
	}
	if n, ok := mm["name"]; ok {
		tsp.Name = n.(string)
	}
	if i, ok := mm["in"]; ok {
		tsp.In = i.(string)
	}
	if r, ok := mm["required"]; ok {
		if r.(bool) {
			tsp.Required = "True"
		} else {
			tsp.Required = "False"
		}
	}
	if s, ok := mm["schema"]; ok {
		tsp.Schema = s.(map[string]interface{})
	}
	return &tsp
}

func checkConf(confPath string) error {
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		return fmt.Errorf("config.json is not exist")
	}
	return nil
}

func getConf(confPath string) parseConfig.Config {
	return parseConfig.New(confPath)
}
