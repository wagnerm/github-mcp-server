package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/google/go-github/v69/github"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func OrgListPackages(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("org_list_packages",
			mcp.WithDescription(t("TOOL_ORG_LIST_PACKAGES", "List packages for an organization")),
			mcp.WithString("organization",
				mcp.Required(),
				mcp.Description("Organization"),
			),
			mcp.WithString("package_type",
				mcp.Required(),
				mcp.Description("Package type"),
			),
			WithPagination(),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			organization, err := requiredParam[string](request, "organization")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			packageType, err := requiredParam[string](request, "package_type")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			pagination, err := OptionalPaginationParams(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			opts := &github.PackageListOptions{
				ListOptions: github.ListOptions{
					Page:    pagination.page,
					PerPage: pagination.perPage,
				},
				PackageType: &packageType,
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}
			packages, resp, err := client.Organizations.ListPackages(ctx, organization, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list packages: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != 200 {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get organization packages: %s", string(body))), nil
			}

			r, err := json.Marshal(packages)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

func OrgPackageGetAllVersions(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("org_package_get_all_versions",
			mcp.WithDescription(t("TOOL_ORG_LIST_PACKAGES", "List package versions for a package owned by an organization")),
			mcp.WithString("organization",
				mcp.Required(),
				mcp.Description("Organization"),
			),
			mcp.WithString("package_type",
				mcp.Required(),
				mcp.Description("Package type"),
			),
			mcp.WithString("package_name",
				mcp.Description("Package name"),
			),
			WithPagination(),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			organization, err := requiredParam[string](request, "organization")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			packageType, err := requiredParam[string](request, "package_type")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			packageName, err := requiredParam[string](request, "package_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			pagination, err := OptionalPaginationParams(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			opts := &github.PackageListOptions{
				ListOptions: github.ListOptions{
					Page:    pagination.page,
					PerPage: pagination.perPage,
				},
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}
			packages, resp, err := client.Organizations.PackageGetAllVersions(ctx, organization, packageType, packageName, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list packages: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != 200 {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get organization packages: %s", string(body))), nil
			}

			r, err := json.Marshal(packages)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
