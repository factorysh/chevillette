package stdin

import (
	"bufio"
	"io"
	"os"

	"github.com/factorysh/chevillette/input"
)

func New(io.Reader) input.LineScanner {
	return bufio.NewScanner(os.Stdin)

}
