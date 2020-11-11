package octopusdeploy

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pemaxim/go-octopusdeploy/octopusdeploy"
)

func getDeploymentStepSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "The name of the step",
					Required:    true,
				},
				constTargetRoles: {
					Description: "The roles that this step run against, or runs on behalf of",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				constPackageRequirement: {
					Type:        schema.TypeString,
					Description: "Whether to run this step before or after package acquisition (if possible)",
					Optional:    true,
					Default:     (string)(octopusdeploy.DeploymentStepPackageRequirementLetOctopusDecide),
					ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
						(string)(octopusdeploy.DeploymentStepPackageRequirementLetOctopusDecide),
						(string)(octopusdeploy.DeploymentStepPackageRequirementBeforePackageAcquisition),
						(string)(octopusdeploy.DeploymentStepPackageRequirementAfterPackageAcquisition),
					}, false)),
				},
				constCondition: {
					Type:        schema.TypeString,
					Description: "When to run the step, one of 'Success', 'Failure', 'Always' or 'Variable'",
					Optional:    true,
					Default:     (string)(octopusdeploy.DeploymentStepConditionTypeSuccess),
					ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
						(string)(octopusdeploy.DeploymentStepConditionTypeSuccess),
						(string)(octopusdeploy.DeploymentStepConditionTypeFailure),
						(string)(octopusdeploy.DeploymentStepConditionTypeAlways),
						(string)(octopusdeploy.DeploymentStepConditionTypeVariable),
					}, false)),
				},
				constConditionExpression: {
					Type:        schema.TypeString,
					Description: "The expression to evaluate to determine whether to run this step when 'condition' is 'Variable'",
					Optional:    true,
				},
				constStartTrigger: {
					Type:        schema.TypeString,
					Description: "Whether to run this step after the previous step ('StartAfterPrevious') or at the same time as the previous step ('StartWithPrevious')",
					Optional:    true,
					Default:     (string)(octopusdeploy.DeploymentStepStartTriggerStartAfterPrevious),
					ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
						(string)(octopusdeploy.DeploymentStepStartTriggerStartAfterPrevious),
						(string)(octopusdeploy.DeploymentStepStartTriggerStartWithPrevious),
					}, false)),
				},
				constWindowSize: {
					Type:        schema.TypeString,
					Description: "The maximum number of targets to deploy to simultaneously",
					Optional:    true,
				},
				constAction:                       getDeploymentActionSchema(),
				constManualInterventionAction:     getManualInterventionActionSchema(),
				constApplyTerraformAction:         getApplyTerraformActionSchema(),
				constDeployPackageAction:          getDeployPackageAction(),
				constDeployWindowsServiceAction:   getDeployWindowsServiceActionSchema(),
				constRunScriptAction:              getRunScriptActionSchema(),
				constRunKubectlScriptAction:       getRunRunKubectlScriptSchema(),
				constDeployKubernetesSecretAction: getDeployKubernetesSecretActionSchema(),
			},
		},
	}
}

func buildDeploymentStepResource(tfStep map[string]interface{}) octopusdeploy.DeploymentStep {
	step := octopusdeploy.DeploymentStep{
		Name:               tfStep[constName].(string),
		PackageRequirement: octopusdeploy.DeploymentStepPackageRequirement(tfStep[constPackageRequirement].(string)),
		Condition:          octopusdeploy.DeploymentStepConditionType(tfStep[constCondition].(string)),
		StartTrigger:       octopusdeploy.DeploymentStepStartTrigger(tfStep[constStartTrigger].(string)),
		Properties:         map[string]string{},
	}

	targetRoles := tfStep[constTargetRoles]
	if targetRoles != nil {
		step.Properties["Octopus.Action.TargetRoles"] = strings.Join(getSliceFromTerraformTypeList(targetRoles), ",")
	}

	conditionExpression := tfStep[constConditionExpression]
	if conditionExpression != nil {
		step.Properties["Octopus.Action.ConditionVariableExpression"] = conditionExpression.(string)
	}

	windowSize := tfStep[constWindowSize]
	if windowSize != nil {
		step.Properties["Octopus.Action.MaxParallelism"] = windowSize.(string)
	}

	if attr, ok := tfStep[constAction]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildDeploymentActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep[constManualInterventionAction]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildManualInterventionActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep[constApplyTerraformAction]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildApplyTerraformActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep[constDeployPackageAction]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildDeployPackageActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep[constDeployWindowsServiceAction]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildDeployWindowsServiceActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep[constRunScriptAction]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildRunScriptActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep[constRunKubectlScriptAction]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildRunKubectlScriptActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep[constDeployKubernetesSecretAction]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildDeployKubernetesSecretActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	return step
}
