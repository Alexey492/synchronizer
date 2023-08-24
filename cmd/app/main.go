package main

import (
	"context"

	"github.com/Alexey492/synchronizer/logic"
)

func main() {
	// адрес директорий
	srcDir := "./output/src"
	dstDir := "./output/dst"

	logger := logic.GetLogger()

	logger.Printf("Syncing %s to %s", srcDir, dstDir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go logic.Monitor(ctx, logger, srcDir, dstDir)

	<-ctx.Done()

	logger.Printf("Syncing finished")
}
