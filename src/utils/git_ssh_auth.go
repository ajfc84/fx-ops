package utils

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/transport"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)

func GitSSHAuth(user, keyPath string) (transport.AuthMethod, error) {
	signer, err := gitssh.NewPublicKeysFromFile(user, keyPath, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load SSH key: %w", err)
	}
	signer.HostKeyCallbackHelper = gitssh.HostKeyCallbackHelper{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return signer, nil
}

func GitSSHAuthAgent(user string) (transport.AuthMethod, error) {
	auth, err := gitssh.NewSSHAgentAuth(user)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH agent: %w", err)
	}
	return auth, nil
}
