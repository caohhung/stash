/*
Copyright AppsCode Inc. and Contributors

Licensed under the PolyForm Noncommercial License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/PolyForm-Noncommercial-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmds

import (
	"fmt"

	cs "stash.appscode.dev/apimachinery/client/clientset/versioned"
	"stash.appscode.dev/apimachinery/pkg/restic"
	"stash.appscode.dev/stash/pkg/status"

	"github.com/appscode/go/flags"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	StashDefaultMetricJob = "stash-prom-metrics"
)

func NewCmdUpdateStatus() *cobra.Command {
	var (
		masterURL      string
		kubeconfigPath string
		opt            = status.UpdateStatusOptions{
			OutputFileName: restic.DefaultOutputFileName,
		}
	)

	cmd := &cobra.Command{
		Use:               "update-status",
		Short:             "Update status of Repository, Backup/Restore Session",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.EnsureRequiredFlags(cmd, "namespace", "output-dir")

			config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
			if err != nil {
				return err
			}
			opt.KubeClient, err = kubernetes.NewForConfig(config)
			if err != nil {
				return err
			}
			opt.StashClient, err = cs.NewForConfig(config)
			if err != nil {
				return err
			}

			opt.Config = config
			if opt.BackupSession != "" {
				return opt.UpdateBackupStatusFromFile()
			}
			if opt.RestoreSession != "" {
				return opt.UpdateRestoreStatusFromFile()
			}
			return fmt.Errorf("respective BackupSession or RestoreSession is not specified")
		},
	}

	cmd.Flags().StringVar(&masterURL, "master", masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")

	cmd.Flags().StringVar(&opt.Namespace, "namespace", "default", "Namespace of Backup/Restore Session")
	cmd.Flags().StringVar(&opt.Repository, "repository", opt.Repository, "Name of the Repository")
	cmd.Flags().StringVar(&opt.TargetRef.Kind, "target-kind", "", "Kind of the target")
	cmd.Flags().StringVar(&opt.TargetRef.Name, "target-name", "", "Name of the target")
	cmd.Flags().StringVar(&opt.BackupSession, "backupsession", opt.BackupSession, "Name of the Backup Session")
	cmd.Flags().StringVar(&opt.RestoreSession, "restoresession", opt.RestoreSession, "Name of the Restore Session")
	cmd.Flags().StringVar(&opt.OutputDir, "output-dir", opt.OutputDir, "Directory where output.json file will be written (keep empty if you don't need to write output in file)")
	cmd.Flags().BoolVar(&opt.Metrics.Enabled, "metrics-enabled", opt.Metrics.Enabled, "Specify whether to export Prometheus metrics")
	cmd.Flags().StringVar(&opt.Metrics.PushgatewayURL, "metrics-pushgateway-url", opt.Metrics.PushgatewayURL, "Pushgateway URL where the metrics will be pushed")
	cmd.Flags().StringSliceVar(&opt.Metrics.Labels, "metrics-labels", opt.Metrics.Labels, "Labels to apply in exported metrics")
	cmd.Flags().StringVar(&opt.Metrics.JobName, "prom-job-name", StashDefaultMetricJob, "Metrics job name")

	return cmd
}
