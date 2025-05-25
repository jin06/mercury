package admin

import "context"

func Run(ctx context.Context) error {
	// Initialize the admin server
	adminServer := &adminServer{
		ctx: ctx,
	}

	// Start the admin server
	if err := adminServer.start(); err != nil {
		return err
	}

	// Wait for the context to be done
	<-ctx.Done()

	// Stop the admin server gracefully
	return adminServer.stop()
}
