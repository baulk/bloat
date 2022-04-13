package main

import (
	"fmt"
	"os"
)

func printf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "format [%s] --> %s\n", format, s)
}

func main() {
	printf("[%20d]", -9999)
	printf("[%-20d]", -9999)
	printf("[%020d]", -9999)
	printf("[%-020d]", -9999)
	printf("[%20X]", -99991111)
	printf("[%-20x]", -99998888)
	printf("[%020X]", -999900000)
	printf("[%2X]", -999900000)
	printf("[%20s]", "abcdefgf")
	printf("[%-20s]", "abcdefgf")
	printf("[%020s]", "abcdefgf")
	printf("[%-020s]", "abcdefgf")
	printf("[%0-20v]", "abcdefgf")
	printf("[%-20v]", "abcdefgf")
	printf("[%20d]", -999900000)
	printf("[%2X]", -999900000)
	printf("[%020v]", true)
	printf("[%20v]", true)
	printf("[%-020v]", true)
	printf("[%020.4f]", 1992.85)
	printf("[%00020.4f]", 1992.85)
	printf("[%20.8f]", 1993.85)
	printf("[%020.8f]", -3.141592654)
	printf("[%v]", -3.141592654)
}
