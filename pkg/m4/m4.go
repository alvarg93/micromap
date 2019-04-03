package m4

import (
	"os/exec"
)

func Merge(inputFile, outputFile string) error {
	_, err := exec.Command("sh", "-c", "m4 "+inputFile+" > "+outputFile).Output()
	return err
}
