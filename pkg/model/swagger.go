package model

//SwaggerModel swagger文件 结构
type SwaggerModel struct {
	SwaggerVersion string       `json:"swagger"`
	Paths          SwaggerPaths `json:"paths"`
	Definitions    definitions  `json:"definitions"`
	Responses      responses    `json:"responses"`
}

//SwaggerPaths paths
type SwaggerPaths struct {
}

type definitions struct {
}

type responses struct {
}
