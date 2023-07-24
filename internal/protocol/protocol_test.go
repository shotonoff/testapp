package protocol

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	ctx := context.Background()
	srv, err := NewServer("0.0.0.0:5001")
	require.NoError(t, err)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		err = srv.Serve(ctx, echoHandler())
		require.NoError(t, err)
	}()
	client := NewClient("localhost:5001")
	err = client.Connect()
	require.NoError(t, err)
	echo := Echo{}
	err = client.Read(ctx, &echo)
	require.NoError(t, err)
	require.Equal(t, []byte("echo"), echo.Data)
	err = client.Disconnect()
	require.NoError(t, err)
	srv.Stop()
	wg.Wait()
}

func echoHandler() HandlerFunc {
	return func(ctx context.Context, tr Connection) error {
		echo := Echo{Data: []byte("echo")}
		data, err := Marshal(echo)
		if err != nil {
			return err
		}
		return tr.Write(ctx, data)
	}
}
