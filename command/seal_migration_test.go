// +build !enterprise

package command

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/shamir"

	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	physInmem "github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func TestSealMigrationAutoToShamir(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace).Named(t.Name())
	phys, err := physInmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	haPhys, err := physInmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	autoSeal := vault.NewAutoSeal(seal.NewTestSeal(nil))
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		Seal:            autoSeal,
		Physical:        phys,
		HAPhysical:      haPhys.(physical.HABackend),
		DisableSealWrap: true,
	}, &vault.TestClusterOptions{
		Logger:      logger,
		HandlerFunc: vaulthttp.Handler,
		SkipInit:    true,
		NumCores:    1,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	initResp, err := client.Sys().Init(&api.InitRequest{
		RecoveryShares:    5,
		RecoveryThreshold: 3,
	})
	if err != nil {
		t.Fatal(err)
	}

	testhelpers.WaitForActiveNode(t, cluster)

	keys := initResp.RecoveryKeysB64
	rootToken := initResp.RootToken
	core := cluster.Cores[0].Core

	client.SetToken(rootToken)
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	shamirSeal := vault.NewDefaultSeal(&seal.Access{
		Wrapper: aeadwrapper.NewWrapper(&wrapping.WrapperOptions{
			Logger: logger.Named("shamir"),
		}),
	})
	shamirSeal.SetCore(core)

	if err := adjustCoreForSealMigration(logger, core, shamirSeal, autoSeal); err != nil {
		t.Fatal(err)
	}

	var resp *api.SealStatusResponse
	unsealOpts := &api.UnsealOpts{}
	for _, key := range keys {
		unsealOpts.Key = key
		unsealOpts.Migrate = false
		resp, err = client.Sys().UnsealWithOptions(unsealOpts)
		if err == nil {
			t.Fatal("expected error due to lack of migrate parameter")
		}
		unsealOpts.Migrate = true
		resp, err = client.Sys().UnsealWithOptions(unsealOpts)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected response")
		}
		if !resp.Sealed {
			break
		}
	}
	if resp.Sealed {
		t.Fatalf("expected unsealed state; got %#v", *resp)
	}

	// Seal and unseal again to verify that things are working fine
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}
	unsealOpts.Migrate = false
	for _, key := range keys {
		unsealOpts.Key = key
		resp, err = client.Sys().UnsealWithOptions(unsealOpts)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected response")
		}
		if !resp.Sealed {
			break
		}
	}
	if resp.Sealed {
		t.Fatalf("expected unsealed state; got %#v", *resp)
	}
}

func TestSealMigration(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace).Named(t.Name())
	phys, err := physInmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	haPhys, err := physInmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	wrapper := vault.NewDefaultSeal(&seal.Access{
		Wrapper: aeadwrapper.NewWrapper(&wrapping.WrapperOptions{
			Logger: logger.Named("shamir"),
		}),
	})
	coreConfig := &vault.CoreConfig{
		Seal:            wrapper,
		Physical:        phys,
		HAPhysical:      haPhys.(physical.HABackend),
		DisableSealWrap: true,
	}
	clusterConfig := &vault.TestClusterOptions{
		Logger:      logger,
		HandlerFunc: vaulthttp.Handler,
		SkipInit:    true,
		NumCores:    1,
	}

	ctx := context.Background()
	var keys []string
	var rootToken string

	{
		logger.Info("integ: start up as normal with shamir seal, init it")
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		client := cluster.Cores[0].Client
		coreConfig = cluster.Cores[0].CoreConfig

		// Init
		resp, err := client.Sys().Init(&api.InitRequest{
			SecretShares:    2,
			SecretThreshold: 2,
		})
		if err != nil {
			t.Fatal(err)
		}
		keys = resp.KeysB64
		rootToken = resp.RootToken

		// Now seal
		cluster.Cleanup()
		// This will prevent cleanup from running again on the defer
		cluster.Cores = nil
	}

	{
		logger.SetLevel(hclog.Trace)
		logger.Info("integ: start up as normal with shamir seal and unseal, make sure everything is normal")
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		var resp *api.SealStatusResponse
		for _, key := range keys {
			resp, err = client.Sys().Unseal(key)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected response")
			}
		}
		if resp.Sealed {
			t.Fatal("expected unsealed state")
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}

	var autoSeal vault.Seal

	{
		logger.SetLevel(hclog.Trace)
		logger.Info("integ: creating an autoseal and activating migration")
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		core := cluster.Cores[0].Core

		newSeal := vault.NewAutoSeal(seal.NewTestSeal(nil))
		newSeal.SetCore(core)
		autoSeal = newSeal
		if err := adjustCoreForSealMigration(logger, core, newSeal, nil); err != nil {
			t.Fatal(err)
		}

		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		var resp *api.SealStatusResponse
		unsealOpts := &api.UnsealOpts{}
		for _, key := range keys {
			unsealOpts.Key = key
			unsealOpts.Migrate = false
			resp, err = client.Sys().UnsealWithOptions(unsealOpts)
			if err == nil {
				t.Fatal("expected error due to lack of migrate parameter")
			}
			unsealOpts.Migrate = true
			resp, err = client.Sys().UnsealWithOptions(unsealOpts)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected response")
			}
		}
		if resp.Sealed {
			t.Fatalf("expected unsealed state; got %#v", *resp)
		}

		// Make sure the seal configs were updated correctly
		b, err := autoSeal.BarrierConfig(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if b.Type != autoSeal.BarrierType() {
			t.Fatalf("bad seal config: %#v", b)
		}
		if b.SecretShares != 1 {
			t.Fatalf("bad seal config: %#v", b)
		}
		if b.SecretThreshold != 1 {
			t.Fatalf("bad seal config: %#v", b)
		}
		if b.StoredShares != 1 {
			t.Fatalf("bad seal config: %#v", b)
		}

		r, err := autoSeal.RecoveryConfig(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if r.Type != wrapping.Shamir {
			t.Fatalf("bad seal config: %#v", r)
		}
		if r.SecretShares != 2 {
			t.Fatalf("bad seal config: %#v", r)
		}
		if r.SecretThreshold != 2 {
			t.Fatalf("bad seal config: %#v", r)
		}
		if r.StoredShares != 0 {
			t.Fatalf("bad seal config: %#v", r)
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}

	{
		logger.SetLevel(hclog.Trace)
		logger.Info("integ: verify autoseal and recovery key usage")
		coreConfig.Seal = autoSeal
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		core := cluster.Cores[0].Core
		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		if err := core.UnsealWithStoredKeys(ctx); err != nil {
			t.Fatal(err)
		}
		resp, err := client.Sys().SealStatus()
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected response")
		}
		if resp.Sealed {
			t.Fatalf("expected unsealed state; got %#v", *resp)
		}

		keyParts := [][]byte{}
		for _, key := range keys {
			raw, err := base64.StdEncoding.DecodeString(key)
			if err != nil {
				t.Fatal(err)
			}
			keyParts = append(keyParts, raw)
		}
		recoveredKey, err := shamir.Combine(keyParts)
		if err != nil {
			t.Fatal(err)
		}
		sealAccess := core.SealAccess()
		if err := sealAccess.VerifyRecoveryKey(ctx, recoveredKey); err != nil {
			t.Fatal(err)
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}

	// We should see stored barrier keys; after the sixth test, we shouldn't
	if entry, err := phys.Get(ctx, vault.StoredBarrierKeysPath); err != nil || entry == nil {
		t.Fatalf("expected nil error and non-nil entry, got error %#v and entry %#v", err, entry)
	}

	altTestSeal := seal.NewTestSeal(nil)
	altTestSeal.SetType("test-alternate")
	altSeal := vault.NewAutoSeal(altTestSeal)

	{
		logger.SetLevel(hclog.Trace)
		logger.Info("integ: migrate from auto-seal to auto-seal")
		coreConfig.Seal = autoSeal
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		core := cluster.Cores[0].Core

		if err := adjustCoreForSealMigration(logger, core, altSeal, autoSeal); err != nil {
			t.Fatal(err)
		}

		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		var resp *api.SealStatusResponse
		unsealOpts := &api.UnsealOpts{}
		for _, key := range keys {
			unsealOpts.Key = key
			unsealOpts.Migrate = true
			resp, err = client.Sys().UnsealWithOptions(unsealOpts)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected response")
			}
		}
		if resp.Sealed {
			t.Fatalf("expected unsealed state; got %#v", *resp)
		}

		// Make sure the seal configs were updated correctly
		b, err := altSeal.BarrierConfig(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if b.Type != altSeal.BarrierType() {
			t.Fatalf("bad seal config: %#v", b)
		}
		if b.SecretShares != 1 {
			t.Fatalf("bad seal config: %#v", b)
		}
		if b.SecretThreshold != 1 {
			t.Fatalf("bad seal config: %#v", b)
		}
		if b.StoredShares != 1 {
			t.Fatalf("bad seal config: %#v", b)
		}

		r, err := altSeal.RecoveryConfig(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if r.Type != wrapping.Shamir {
			t.Fatalf("bad seal config: %#v", r)
		}
		if r.SecretShares != 2 {
			t.Fatalf("bad seal config: %#v", r)
		}
		if r.SecretThreshold != 2 {
			t.Fatalf("bad seal config: %#v", r)
		}
		if r.StoredShares != 0 {
			t.Fatalf("bad seal config: %#v", r)
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}

	{
		logger.SetLevel(hclog.Trace)
		logger.Info("integ: create a Shamir seal and activate migration; verify it doesn't work if disabled isn't set.")
		coreConfig.Seal = altSeal
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		core := cluster.Cores[0].Core

		wrapper := vault.NewDefaultSeal(&seal.Access{
			Wrapper: aeadwrapper.NewWrapper(&wrapping.WrapperOptions{
				Logger: logger.Named("shamir"),
			}),
		})

		if err := adjustCoreForSealMigration(logger, core, wrapper, altSeal); err != nil {
			t.Fatal(err)
		}

		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		var resp *api.SealStatusResponse
		unsealOpts := &api.UnsealOpts{}
		for _, key := range keys {
			unsealOpts.Key = key
			unsealOpts.Migrate = false
			resp, err = client.Sys().UnsealWithOptions(unsealOpts)
			if err == nil {
				t.Fatal("expected error due to lack of migrate parameter")
			}
			unsealOpts.Migrate = true
			resp, err = client.Sys().UnsealWithOptions(unsealOpts)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected response")
			}
		}
		if resp.Sealed {
			t.Fatalf("expected unsealed state; got %#v", *resp)
		}

		// Make sure the seal configs were updated correctly
		b, err := wrapper.BarrierConfig(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if b.Type != wrapping.Shamir {
			t.Fatalf("bad seal config: %#v", b)
		}
		if b.SecretShares != 2 {
			t.Fatalf("bad seal config: %#v", b)
		}
		if b.SecretThreshold != 2 {
			t.Fatalf("bad seal config: %#v", b)
		}
		if b.StoredShares != 1 {
			t.Fatalf("bad seal config: %#v", b)
		}

		_, err = wrapper.RecoveryConfig(context.Background())
		if err == nil {
			t.Fatal("expected error")
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}

	{
		logger.SetLevel(hclog.Trace)
		logger.Info("integ: verify autoseal is off and the expected key shares work")
		coreConfig.Seal = wrapper
		cluster := vault.NewTestCluster(t, coreConfig, clusterConfig)
		cluster.Start()
		defer cluster.Cleanup()

		core := cluster.Cores[0].Core
		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		if err := core.UnsealWithStoredKeys(ctx); err != nil {
			t.Fatal(err)
		}
		resp, err := client.Sys().SealStatus()
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected response")
		}
		if !resp.Sealed {
			t.Fatalf("expected sealed state; got %#v", *resp)
		}

		for _, key := range keys {
			resp, err = client.Sys().Unseal(key)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected response")
			}
		}
		if resp.Sealed {
			t.Fatal("expected unsealed state")
		}

		cluster.Cleanup()
		cluster.Cores = nil
	}
}
