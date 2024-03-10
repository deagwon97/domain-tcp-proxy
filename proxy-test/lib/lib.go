package lib

import "os"

func CheckHostsExist(hsotnameWithIpPort string) bool {
	hostFile, err := os.Open("/etc/hosts")
	return true
}