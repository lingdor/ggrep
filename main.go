package main

import "flag"

func main() {

	var cmd command
	cmd.ready()
	flag.Parse()
	cmd.run()

}
