package arkeocli

import (
	"bufio"
	"strings"

	"github.com/spf13/cobra"
)

func promptForArg(cmd *cobra.Command, prompt string) (string, error) {
	cmd.Print(prompt)
	reader := bufio.NewReader(cmd.InOrStdin())
	read, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	read = strings.TrimSpace(read)
	return read, nil
}
