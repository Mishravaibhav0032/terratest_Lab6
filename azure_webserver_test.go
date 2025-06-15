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
	subscriptionId := "eed2085b-133e-4466-8ed6-6da86c403fc0" // your Azure subscription ID
	labelPrefix := "mish0032A05VM"                           // your labelPrefix used in Terraform

	terraformOptions := &terraform.Options{
		TerraformDir: "../", // path to Terraform code
		Vars: map[string]interface{}{
			"labelPrefix": labelPrefix,
		},
		EnvVars: map[string]string{
			"ARM_SUBSCRIPTION_ID": subscriptionId,
		},
	}

	defer func() {
		time.Sleep(20 * time.Second)
		// terraform.Destroy(t, terraformOptions) // enable if you want cleanup
	}()

	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")
	resourceGroup := terraform.Output(t, terraformOptions, "resource_group_name")

	assert.NotEmpty(t, vmName)

	// ✅ Confirm NIC is attached
	nicList, err := azure.GetVirtualMachineNicsE(vmName, resourceGroup, subscriptionId)
	assert.NoError(t, err)
	assert.NotEmpty(t, nicList)
	assert.Contains(t, nicList, nicName)

	networkInterface := azure.GetNetworkInterface(t, nicName, resourceGroup, subscriptionId)
	vm := azure.GetLinuxVirtualMachine(t, vmName, resourceGroup, subscriptionId)
	assert.Equal(t, *vm.ID, *networkInterface.VirtualMachine.ID)

	// ✅ Confirm the correct Ubuntu version
	vmImage, err := azure.GetVirtualMachineImageE(vmName, resourceGroup, subscriptionId)
	assert.NoError(t, err)

	imageSKU := vmImage.SKU
	assert.True(t,
		strings.Contains(imageSKU, "18.04") || strings.Contains(imageSKU, "20_04") || strings.Contains(imageSKU, "22_04"),
		"VM image SKU should be Ubuntu 18.04, 20.04, or 22.04. Got: %s", imageSKU)
}
