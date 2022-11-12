package disk

import "fmt"

func validateDisks(disks int) error {
	if disks < 20 || disks > 240 {
		return fmt.Errorf("VM validatation failed: ram size must be between 1024 and 65536 MB ")
	}
	return nil
}