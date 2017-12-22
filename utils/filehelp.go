package utils

import (
	"strings"
	"path/filepath"
	"os"
	"io/ioutil"
	"path"
	"compress/gzip"
	"archive/tar"
	"io"
	"github.com/astaxie/beego"
	"os/exec"
	"archive/zip"
	"sync"
	"github.com/pkg/errors"
	"log"
)

var lock *sync.RWMutex

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetParentDir(dirctory string) string {
	return Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func GetAbsDir(path string) string {
	if len(path) == 0 {
		return GetCurrentDir("")
	}
	if err := Exist(path); err == nil {
		return path
	}
	os.MkdirAll(path, os.ModePerm)
	return path
}

func GetCurrentDir(path string) string {
	fileName := ""
	if strings.Contains(path, ".") {
		fileName = path[strings.LastIndex(path, "/"):]
		path = path[:strings.LastIndex(path, "/")]
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	if len(path) > 0 {
		path = strings.Replace(dir, "\\", "/", -1) + path

		if err := Exist(path); err != nil {
			os.MkdirAll(path, os.ModePerm)
		}
		return path + fileName
	} else {
		return strings.Replace(dir, "\\", "/", -1)
	}
}

func ListFiles(params ...string) (files []string, err error) {
	dirPath := ""
	suffix := ""
	if len(params) == 0 {
		return nil, errors.New("params is null")
	}
	if len(params) == 1 {
		dirPath = params[0]
	}
	if len(params) == 2 {
		dirPath = params[0]
		suffix = params[1]
	}
	fileArr, err := ioutil.ReadDir(dirPath)

	if err != nil {
		return nil, err
	}
	for _, file := range fileArr {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if len(suffix) > 0 {
			if strings.HasSuffix(suffix, path.Ext(fileName)) {
				files = append(files, dirPath+"/"+fileName)
			}
			continue
		}
		files = append(files, dirPath+"/"+file.Name())
	}
	return files, err
}

func ListDir(dirPath string) (files []string, err error) {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		files = append(files, dirPath+"/"+fi.Name())
	}
	return files, nil
}

func WalkDir(dirPath, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	err = filepath.Walk(dirPath, func(fileName string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		// if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
		// 	files = append(files, fileName)
		// }
		files = append(files, fileName)
		return nil
	})
	return files, nil
}

func UnGzip(zipPath, destFolderPath string) (absolutePath string, err error) {
	file, err := os.Open(zipPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	gzipRead, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzipRead.Close()
	tarRead := tar.NewReader(gzipRead)
	if err = os.MkdirAll(destFolderPath, os.ModePerm); err != nil {
		return absolutePath, err
	}
	for {
		header, err := tarRead.Next()
		if err == io.EOF || header == nil {
			break
		}
		guid, _ := GetGuid()
		fw, err := os.Create(destFolderPath + guid)
		if err != nil {
			beego.Info("-------err:", err)
			continue
		}
		defer fw.Close()
		absolutePath = fw.Name()
		io.Copy(fw, tarRead)
		fw.Close()

	}
	//os.Remove(zipPath)
	return absolutePath, err
}

func Unzip(zipPath, destPath string) (err error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	beego.Info("--------zipPath:", zipPath, " | destPath:", destPath)
	for _, k := range r.Reader.File {
		fileName := k.Name
		if k.FileInfo().IsDir() {
			os.MkdirAll(k.Name, 0644)
			continue
		}
		if strings.Contains(fileName, "20170212") {
			fileName = strings.Replace(fileName, "20170212/", "", 1)
		}
		beego.Info("---------k:", fileName, " | k:", k)
		file, err := k.Open()
		if err != nil {
			continue
		}
		defer file.Close()
		// if strings.Contains(k.Name, "general") {
		// 	destPath = destPath + "/temp/general/" + strconv.FormatInt(time.Now().Unix(), 10) + "/"
		// } else {
		// 	destPath = destPath + "temp/crash/" + strconv.FormatInt(time.Now().Unix(), 10) + "/"
		// }
		tempPath := destPath
		if strings.Contains(fileName, "log_unexception") {
			tempPath = tempPath + "unexception/"
			if strings.Contains(fileName, ".info") {
				tempPath = tempPath + "info/"
			}
		} else if strings.Contains(fileName, "log_exception") {
			tempPath = tempPath + "exception/"
			if strings.Contains(fileName, ".info") {
				tempPath = tempPath + "info/"
			}
		} else if strings.Contains(fileName, "log_normal") {
			tempPath = tempPath + "normal/"
			if strings.Contains(fileName, ".info") {
				tempPath = tempPath + "info/"
			}
		}
		if err := os.MkdirAll(destPath, os.ModePerm); err == nil {
			guid, _ := GetGuid()
			f, err := os.Create(destPath + guid)
			beego.Info("----------destPath:", tempPath, " | app_name:", fileName, " | f:", f)
			tempPath = f.Name()
			if err != nil {
				beego.Error("------1-err:", err)
				continue
			}
			beego.Info("--------f:", f, " | file:", file)
			io.Copy(f, file)
			f.Close()
		} else {
			return err
		}

	}
	return err
}

func WriterFile(sourceFile, destFile string, isDeleteSourceFile bool) (errs error) {
	lock = new(sync.RWMutex)
	lock.Lock()
	file, err := os.Open(sourceFile)
	defer file.Close()
	if err != nil {
		lock.Unlock()
		return err
	}
	fWrite, err := os.OpenFile(destFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	defer fWrite.Close()
	if err != nil {
		lock.Unlock()
		return err
	}
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			//panic(err)
			errs = err
			continue
		}
		if n == 0 {
			break
		}

		if n2, err := fWrite.Write(buf[:n]); err != nil {
			//panic(err)
			errs = err
			continue
		} else if n2 != n {
			panic("error in writing")
		}
	}
	if isDeleteSourceFile {
		err = os.Remove(sourceFile)
	}
	lock.Unlock()
	return errs
}

func FileSize(path string) (int64, string) {
	sufix := "default"
	if strings.Contains(path, ".") {
		sufix = strings.Split(path, ".")[1]
	}
	if !strings.Contains(path, GetCurrentDir("")) {
		path = GetCurrentDir(path)
	}
	if fileInfo, err := os.Stat(path); err == nil {
		return fileInfo.Size(), sufix
	} else {
		beego.Error(err)
	}
	return 0, sufix
}

func AbsFileSize(path string) (int64, string) {
	sufix := "default"
	if strings.Contains(path, ".") {
		sufix = strings.Split(path, ".")[1]
	}
	if fileInfo, err := os.Stat(path); err == nil {
		return fileInfo.Size(), sufix
	}
	return 0, sufix
}

/**
追加文件内容
*/
func AppendStrFile(content, filePath string) (n int, err error) {
	fRead, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer fRead.Close()
	buf := []byte(content)
	return fRead.Write(buf)
}

/**
Copyfile
*/
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func ReadAll(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func Exist(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else {
		return err
	}
}

func TarGz(srcDirPath, destFilePath string) error {
	fw, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer fw.Close()
	gWrite := gzip.NewWriter(fw)
	defer gWrite.Close()
	tWrite := tar.NewWriter(gWrite)
	defer tWrite.Close()

	file, err := os.Open(srcDirPath)
	if err != nil {
		return err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		tarGzDir(srcDirPath, path.Base(srcDirPath), tWrite)
	} else {
		tarGzFile(srcDirPath, fileInfo.Name(), tWrite, fileInfo)
	}
	return nil
}

/**
Deal with directories
if find files, handle them with tarGzFile
Every recurrence append the base path to the recPath
recPath is the path inside of tar.gz
*/
func tarGzDir(srcDirPath, recPath string, tw *tar.Writer) (err error) {
	dir, err := os.Open(srcDirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	// Get file info slice
	fileInfos, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, fileInfo := range fileInfos {
		curPath := srcDirPath + "/" + fileInfo.Name()
		if fileInfo.IsDir() {
			tarGzDir(curPath, recPath+"/"+fileInfo.Name(), tw)
		} else {
			beego.Info("Adding file...%s\n", curPath)
		}
		err = tarGzFile(curPath, recPath+"/"+fileInfo.Name(), tw, fileInfo)
	}
	return err
}

/**
Deal with files
*/
func tarGzFile(srcFile, recPath string, tw *tar.Writer, fi os.FileInfo) (err error) {
	if fi.IsDir() {
		// Create tar header
		hdr := new(tar.Header)
		// if last character of header app_name is '/' it also can be directory
		// but if you don't set Typeflag, error will occur when you untargz
		hdr.Name = recPath + "/"
		hdr.Typeflag = tar.TypeDir
		hdr.Size = 0
		//hdr.Mode = 0755 | c_ISDIR
		hdr.Mode = int64(fi.Mode())
		hdr.ModTime = fi.ModTime()

		// Write hander
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}
	} else {
		// File reader
		fr, err := os.Open(srcFile)
		if err != nil {
			return err
		}
		defer fr.Close()

		// Create tar header
		hdr := new(tar.Header)
		hdr.Name = recPath
		hdr.Size = fi.Size()
		hdr.Mode = int64(fi.Mode())
		hdr.ModTime = fi.ModTime()

		// Write hander
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}

		// Write file data
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}
	}
	return err
}

/**
Ungzip and untar from source file to destination directory
you need check file exist before you call this function
*/
func UnTarGz(srcFilePath, destDirPath string) {
	beego.Info("UnTarGzing " + srcFilePath + "...")
	// Create destination directory
	os.Mkdir(destDirPath, os.ModePerm)

	fr, err := os.Open(srcFilePath)
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	// Gzip reader
	gr, err := gzip.NewReader(fr)
	if err != nil {
		panic(err)
	}
	defer gr.Close()

	// Tar reader
	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		//handleError(err)
		beego.Info("UnTarGzing file..." + hdr.Name)
		// Check if it is diretory or file
		if hdr.Typeflag != tar.TypeDir {
			// Get files from archive
			// Create diretory before create file
			os.MkdirAll(destDirPath+"/"+path.Dir(hdr.Name), os.ModePerm)
			// Write data to file
			fw, _ := os.Create(destDirPath + "/" + hdr.Name)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				panic(err)
			}
		}
	}
	beego.Info("Well done!")
}

func ExecCommand(command string) ([]byte, error) {
	res, err := exec.Command("/bin/sh", "-c", ``+command+``).Output()
	return res, err
}
