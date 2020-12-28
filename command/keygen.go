package command

import (
	"flag"
	"fmt"

	"github.com/mouuff/go-rocket-update/crypto"
)

type Keygen struct {
	flagSet *flag.FlagSet

	keyName string
}

func (cmd *Keygen) Name() string {
	return "keygen"
}

func (cmd *Keygen) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)

	cmd.flagSet.StringVar(&cmd.keyName, "name", "id_rsa", "name of the key to generate")

	return cmd.flagSet.Parse(args)
}

func (cmd *Keygen) Run() error {
	priv, err := crypto.GeneratePrivateKey()
	if err != nil {
		return err
	}

	fmt.Println("rsa: ", crypto.ExportPrivateKey(priv), "!")
	fmt.Println("name: ", cmd.keyName, "!")

	return nil
}
