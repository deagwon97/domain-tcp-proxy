package lib

import (
	"fmt"
	"os"
	"strings"
)

func AddHostEntry(host, ip string) error {
	// fmt.Println("Adding new entry to /etc/hosts")
	entry := fmt.Sprintf("%s\t%s", ip, host)

	// Check if the entry already exists in the hosts file
	hostsFile, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return err
	}

	if strings.Contains(string(hostsFile), entry) {
		return fmt.Errorf("entry already exists in /etc/hosts")
	}

	// Append the new entry to the hosts file
	file, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// fmt.Println("Adding new entry to /etc/hosts")
	_, err = fmt.Fprintln(file, entry)
	if err != nil {
		return err
	}

	return nil
}
