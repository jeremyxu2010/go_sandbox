package main
import (
	"path/filepath"
	"os/exec"
	"os"
	"log"
)
// 此命令工具用于将hexo部署至服务器
func main() {
	blogDir := filepath.Join(`W:\`, `gits`, `blog`)
	deployBlog(blogDir)
}
func deployBlog(blogDir string) {
	cmd := exec.Command(`hexo`, `deploy`)
	cmd.Dir = blogDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
