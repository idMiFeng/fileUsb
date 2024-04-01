package utils

import (
	"context"
	"log"
	"os/exec"
	"strings"
	"time"
)

// 封装的cmd函数，目录，名字，参数，返回字节数据

func RunCommandWithOutput(timeout uint64, dir string, args ...string) ([]byte, error) {
	cmdStr := strings.Join(args, " ")
	var cmd *exec.Cmd
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	if dir != "" {
		cmd.Dir = dir
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("执行错误: %v, 详细输出: %s\n", err, string(out))
	}

	return out, err
}
