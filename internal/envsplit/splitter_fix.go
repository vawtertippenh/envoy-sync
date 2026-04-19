package envsplit

import (
	"fmt"
	"strings"
)

// ensure imports used across files are satisfied — real logic is in splitter.go
var _ = fmt.Errorf
var _ = strings.HasPrefix
