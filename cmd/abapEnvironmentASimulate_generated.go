// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperenv"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/spf13/cobra"
)

type abapEnvironmentASimulateOptions struct {
	PackageType            string `json:"PackageType,omitempty"`
	PackageName            string `json:"PackageName,omitempty"`
	SWC                    string `json:"SWC,omitempty"`
	CVERS                  string `json:"CVERS,omitempty"`
	Namespace              string `json:"Namespace,omitempty"`
	PreviousDeliveryCommit string `json:"PreviousDeliveryCommit,omitempty"`
}

type abapEnvironmentASimulateCommonPipelineEnvironment struct {
	PackageName            string
	PackageType            string
	SWC                    string
	CVERS                  string
	Namespace              string
	PreviousDeliveryCommit string
}

func (p *abapEnvironmentASimulateCommonPipelineEnvironment) persist(path, resourceName string) {
	content := []struct {
		category string
		name     string
		value    string
	}{
		{category: "", name: "PackageName", value: p.PackageName},
		{category: "", name: "PackageType", value: p.PackageType},
		{category: "", name: "SWC", value: p.SWC},
		{category: "", name: "CVERS", value: p.CVERS},
		{category: "", name: "Namespace", value: p.Namespace},
		{category: "", name: "PreviousDeliveryCommit", value: p.PreviousDeliveryCommit},
	}

	errCount := 0
	for _, param := range content {
		err := piperenv.SetResourceParameter(path, resourceName, filepath.Join(param.category, param.name), param.value)
		if err != nil {
			log.Entry().WithError(err).Error("Error persisting piper environment.")
			errCount++
		}
	}
	if errCount > 0 {
		log.Entry().Fatal("failed to persist Piper environment")
	}
}

// AbapEnvironmentASimulateCommand simulates writing to the commonPipelineEnvironment
func AbapEnvironmentASimulateCommand() *cobra.Command {
	const STEP_NAME = "abapEnvironmentASimulate"

	metadata := abapEnvironmentASimulateMetadata()
	var stepConfig abapEnvironmentASimulateOptions
	var startTime time.Time
	var commonPipelineEnvironment abapEnvironmentASimulateCommonPipelineEnvironment

	var createAbapEnvironmentASimulateCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "simulates writing to the commonPipelineEnvironment",
		Long:  `simulates writing to the commonPipelineEnvironment`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			telemetryData := telemetry.CustomData{}
			telemetryData.ErrorCode = "1"
			handler := func() {
				commonPipelineEnvironment.persist(GeneralConfig.EnvRootPath, "commonPipelineEnvironment")
				telemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				telemetry.Send(&telemetryData)
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetry.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			abapEnvironmentASimulate(stepConfig, &telemetryData, &commonPipelineEnvironment)
			telemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addAbapEnvironmentASimulateFlags(createAbapEnvironmentASimulateCmd, &stepConfig)
	return createAbapEnvironmentASimulateCmd
}

func addAbapEnvironmentASimulateFlags(cmd *cobra.Command, stepConfig *abapEnvironmentASimulateOptions) {
	cmd.Flags().StringVar(&stepConfig.PackageType, "PackageType", `AOI`, "Package Types. Valid values: 'AOI', 'CSP', 'CPK'.")
	cmd.Flags().StringVar(&stepConfig.PackageName, "PackageName", os.Getenv("PIPER_PackageName"), "Package Name")
	cmd.Flags().StringVar(&stepConfig.SWC, "SWC", os.Getenv("PIPER_SWC"), "SWC")
	cmd.Flags().StringVar(&stepConfig.CVERS, "CVERS", os.Getenv("PIPER_CVERS"), "Cvers")
	cmd.Flags().StringVar(&stepConfig.Namespace, "Namespace", os.Getenv("PIPER_Namespace"), "Namespace")
	cmd.Flags().StringVar(&stepConfig.PreviousDeliveryCommit, "PreviousDeliveryCommit", os.Getenv("PIPER_PreviousDeliveryCommit"), "Previous delivery commit")

	cmd.MarkFlagRequired("PackageType")
	cmd.MarkFlagRequired("PackageName")
	cmd.MarkFlagRequired("SWC")
	cmd.MarkFlagRequired("CVERS")
	cmd.MarkFlagRequired("Namespace")
	cmd.MarkFlagRequired("PreviousDeliveryCommit")
}

// retrieve step metadata
func abapEnvironmentASimulateMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:    "abapEnvironmentASimulate",
			Aliases: []config.Alias{},
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Parameters: []config.StepParameters{
					{
						Name:        "PackageType",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "PackageName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "SWC",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "CVERS",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "Namespace",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "PreviousDeliveryCommit",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
				},
			},
		},
	}
	return theMetaData
}