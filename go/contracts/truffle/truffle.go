package truffle

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Deploy() (string, error) {
	cryptoRoot, err := exec.Command("/usr/bin/git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	log.Infof("Repo root at: %s", cryptoRoot)
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if err = os.Chdir(string(cryptoRoot[:len(cryptoRoot)-1]) + "/truffle-core"); err != nil {
		return "", err
	}
	cmd := exec.Command("./node_modules/.bin/truffle", "migrate", "--network", "local", "--reset")
	log.Infof("Deploy by truffle: %v", cmd.String())
	o, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// log.Infof("%s", o)
	output := string(o)
	startIdx := strings.Index(output, "-----------AddressesStart-----------") + len("-----------AddressesStart-----------")
	endIdx := strings.Index(output, "-----------AddressesEnd-----------")
	var contracts map[string]string
	contracts = make(map[string]string)
	addresses := strings.Split(output[startIdx+1:endIdx-1], "\n")
	for _, address := range addresses {
		log.Info(address)
		pair := strings.SplitN(address, "=", 2)
		contracts[pair[0]] = pair[1]
	}
	if err = os.Chdir(workingDir); err != nil {
		return "", err
	}
	return contracts["DIRECTORY_CONTRACT_ADDRESS"], nil
}
