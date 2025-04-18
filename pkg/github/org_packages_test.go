package github

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/google/go-github/v69/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_OrgListPackages(t *testing.T) {
	// Verify tool definition once
	mockClient := github.NewClient(nil)
	tool, _ := OrgListPackages(stubGetClientFn(mockClient), translations.NullTranslationHelper)

	assert.Equal(t, "org_list_packages", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.Equal(t, tool.InputSchema.Properties, "organization")
	assert.Equal(t, tool.InputSchema.Properties, "package_type")
	assert.ElementsMatch(t, tool.InputSchema.Required, []string{"organization", "package_type"})

	mockPackage := &github.Package{
		Name:        github.Ptr("test-package"),
		PackageType: github.Ptr("container"),
		Owner: &github.User{
			Login: github.Ptr("test-org"),
		},
		HTMLURL: github.Ptr("https://github.com/orgs/test-org/packages/container/test-package"),
	}

	tests := []struct {
		name            string
		mockedClient    *http.Client
		requestArgs     map[string]interface{}
		expectedPackage *github.Package
		expectedErrMsg  string
		expectError     bool
	}{
		{
			name: "Test with valid parameters",
			mockedClient: mock.NewMockedHTTPClient(
				mock.WithRequestMatch(
					mock.GetOrgsPackagesByOrg,
					mockPackage,
				),
			),
			requestArgs: map[string]interface{}{
				"organization": "test-org",
				"package_type": "container",
			},
			expectedPackage: mockPackage,
			expectError:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup client with mock
			client := github.NewClient(tc.mockedClient)
			_, handler := OrgListPackages(stubGetClientFn(client), translations.NullTranslationHelper)

			// Create call request
			request := createMCPRequest(tc.requestArgs)

			// Call handler
			result, err := handler(context.Background(), request)

			// Verify results
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
				return
			}

			require.NoError(t, err)

			// Parse the result and get the text content if no error
			textContent := getTextResult(t, result)

			// Unmarshal and verify the result
			var returnedPackage github.Package
			err = json.Unmarshal([]byte(textContent.Text), &returnedPackage)
			require.NoError(t, err)
			assert.Equal(t, *tc.expectedPackage.Name, *returnedPackage.Name)
			assert.Equal(t, *tc.expectedPackage.PackageType, *returnedPackage.PackageType)
			assert.Equal(t, *tc.expectedPackage.Owner.Login, *returnedPackage.Owner.Login)
			assert.Equal(t, *tc.expectedPackage.HTMLURL, *returnedPackage.HTMLURL)
		})
	}
}
