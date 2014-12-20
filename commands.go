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
	Usage: "UF check <root_dir>",
	Description: `
	input target dir.
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
type fileInfo struct{
	path string
	size int64
	md5	[16]byte
	isSame bool
}

type fileInfoList []fileInfo

// for Sort function of sort package
func (p fileInfoList) Len() int {
    return len(p)
}

// for Sort function of sort package
func (p fileInfoList) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

// for Sort function of sort package
func (p fileInfoList) Less(i, j int) bool {
	if p[i].size != p[j].size{
	    return p[i].size < p[j].size
	}else{
		return p[i].path > p[j].path
	}
}

// set fileInfo to fileInfoList
func (f *fileInfoList)setFile(path string)(err error){
	file, err := os.Open(path)
	defer file.Close()
    if err != nil{
    	assert(err)
    	return err
   	}
    
    fstat, err := file.Stat()
    if err != nil{
    	assert(err)
    	return err
    }
    
    var buff []byte
    _, err = file.Read(buff)
    if err != nil {
    	assert(err)
    	return err
    }
    
    filepack := fileInfo{path, fstat.Size(), md5.Sum(buff), false}
    *f = append(*f, filepack)
    return nil
}


func (f *fileInfoList)setup(root string)(err error){
     err = filepath.Walk(root, 
        func(path string, info os.FileInfo, err error) error {
            if info.IsDir() {
                // 特定のディレクトリ以下を無視する場合は
                // return filepath.SkipDir
                return nil
            }
            
           	rel, err := filepath.Rel(root, path)
           	if err != nil{
           		assert(err)
           		return err
           	}
            f.setFile(root + rel)
            return err
        })

    if err != nil {
        assert(err)
        return err
    }
    
    return nil
}

func (f *fileInfoList)isSame(a, b fileInfo)(bool){
	if a.md5 == b.md5 && a.size == b.size{
		return true
	}
	return false
}

func (f *fileInfoList)check()(err error){
	sort.Sort(f)
	
	for i :=0; i<len(*f)-1; i++{
		if f.isSame((*f)[i], (*f)[i+1]){ 
			(*f)[i+1].isSame = true
		}
	} 
	return nil
}

func (f *fileInfoList)delete()(err error){
	remain := make([]fileInfo, 0, 10000)
	willDelete := make([]fileInfo, 0, 10000)
	
	for _, file := range *f{
		if file.isSame{
			willDelete = append(willDelete, file)
		}else{
			remain = append(remain, file)
		}
	}
	
	for _, file := range willDelete{
		log.Println("delete: ", file)
		err := os.Remove(file.path)
		if err != nil{
			assert(err)
		}
	}
	
	(*f) = remain
	
	return nil
}

func doCheck(c *cli.Context) {
	
	var MaxFileNum = 10000
	//list := make([]fileInfo, 0, 10000)
	for _, path := range c.Args(){
		file, err := os.Open(path)
		if err != nil {
			assert(err)
			continue
		}
		file.Close()
		
		var list fileInfoList = make([]fileInfo, 0, MaxFileNum)
		
		err = list.setup(path)
		if err != nil{
			assert(err)
			return
		}
		
		err = list.check()
		if err != nil{
			assert(err)
			return
		}
		
		err = list.delete()
		if err != nil{
			assert(err)
			return
		}
	}
}
