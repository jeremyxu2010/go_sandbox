package main
import (
	"log"
	"os/exec"
	"time"
	"os"
	"path/filepath"
)
// 此命令工具用于辅助hexo进行博客写作
// 原理：
// 1. 根据source目录下的文件提交变化自动重新生成博客静态文件
// 2. 启动hexo服务器，端口为5000
func main() {
	blogDir := filepath.Join(`W:\`, `gits`, `blog`)
	//首先生成一次，然后取得当前的commitID
	generateBlog(blogDir)
	lastCommitID := getLastCommitID(blogDir)
	//运行hexo本地服务器
	go func() {
		runBlog(blogDir)
	}()
	//不停地检查当前commitID是否与保存的是否一致，如不一致，则重新生成
	for {
		newCommitID := getLastCommitID(blogDir)
		if newCommitID != lastCommitID {
			generateBlog(blogDir)
			lastCommitID = newCommitID
		}
		time.Sleep(10 * time.Second)
	}
}
func runBlog(blogDir string) {
	cmd := exec.Command(`hexo`, `server`, `-s`, `-p`, `5000`)
	cmd.Dir = blogDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func getLastCommitID(blogDir string) string {
	cmd := exec.Command(`git`, `rev-parse`, `HEAD`)
	cmd.Dir = filepath.Join(blogDir, `source`)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}
func generateBlog(blogDir string) {
	cmd := exec.Command(`hexo`, `generate`)
	cmd.Dir = blogDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}