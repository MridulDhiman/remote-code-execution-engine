package types

import (
	"path"
	"strings"
)

type LangConfig struct {
	Lang          string
	IsCompiled    bool
	Runner       string
	FileName string
}

const (
	DOCKER_WORKDIR = "/judge/"
	SCRIPT_DIR = "/scripts/"
)


func NewLangConfig(lang string, isCompiled bool, Runner string, FileName string) *LangConfig {
	return &LangConfig{
		Lang:          lang,
		IsCompiled:    isCompiled,
		Runner:       Runner,
		FileName: FileName,
	}
}

func (l* LangConfig) GetImageAndContainerName () (string, string) {
	var imageName string = l.Lang + "-image";
	var contName string = l.Lang + "-cont";

	if(l.Lang == "c" || l.Lang == "cpp") {
		imageName = "gcc-image"
		contName = "gcc-cont"
	}

	return imageName, contName
}



func (l *LangConfig) DockerExecCmd(hasInput bool) string {
	  inputCmd := ""
	  inputCmdForCompiledLang:= ""
	  _,contName := l.GetImageAndContainerName()

	  if hasInput {
		inputCmd = strings.Join([]string{"cat", path.Join(DOCKER_WORKDIR, "input.txt"), "|"}, " ")
		inputCmdForCompiledLang = strings.Join([]string{"&&", strings.Trim(inputCmd, " "), path.Join(DOCKER_WORKDIR, "/code")}, " ")
	  }

	  cmd:= ""

if(!l.IsCompiled) {
	cmd =  strings.Join([]string{inputCmd, l.Runner, path.Join(DOCKER_WORKDIR, l.FileName)}, " ")
} else if l.Lang != "rust" {
	cmd =  strings.Join([]string{l.Runner, path.Join(DOCKER_WORKDIR, l.FileName), "-o code", inputCmdForCompiledLang}, " ")
} 
if l.Lang == "rust" {
	cmd =  strings.Join([]string{l.Runner, path.Join(DOCKER_WORKDIR, l.FileName, inputCmdForCompiledLang)}, " ")
}


return strings.Join([]string{"docker exec -i", contName, "sh -c", "'", strings.Trim(cmd, " "), "'"}, " ")
}



func (l* LangConfig) DockerRunCmd () (string) {
	imageName, contName := l.GetImageAndContainerName()

	if !l.IsCompiled {
 return  strings.Join([]string{"docker run -dit --name", contName, imageName}, " ");
} else {
	return strings.Join([]string{"docker run -d --name", contName, "--privileged", imageName, "tail -f /dev/null"}, " ");
}
}




