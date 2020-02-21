package config

import (
	"os"
	"testing"
	"time"

	"github.com/go-test/deep"
	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
)

func TestLoadConfigFile_AgentCache(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				&Sink{
					Type:   "file",
					DHType: "curve25519",
					DHPath: "/tmp/file-foo-dhpath",
					AAD:    "foobar",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
			},
		},
		Cache: &Cache{
			UseAutoAuthToken:        true,
			UseAutoAuthTokenEnforce: false,
			UseAutoAuthTokenRaw:     true,
		},
		Listeners: []*Listener{
			&Listener{
				Type: "unix",
				Config: map[string]interface{}{
					"address":      "/path/to/socket",
					"tls_disable":  true,
					"socket_mode":  "configmode",
					"socket_user":  "configuser",
					"socket_group": "configgroup",
				},
			},
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":     "127.0.0.1:8300",
					"tls_disable": true,
				},
			},
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":       "127.0.0.1:8400",
					"tls_key_file":  "/path/to/cakey.pem",
					"tls_cert_file": "/path/to/cacert.pem",
				},
			},
		},
		Vault: &Vault{
			Address:          "http://127.0.0.1:1111",
			CACert:           "config_ca_cert",
			CAPath:           "config_ca_path",
			TLSSkipVerifyRaw: interface{}("true"),
			TLSSkipVerify:    true,
			ClientCert:       "config_client_cert",
			ClientKey:        "config_client_key",
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}

	config, err = LoadConfig("./test-fixtures/config-cache-embedded-type.hcl")
	if err != nil {
		t.Fatal(err)
	}
	expected.Vault.TLSSkipVerifyRaw = interface{}(true)

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile(t *testing.T) {
	if err := os.Setenv("TEST_AAD_ENV", "aad"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Unsetenv("TEST_AAD_ENV"); err != nil {
			t.Fatal(err)
		}
	}()

	config, err := LoadConfig("./test-fixtures/config.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				&Sink{
					Type:   "file",
					DHType: "curve25519",
					DHPath: "/tmp/file-foo-dhpath",
					AAD:    "foobar",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
				&Sink{
					Type:    "file",
					WrapTTL: 5 * time.Minute,
					DHType:  "curve25519",
					DHPath:  "/tmp/file-foo-dhpath2",
					AAD:     "aad",
					Config: map[string]interface{}{
						"path": "/tmp/file-bar",
					},
				},
			},
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}

	config, err = LoadConfig("./test-fixtures/config-embedded-type.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Method_Wrapping(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-method-wrapping.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				WrapTTL:   5 * time.Minute,
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				&Sink{
					Type: "file",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
			},
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_NoAutoAuth(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-no-auto_auth.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Cache: &Cache{},
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":     "127.0.0.1:8300",
					"tls_disable": true,
				},
			},
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Bad_AgentCache_InconsisentAutoAuth(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-cache-inconsistent-auto_auth.hcl")
	if err == nil {
		t.Fatal("LoadConfig should return an error when use_auto_auth_token=true and no auto_auth section present")
	}
}

func TestLoadConfigFile_Bad_AgentCache_NoListeners(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-cache-no-listeners.hcl")
	if err == nil {
		t.Fatal("LoadConfig should return an error when cache section present and no listeners present")
	}
}

func TestLoadConfigFile_Bad_AutoAuth_Wrapped_Multiple_Sinks(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-auto_auth-wrapped-multiple-sinks")
	if err == nil {
		t.Fatal("LoadConfig should return an error when auth_auth.method.wrap_ttl nonzero and multiple sinks defined")
	}
}

func TestLoadConfigFile_Bad_AutoAuth_Both_Wrapping_Types(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-method-wrapping-and-sink-wrapping.hcl")
	if err == nil {
		t.Fatal("LoadConfig should return an error when auth_auth.method.wrap_ttl nonzero and sinks.wrap_ttl nonzero")
	}
}

func TestLoadConfigFile_Bad_AgentCache_AutoAuth_Method_wrapping(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-cache-auto_auth-method-wrapping.hcl")
	if err == nil {
		t.Fatal("LoadConfig should return an error when auth_auth.method.wrap_ttl nonzero and cache.use_auto_auth_token=true")
	}
}

func TestLoadConfigFile_AgentCache_AutoAuth_NoSink(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-auto_auth-no-sink.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
		},
		Cache: &Cache{
			UseAutoAuthToken:        true,
			UseAutoAuthTokenEnforce: false,
			UseAutoAuthTokenRaw:     true,
		},
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":     "127.0.0.1:8300",
					"tls_disable": true,
				},
			},
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_AutoAuth_Force(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-auto_auth-force.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
		},
		Cache: &Cache{
			UseAutoAuthToken:        true,
			UseAutoAuthTokenEnforce: true,
			UseAutoAuthTokenRaw:     "force",
		},
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":     "127.0.0.1:8300",
					"tls_disable": true,
				},
			},
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_AutoAuth_True(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-auto_auth-true.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
		},
		Cache: &Cache{
			UseAutoAuthToken:        true,
			UseAutoAuthTokenEnforce: false,
			UseAutoAuthTokenRaw:     "true",
		},
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":     "127.0.0.1:8300",
					"tls_disable": true,
				},
			},
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_AutoAuth_False(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-auto_auth-false.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				&Sink{
					Type:   "file",
					DHType: "curve25519",
					DHPath: "/tmp/file-foo-dhpath",
					AAD:    "foobar",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
			},
		},
		Cache: &Cache{
			UseAutoAuthToken:        false,
			UseAutoAuthTokenEnforce: false,
			UseAutoAuthTokenRaw:     "false",
		},
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":     "127.0.0.1:8300",
					"tls_disable": true,
				},
			},
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

// TestLoadConfigFile_Template tests template definitions in Vault Agent
func TestLoadConfigFile_Template(t *testing.T) {
	testCases := map[string]struct {
		fixturePath       string
		expectedTemplates []*ctconfig.TemplateConfig
	}{
		"min": {
			fixturePath: "./test-fixtures/config-template-min.hcl",
			expectedTemplates: []*ctconfig.TemplateConfig{
				&ctconfig.TemplateConfig{
					Source:      pointerutil.StringPtr("/path/on/disk/to/template.ctmpl"),
					Destination: pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
				},
			},
		},
		"full": {
			fixturePath: "./test-fixtures/config-template-full.hcl",
			expectedTemplates: []*ctconfig.TemplateConfig{
				&ctconfig.TemplateConfig{
					Backup:         pointerutil.BoolPtr(true),
					Command:        pointerutil.StringPtr("restart service foo"),
					CommandTimeout: pointerutil.TimeDurationPtr("60s"),
					Contents:       pointerutil.StringPtr("{{ keyOrDefault \"service/redis/maxconns@east-aws\" \"5\" }}"),
					CreateDestDirs: pointerutil.BoolPtr(true),
					Destination:    pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
					ErrMissingKey:  pointerutil.BoolPtr(true),
					LeftDelim:      pointerutil.StringPtr("<<"),
					Perms:          pointerutil.FileModePtr(0655),
					RightDelim:     pointerutil.StringPtr(">>"),
					SandboxPath:    pointerutil.StringPtr("/path/on/disk/where"),

					Wait: &ctconfig.WaitConfig{
						Min: pointerutil.TimeDurationPtr("10s"),
						Max: pointerutil.TimeDurationPtr("40s"),
					},
				},
			},
		},
		"many": {
			fixturePath: "./test-fixtures/config-template-many.hcl",
			expectedTemplates: []*ctconfig.TemplateConfig{
				&ctconfig.TemplateConfig{
					Source:         pointerutil.StringPtr("/path/on/disk/to/template.ctmpl"),
					Destination:    pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
					ErrMissingKey:  pointerutil.BoolPtr(false),
					CreateDestDirs: pointerutil.BoolPtr(true),
					Command:        pointerutil.StringPtr("restart service foo"),
					Perms:          pointerutil.FileModePtr(0600),
				},
				&ctconfig.TemplateConfig{
					Source:      pointerutil.StringPtr("/path/on/disk/to/template2.ctmpl"),
					Destination: pointerutil.StringPtr("/path/on/disk/where/template/will/render2.txt"),
					Backup:      pointerutil.BoolPtr(true),
					Perms:       pointerutil.FileModePtr(0755),
					Wait: &ctconfig.WaitConfig{
						Min: pointerutil.TimeDurationPtr("2s"),
						Max: pointerutil.TimeDurationPtr("10s"),
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			config, err := LoadConfig(tc.fixturePath)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			expected := &Config{
				AutoAuth: &AutoAuth{
					Method: &Method{
						Type:      "aws",
						MountPath: "auth/aws",
						Namespace: "my-namespace/",
						Config: map[string]interface{}{
							"role": "foobar",
						},
					},
					Sinks: []*Sink{
						&Sink{
							Type:   "file",
							DHType: "curve25519",
							DHPath: "/tmp/file-foo-dhpath",
							AAD:    "foobar",
							Config: map[string]interface{}{
								"path": "/tmp/file-foo",
							},
						},
					},
				},
				Templates: tc.expectedTemplates,
				PidFile:   "./pidfile",
			}

			if diff := deep.Equal(config, expected); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}
