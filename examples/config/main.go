package main

import (
	"fmt"
	"time"

	"github.com/shipperizer/kilo-franz/config"
)

func main() {
	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)

	fmt.Printf("%+v\n", cfg)
	fmt.Printf("refresh timeout: %s\n", cfg.GetRefreshTimeout())
	fmt.Printf("logger: %v\n", cfg.GetLogger())
}
