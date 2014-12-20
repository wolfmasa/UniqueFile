package main

import (
	"log"
	"path/filepath"
	"os"
	"sort"
	"crypto/md5"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandCheck,
}

var commandCheck = cli.Command{
	Name:  "check",
	Usage: "",
	Description: `
`,
	Action: doCheck,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
type filePackage struct{
	path string
	size int64
	md5	[16]byte
}

type fileList []filePackage

func (p fileList) Len() int {
    return len(p)
}

func (p fileList) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

func (p fileList) Less(i, j int) bool {
    return p[i].size < p[j].size
}

func (f *fileList)setFile(path string){
	
            ff, err := os.Open(path)
            if err != nil{
            	log.Println(err)
            }
            if err == nil {
            	fstat, _ := ff.Stat()
            	var buff []byte
            	_, err := ff.Read(buff)
            	if err != nil {
            		log.Println(err)
            	}
            	file := filePackage{path, fstat.Size(), md5.Sum(buff)}
          	  	*f = append(*f, file)
          	  	//log.Println("add path : ", file.path, file.size)
            }
}

func (f *fileList)listupFiles(root string){
     err := filepath.Walk(root, 
        func(path string, info os.FileInfo, err error) error {
            if info.IsDir() {
                // 特定のディレクトリ以下を無視する場合は
                // return filepath.SkipDir
                return nil
            }
            
            rel, err := filepath.Rel(root, path)
            if err != nil{
            	log.Println(err)
            }

            f.setFile(root + rel)

            return nil
        })

    if err != nil {
        log.Println(1, err)
    }
}

func doCheck(c *cli.Context) {
	log.Println(c.Args())
	
	var MaxFileNum = 10000
	//list := make([]filePackage, 0, 10000)
	for _, path := range c.Args(){
		var list fileList = make([]filePackage, 0, MaxFileNum)
		//log.Println(path)
		list.listupFiles(path)
		sort.Sort(list)
		
		for _, l := range list{
			log.Println(l)
		}
	}
}
