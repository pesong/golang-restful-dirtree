package main

import (
	"os"
	"os/user"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"encoding/json"	

	"github.com/gorilla/mux"     // go get   (mux router)
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", DefaultRouter)
	router.HandleFunc("/files", FilesDirRouter)     

	log.Fatal(http.ListenAndServe(":8090", router))
}


//different router mapping 
func DefaultRouter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome! This is my restful api program! \n")
}

func FilesDirRouter(w http.ResponseWriter, r *http.Request) {
/*  //for windows:  
  usrdir := "C:\\mygo\\src\\rest"*/
	 
//get user's home dir 
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }

//call ParseDirTree
    tree,err := ParseDirTree(usr.HomeDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error traversing the filesystem: %v\n", err)
		os.Exit(3)
	}else{
		if err := json.NewEncoder(w).Encode(tree); err != nil {
		panic(err)
	}
	}




}



//def struct DirTree
type DirTree struct {
	IsDir 		bool		`json:"IsDir"`
	Name 		string 		`json:"name"`
	Path 		string    	`json:"path"`
	Children 	[]*DirTree  `json:"children"`
}



// build the directory tree 
func ParseDirTree(root string) (result *DirTree, err error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return
	}
	parents := make(map[string]*DirTree)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		parents[path] = &DirTree{
			IsDir: info.IsDir(),
			Name: info.Name(),
			Path: path,
			Children: make([]*DirTree, 0),
		}
		return nil
	}
	if err = filepath.Walk(absRoot, walkFunc); err != nil {
		return
	}
	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { 
			result = node
		} else {
			parent.Children = append(parent.Children, node)
		}
	}
	return
}



// convert struct to json 
func (parsed *DirTree) ToJson() string {
	j, err := json.Marshal(parsed)
	if err!=nil {
		log.Println("JSON ERROR: " + err.Error())
		return "JSON ERROR: " + err.Error()
	}
	return string(j)
}
