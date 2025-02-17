package kube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	mtypes "github.com/akash-network/akash-api/go/node/market/v1beta3"
	"github.com/akash-network/node/sdl"
	"github.com/akash-network/node/testutil"

	"github.com/akash-network/provider/cluster/kube/builder"
)

const (
	randDSeq uint64 = 1
	randGSeq uint32 = 2
	randOSeq uint32 = 3
)

func TestDeploy(t *testing.T) {
	t.Skip()
	ctx := context.Background()

	owner := ed25519.GenPrivKey().PubKey().Address()
	provider := ed25519.GenPrivKey().PubKey().Address()

	leaseID := mtypes.LeaseID{
		Owner:    sdk.AccAddress(owner).String(),
		DSeq:     randDSeq,
		GSeq:     randGSeq,
		OSeq:     randOSeq,
		Provider: sdk.AccAddress(provider).String(),
	}

	sdl, err := sdl.ReadFile("../../_run/kube/deployment.yaml")
	require.NoError(t, err)

	mani, err := sdl.Manifest()
	require.NoError(t, err)

	log := testutil.Logger(t)
	client, err := NewClient(ctx, log, "lease", "")
	require.NoError(t, err)

	ctx = context.WithValue(ctx, builder.SettingsKey, builder.NewDefaultSettings())
	err = client.Deploy(ctx, leaseID, &mani.GetGroups()[0])
	require.NoError(t, err)
}
