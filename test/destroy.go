package test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func destroy(t *testing.T, terraformTempDir string, vars *map[string]interface{}) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path) // for example /home/user
	var terraformDestroyOptions *terraform.Options
	if vars != nil {
		terraformDestroyOptions = &terraform.Options{
			TerraformDir: terraformTempDir,
			//VarFiles:     []string{"test/terratest.tfvars"},
			VarFiles: []string{path + "/terratest.tfvars"},
			Vars:     *vars,
			EnvVars: map[string]string{
				"AWS_DEFAULT_REGION": "us-east-1",
				//"TF_LOG":             "TRACE",
			},
			Reconfigure: true,
		}
	} else {
		terraformDestroyOptions = &terraform.Options{
			TerraformDir: terraformTempDir,
			//VarFiles:     []string{"test/terratest.tfvars"},
			VarFiles: []string{path + "/terratest.tfvars"},
			EnvVars: map[string]string{
				"AWS_DEFAULT_REGION": "us-east-1",
				//"TF_LOG":             "TRACE",
			},
			Reconfigure: true,
		}
	}
	log.Println("Destroying resources...")
	terraform.Destroy(t, terraformDestroyOptions)

}
