package test

import (
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAzureLinuxVMCreation(t *testing.T) {
	subscriptionId := "eed2085b-133e-4466-8ed6-6da86c403fc0" // your actual Azure subscription ID
	labelPrefix := "mish0032A05VM"                           // your actual labelPrefix used in Terraform

	terraformOptions := &terraform.Options{
		TerraformDir: "../", // path to your Terraform code
		Vars: map[string]interface{}{
			"labelPrefix": labelPrefix,
		},
		EnvVars: map[string]string{
			"ARM_SUBSCRIPTION_ID": subscriptionId,
		},
	}

	defer func() {
		time.Sleep(20 * time.Second)
		// terraform.Destroy(t, terraformOptions)
	}()

	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")
	resourceGroup := terraform.Output(t, terraformOptions, "resource_group_name")

	assert.NotEmpty(t, vmName)

	nicList, err := azure.GetVirtualMachineNicsE(vmName, resourceGroup, subscriptionId)
	assert.NoError(t, err)
	assert.NotEmpty(t, nicList)
	assert.Contains(t, nicList, nicName)

	vmImage, err := azure.GetVirtualMachineImageE(vmName, resourceGroup, subscriptionId)
	assert.NoError(t, err)
	assert.True(t,
		strings.Contains(vmImage.SKU, "18.04") || strings.Contains(vmImage.SKU, "22_04"),
		"VM image SKU should be either 18.04 or 22_04, got: %s", vmImage.SKU)
}
