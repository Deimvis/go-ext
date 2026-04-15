package ext

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Shred execs `shred -u -z "$filePath"`.
func Shred(filePath string) error {
	cmd := exec.Command("shred", "-u", "-z", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CopyFile(srcFp string, dstFp string) error {
	srcStat, err := os.Stat(srcFp)
	if err != nil {
		return err
	}

	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", srcFp)
	}

	src, err := os.Open(srcFp)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstFp)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}
