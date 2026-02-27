package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	jiratypes "github.com/surajrajput1024/go-atlassian-cloud/types"
)

func projectResponseToResourceModel(project *jiratypes.ProjectResponse) JiraProjectResourceModel {
	var state JiraProjectResourceModel
	state.ID = types.StringValue(project.ID)
	state.Key = types.StringValue(project.Key)
	state.Name = types.StringValue(project.Name)
	if project.Description != "" {
		state.Description = types.StringValue(project.Description)
	} else {
		state.Description = types.StringNull()
	}
	if project.Lead != nil && project.Lead.AccountID != "" {
		state.LeadAccountID = types.StringValue(project.Lead.AccountID)
	} else {
		state.LeadAccountID = types.StringNull()
	}
	return state
}
func createProjectErrorMessage(apiErr error) string {
	errMsg := apiErr.Error()
	if strings.Contains(errMsg, "already exists") ||
		strings.Contains(errMsg, "projectName") ||
		strings.Contains(errMsg, "projectKey") {
		return fmt.Sprintf("%s. If the project already exists in Jira, import it with: terraform import 'atlassian_jira_project.<resource>' <PROJECT_ID_OR_KEY>. Otherwise use a unique project key and name.", errMsg)
	}
	return errMsg
}
