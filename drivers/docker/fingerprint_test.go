package docker

import (
	"context"
	"testing"

	"github.com/hashicorp/nomad/ci"
	"github.com/hashicorp/nomad/client/testutil"
	"github.com/hashicorp/nomad/helper/testlog"
	"github.com/hashicorp/nomad/plugins/drivers"
	"github.com/stretchr/testify/require"
)

// TestDockerDriver_FingerprintHealth asserts that docker reports healthy
// whenever Docker is supported.
//
// In Linux CI and AppVeyor Windows environment, it should be enabled.
func TestDockerDriver_FingerprintHealth(t *testing.T) {
	ci.Parallel(t)
	testutil.DockerCompatible(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := NewDockerDriver(ctx, testlog.HCLogger(t)).(*Driver)

	fp := d.buildFingerprint()
	require.Equal(t, drivers.HealthStateHealthy, fp.Health)
}
