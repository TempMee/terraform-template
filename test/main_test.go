package test

import (
	"errors"
	"fmt"
	test_structure "github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"log"
	"os"
	"testing"
)

const localBackend = `
terraform {
	backend "local" {}
}
`

func setupTest() (string, error) {
	terraformTempDir, errCopying := test_structure.CopyTerragruntFolderToTemp("../", "terratest-")
	if errCopying != nil {
		return "", errCopying
	}

	backendFilePath := fmt.Sprintf("%s/%s", terraformTempDir, "backend.tf")
	// check if backendFilePath exists
	_, err := os.Stat(backendFilePath)
	if !errors.Is(err, os.ErrNotExist) {
		// if it exists, remove it
		errRemoving := os.Remove(backendFilePath)
		if errRemoving != nil {
			return "", errRemoving
		}
	}

	errWritingFile := os.WriteFile(backendFilePath, []byte(localBackend), os.ModeAppend)
	if errWritingFile != nil {
		return "", errWritingFile
	}
	os.Chmod(backendFilePath, fs.FileMode(0777))
	return terraformTempDir, nil
}

func TestTerraformCodeInfrastructureInitialCredentials(t *testing.T) {
	terraformTempDir, errSettingUpTest := setupTest()
	if errSettingUpTest != nil {
		t.Fatalf("Error setting up test :%v", errSettingUpTest)
	}
	defer os.RemoveAll(terraformTempDir)
	log.Printf("Temp folder: %s", terraformTempDir)
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	tfvars := map[string]interface{}{
		"aws_region": "us-east-1",
	}
	terraformInitOptions := &terraform.Options{
		TerraformDir: terraformTempDir,
		Vars:         tfvars,
		VarFiles:     []string{path + "/terratest.tfvars"},
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": "us-east-1",
			//"TF_LOG":             "TRACE",
		},
		Reconfigure: true,
	}

	defer destroy(t, terraformTempDir, &tfvars)
	terraform.Init(t, terraformInitOptions)
	terraformValidateOptions := &terraform.Options{
		TerraformDir: terraformTempDir,
	}
	terraform.Validate(t, terraformValidateOptions)
	plan, errApplyingIdempotent := terraform.ApplyE(t, terraformInitOptions)
	if errApplyingIdempotent != nil {
		t.Logf("Error applying plan: %v", errApplyingIdempotent)
		t.Fail()
	} else {
		t.Log(fmt.Sprintf("Plan worked: %s", plan))
	}

	dummy := terraform.Output(t, terraformInitOptions, "dummy")

	t.Run("Verify output", func(t *testing.T) {
		a := assert.New(t)
		// check that rabbitmq exists as a resource

		a.Equal("dummy", dummy)
	})

}
