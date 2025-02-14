package vault_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/argoproj-labs/argocd-vault-plugin/pkg/auth/vault"
	"github.com/argoproj-labs/argocd-vault-plugin/pkg/utils"
	"github.com/argoproj-labs/argocd-vault-plugin/pkg/helpers"
)

func TestAppRoleLogin(t *testing.T) {
	cluster, roleID, secretID := helpers.CreateTestAppRoleVault(t)
	defer cluster.Cleanup()

	appRole := vault.NewAppRoleAuth(roleID, secretID, "", "default")

	err := appRole.Authenticate(cluster.Cores[0].Client)
	if err != nil {
		t.Fatalf("expected no errors but got: %s", err)
	}

	cachedToken, err := utils.ReadExistingToken(fmt.Sprintf("approle_default_%s", roleID))
	if err != nil {
		t.Fatalf("expected cached vault token but got: %s", err)
	}

	err = appRole.Authenticate(cluster.Cores[0].Client)
	if err != nil {
		t.Fatalf("expected no errors but got: %s", err)
	}

	newCachedToken, err := utils.ReadExistingToken(fmt.Sprintf("approle_default_%s", roleID))
	if err != nil {
		t.Fatalf("expected cached vault token but got: %s", err)
	}

	if bytes.Compare(cachedToken, newCachedToken) != 0 {
		t.Fatalf("expected same token %s but got %s", cachedToken, newCachedToken)
	}

	// We create a new connection with a different approle and create a different cache
	secondCluster, secondRoleID, secondSecretID := helpers.CreateTestAppRoleVault(t)
	defer secondCluster.Cleanup()

	secondAppRole := vault.NewAppRoleAuth(secondRoleID, secondSecretID, "", "default")

	err = secondAppRole.Authenticate(secondCluster.Cores[0].Client)
	if err != nil {
		t.Fatalf("expected no errors but got: %s", err)
	}

	secondCachedToken, err := utils.ReadExistingToken(fmt.Sprintf("approle_default_%s", secondRoleID))
	if err != nil {
		t.Fatalf("expected cached vault token but got: %s", err)
	}

	// Both cache should be different
	if bytes.Compare(cachedToken, secondCachedToken) == 0 {
		t.Fatalf("expected different tokens but got %s", secondCachedToken)
	}

	// We create a new connection to a specific namespace and create a different cache
	namespaceCluster, namespaceRoleID, namespaceSecretID := helpers.CreateTestAppRoleVault(t)
	defer namespaceCluster.Cleanup()

	namespaceAppRole := vault.NewAppRoleAuth(namespaceRoleID, namespaceSecretID, "", "my-other-namespace")

	err = namespaceAppRole.Authenticate(namespaceCluster.Cores[0].Client)
	if err != nil {
		t.Fatalf("expected no errors but got: %s", err)
	}

	_, err = utils.ReadExistingToken(fmt.Sprintf("approle_my-other-namespace_%s", namespaceRoleID))
	if err != nil {
		t.Fatalf("expected cached vault token but got: %s", err)
	}
}
