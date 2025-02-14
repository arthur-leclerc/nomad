package api

import (
	"testing"

	"github.com/hashicorp/nomad/api/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestACLPolicies_ListUpsert(t *testing.T) {
	testutil.Parallel(t)
	c, s, _ := makeACLClient(t, nil, nil)
	defer s.Stop()
	ap := c.ACLPolicies()

	// Listing when nothing exists returns empty
	result, qm, err := ap.List(nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if qm.LastIndex != 1 {
		t.Fatalf("bad index: %d", qm.LastIndex)
	}
	if n := len(result); n != 0 {
		t.Fatalf("expected 0 policies, got: %d", n)
	}

	// Register a policy
	policy := &ACLPolicy{
		Name:        "test",
		Description: "test",
		Rules: `namespace "default" {
			policy = "read"
		}
		`,
	}
	wm, err := ap.Upsert(policy, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)

	// Check the list again
	result, qm, err = ap.List(nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assertQueryMeta(t, qm)
	if len(result) != 1 {
		t.Fatalf("expected policy, got: %#v", result)
	}
}

func TestACLPolicies_Delete(t *testing.T) {
	testutil.Parallel(t)
	c, s, _ := makeACLClient(t, nil, nil)
	defer s.Stop()
	ap := c.ACLPolicies()

	// Register a policy
	policy := &ACLPolicy{
		Name:        "test",
		Description: "test",
		Rules: `namespace "default" {
			policy = "read"
		}
		`,
	}
	wm, err := ap.Upsert(policy, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)

	// Delete the policy
	wm, err = ap.Delete(policy.Name, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)

	// Check the list again
	result, qm, err := ap.List(nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assertQueryMeta(t, qm)
	if len(result) != 0 {
		t.Fatalf("unexpected policy, got: %#v", result)
	}
}

func TestACLPolicies_Info(t *testing.T) {
	testutil.Parallel(t)
	c, s, _ := makeACLClient(t, nil, nil)
	defer s.Stop()
	ap := c.ACLPolicies()

	// Register a policy
	policy := &ACLPolicy{
		Name:        "test",
		Description: "test",
		Rules: `namespace "default" {
			policy = "read"
		}
		`,
	}
	wm, err := ap.Upsert(policy, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)

	// Query the policy
	out, qm, err := ap.Info(policy.Name, nil)
	assert.Nil(t, err)
	assertQueryMeta(t, qm)
	assert.Equal(t, policy.Name, out.Name)
}

func TestACLTokens_List(t *testing.T) {
	testutil.Parallel(t)
	c, s, _ := makeACLClient(t, nil, nil)
	defer s.Stop()
	at := c.ACLTokens()

	// Expect out bootstrap token
	result, qm, err := at.List(nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if qm.LastIndex == 0 {
		t.Fatalf("bad index: %d", qm.LastIndex)
	}
	if n := len(result); n != 1 {
		t.Fatalf("expected 1 token, got: %d", n)
	}
}

func TestACLTokens_CreateUpdate(t *testing.T) {
	testutil.Parallel(t)
	c, s, _ := makeACLClient(t, nil, nil)
	defer s.Stop()
	at := c.ACLTokens()

	token := &ACLToken{
		Name:     "foo",
		Type:     "client",
		Policies: []string{"foo1"},
	}

	// Create the token
	out, wm, err := at.Create(token, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)
	assert.NotNil(t, out)

	// Update the token
	out.Name = "other"
	out2, wm, err := at.Update(out, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)
	assert.NotNil(t, out2)

	// Verify the change took hold
	assert.Equal(t, out.Name, out2.Name)
}

func TestACLTokens_Info(t *testing.T) {
	testutil.Parallel(t)
	c, s, _ := makeACLClient(t, nil, nil)
	defer s.Stop()
	at := c.ACLTokens()

	token := &ACLToken{
		Name:     "foo",
		Type:     "client",
		Policies: []string{"foo1"},
	}

	// Create the token
	out, wm, err := at.Create(token, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)
	assert.NotNil(t, out)

	// Query the token
	out2, qm, err := at.Info(out.AccessorID, nil)
	assert.Nil(t, err)
	assertQueryMeta(t, qm)
	assert.Equal(t, out, out2)
}

func TestACLTokens_Self(t *testing.T) {
	testutil.Parallel(t)
	c, s, _ := makeACLClient(t, nil, nil)
	defer s.Stop()
	at := c.ACLTokens()

	token := &ACLToken{
		Name:     "foo",
		Type:     "client",
		Policies: []string{"foo1"},
	}

	// Create the token
	out, wm, err := at.Create(token, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)
	assert.NotNil(t, out)

	// Set the clients token to the new token
	c.SetSecretID(out.SecretID)
	at = c.ACLTokens()

	// Query the token
	out2, qm, err := at.Self(nil)
	if assert.Nil(t, err) {
		assertQueryMeta(t, qm)
		assert.Equal(t, out, out2)
	}
}

func TestACLTokens_Delete(t *testing.T) {
	testutil.Parallel(t)
	c, s, _ := makeACLClient(t, nil, nil)
	defer s.Stop()
	at := c.ACLTokens()

	token := &ACLToken{
		Name:     "foo",
		Type:     "client",
		Policies: []string{"foo1"},
	}

	// Create the token
	out, wm, err := at.Create(token, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)
	assert.NotNil(t, out)

	// Delete the token
	wm, err = at.Delete(out.AccessorID, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)
}

func TestACL_OneTimeToken(t *testing.T) {
	testutil.Parallel(t)
	c, s, _ := makeACLClient(t, nil, nil)
	defer s.Stop()
	at := c.ACLTokens()

	token := &ACLToken{
		Name:     "foo",
		Type:     "client",
		Policies: []string{"foo1"},
	}

	// Create the ACL token
	out, wm, err := at.Create(token, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)
	assert.NotNil(t, out)

	// Get a one-time token
	c.SetSecretID(out.SecretID)
	out2, wm, err := at.UpsertOneTimeToken(nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)
	assert.NotNil(t, out2)

	// Exchange the one-time token
	out3, wm, err := at.ExchangeOneTimeToken(out2.OneTimeSecretID, nil)
	assert.Nil(t, err)
	assertWriteMeta(t, wm)
	assert.NotNil(t, out3)
	assert.Equal(t, out3.AccessorID, out.AccessorID)
}

func TestACLTokens_BootstrapInvalidToken(t *testing.T) {
	testutil.Parallel(t)
	c, s := makeClient(t, nil, func(c *testutil.TestServerConfig) {
		c.ACL.Enabled = true
	})
	defer s.Stop()
	at := c.ACLTokens()

	bootkn := "badtoken"
	// Bootstrap with invalid token
	_, _, err := at.BootstrapOpts(bootkn, nil)
	assert.EqualError(t, err, "Unexpected response code: 400 (invalid acl token)")
}

func TestACLTokens_BootstrapValidToken(t *testing.T) {
	testutil.Parallel(t)
	c, s := makeClient(t, nil, func(c *testutil.TestServerConfig) {
		c.ACL.Enabled = true
	})
	defer s.Stop()
	at := c.ACLTokens()

	bootkn := "2b778dd9-f5f1-6f29-b4b4-9a5fa948757a"
	// Bootstrap with Valid token
	out, wm, err := at.BootstrapOpts(bootkn, nil)
	assert.NoError(t, err)
	assertWriteMeta(t, wm)
	assert.Equal(t, bootkn, out.SecretID)
}
