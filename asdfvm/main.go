// main is the entry point
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// ModuleArgs holds the arguments for the asdf-vm module
type ModuleArgs struct {
	Name  string `json:"name"`
	URL   string `json:"url,omitempty"`
	State string `json:"state,omitempty"`
  Version string `json:"version,omitempty"`
  Default *bool `json:"default"`
}


//Response is what ansible expects as module output
type Response struct {
	Msg     string `json:"msg"`
	Changed bool   `json:"changed"`
	Failed  bool   `json:"failed"`
}

//FailJSON returns a failed response
func FailJSON(responseBody Response) {
	responseBody.Failed = true
	returnResponse(responseBody)
}

func returnResponse(responseBody Response) {
	var response []byte
	var err error
	response, err = json.Marshal(responseBody)
	if err != nil {
		response, _ = json.Marshal(Response{Msg: "Invalid response object"})
	}
	fmt.Println(string(response))
	if responseBody.Failed {
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	var response Response

	if len(os.Args) != 2 {
		response.Msg = "No argument file provided"
		FailJSON(response)
	}

	argsFile := os.Args[1]

	text, err := os.ReadFile(argsFile)
	if err != nil {
		response.Msg = "Could not read configuration file: " + argsFile
		FailJSON(response)
	}

	var moduleArgs ModuleArgs
	err = json.Unmarshal(text, &moduleArgs)
	if err != nil {
		response.Msg = "Configuration file not valid JSON: " + argsFile
		FailJSON(response)
	}

  if moduleArgs.Name == "" {
    response.Msg = "'name' needs to be set."
    FailJSON(response)
  }
  if moduleArgs.State == ""  {
    moduleArgs.State = "present"
  }
  
  if moduleArgs.Version == "" {
    moduleArgs.Version = "latest"
  }

  if moduleArgs.Default == nil {
    setDefault := true 
    moduleArgs.Default = &setDefault
  }

  changed, err := EnsureAsdfPlugin(moduleArgs.Name, moduleArgs.URL, moduleArgs.State, moduleArgs.Version, *moduleArgs.Default)
  response.Changed = changed
  if err != nil {
    response.Failed = true
    response.Msg = fmt.Sprintf("Error running asdf: %v", err)
  } else {
    response.Failed = false
    response.Msg = fmt.Sprintf("Success running module")
  }
  
  returnResponse(response)
}
