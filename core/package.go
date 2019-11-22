package core

import "io/ioutil"

import "encoding/json"

type PackageJson struct {
	Main            string          `json:"main"`
	Version         string          `json:"version"`
	Name            string          `json:"name"`
	DevDependencies DevDependencies `json:"devDependencies"`
	Scripts         Scripts         `json:"scripts"`
}
type DevDependencies struct {
	Electron string `json:"electron"`
}
type Scripts struct {
	Start string `json:"start"`
}

func CreatePackageJson(old string) []byte {
	f := ReadPackageJson(old)

	f.Scripts.Start = "electron ."

	data, _ := json.Marshal(f)

	return data
}

func ReadPackageJson(old string) PackageJson {
	file, _ := ioutil.ReadFile(old)

	var f PackageJson
	json.Unmarshal(file, &f)
	return f
}
