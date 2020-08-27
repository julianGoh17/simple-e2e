package main

import "github.com/julianGoh17/simple-e2e/framework/cmd"

func main() {
	rootCmd := cmd.NewRootCmd()
	cmd.InitRootCmd(rootCmd)
	rootCmd.Execute()
}
