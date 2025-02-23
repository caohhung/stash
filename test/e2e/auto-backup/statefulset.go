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

package auto_backup

import (
	"context"
	"fmt"

	"stash.appscode.dev/apimachinery/apis"
	"stash.appscode.dev/apimachinery/apis/stash/v1alpha1"
	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	"stash.appscode.dev/stash/pkg/eventer"
	"stash.appscode.dev/stash/test/e2e/framework"
	"stash.appscode.dev/stash/test/e2e/matcher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	store "kmodules.xyz/objectstore-api/api/v1"
)

var _ = Describe("Auto-Backup", func() {

	var f *framework.Invocation

	BeforeEach(func() {
		f = framework.NewInvocation()
	})

	JustAfterEach(func() {
		f.PrintDebugInfoOnFailure()
	})

	AfterEach(func() {
		err := f.CleanupTestResources()
		Expect(err).NotTo(HaveOccurred())
	})

	annotations := func(backupBlueprintName string) map[string]string {
		return map[string]string{
			v1beta1.KeyBackupBlueprint: backupBlueprintName,
			v1beta1.KeyTargetPaths:     framework.TestSourceDataTargetPath,
			v1beta1.KeyVolumeMounts:    framework.TestSourceVolumeAndMount,
		}
	}

	Context("StatefulSet", func() {

		Context("Success Case", func() {

			It("should success auto-backup for the StatefulSet", func() {
				// Create BackupBlueprint
				bb, err := f.CreateBackupBlueprintForWorkload(framework.WorkloadBackupBlueprint)
				Expect(err).NotTo(HaveOccurred())

				// Deploy a StatefulSet
				ss, err := f.DeployStatefulSet(framework.SourceStatefulSet, int32(1), framework.SourceVolume)
				Expect(err).NotTo(HaveOccurred())

				// Generate Sample Data
				_, err = f.GenerateSampleData(ss.ObjectMeta, apis.KindStatefulSet)
				Expect(err).NotTo(HaveOccurred())

				// Add and Ensure annotations to Target
				err = f.AddAutoBackupAnnotations(annotations(bb.Name), ss)
				Expect(err).NotTo(HaveOccurred())

				// ensure Repository and BackupConfiguration
				backupConfig, err := f.VerifyAutoBackupConfigured(ss.ObjectMeta, apis.KindStatefulSet)
				Expect(err).NotTo(HaveOccurred())

				// Take an Instant Backup of the Sample Data
				backupSession, err := f.TakeInstantBackup(backupConfig.ObjectMeta, v1beta1.BackupInvokerRef{
					Name: backupConfig.Name,
					Kind: v1beta1.ResourceKindBackupConfiguration,
				})
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that BackupSession has succeeded")
				completedBS, err := f.StashClient.StashV1beta1().BackupSessions(backupSession.Namespace).Get(context.TODO(), backupSession.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(completedBS.Status.Phase).Should(Equal(v1beta1.BackupSessionSucceeded))
			})
		})

		Context("Failure Case", func() {

			Context("Missing AutoBackup resource credential in BackupBlueprint", func() {
				It("should fail BackupSession for missing Backend credential", func() {
					// Create Secret for BackupBlueprint
					secret, err := f.CreateBackendSecretForMinio()
					Expect(err).NotTo(HaveOccurred())

					// Generate BackupBlueprint definition
					bb := f.BackupBlueprint(secret.Name)
					bb.Spec.Backend.S3 = &store.S3Spec{}
					By(fmt.Sprintf("Creating BackupBlueprint: %s", bb.Name))
					_, err = f.CreateBackupBlueprint(bb)
					Expect(err).NotTo(HaveOccurred())
					f.AppendToCleanupList(bb)

					// Deploy a StatefulSet
					ss, err := f.DeployStatefulSet(framework.SourceStatefulSet, int32(1), framework.SourceVolume)
					Expect(err).NotTo(HaveOccurred())

					// Generate Sample Data
					_, err = f.GenerateSampleData(ss.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Add and Ensure annotations to Target
					err = f.AddAutoBackupAnnotations(annotations(bb.Name), ss)
					Expect(err).NotTo(HaveOccurred())

					// ensure Repository and BackupConfiguration
					backupConfig, err := f.VerifyAutoBackupConfigured(ss.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Take an Instant Backup of the Sample Data
					backupSession, err := f.TakeInstantBackup(backupConfig.ObjectMeta, v1beta1.BackupInvokerRef{
						Name: backupConfig.Name,
						Kind: v1beta1.ResourceKindBackupConfiguration,
					})
					Expect(err).NotTo(HaveOccurred())

					By("Verifying that BackupSession has failed")
					completedBS, err := f.StashClient.StashV1beta1().BackupSessions(backupSession.Namespace).Get(context.TODO(), backupSession.Name, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					Expect(completedBS.Status.Phase).Should(Equal(v1beta1.BackupSessionFailed))
				})
				It("should fail BackupSession for missing RetentionPolicy", func() {
					// Create Storage Secret for Minio
					secret, err := f.CreateBackendSecretForMinio()
					Expect(err).NotTo(HaveOccurred())

					// Generate BackupBlueprint definition
					bb := f.BackupBlueprint(secret.Name)
					bb.Spec.RetentionPolicy = v1alpha1.RetentionPolicy{}
					By(fmt.Sprintf("Creating BackupBlueprint: %s", bb.Name))
					_, err = f.CreateBackupBlueprint(bb)
					Expect(err).NotTo(HaveOccurred())

					// Deploy a StatefulSet
					ss, err := f.DeployStatefulSet(framework.SourceStatefulSet, int32(1), framework.SourceVolume)
					Expect(err).NotTo(HaveOccurred())

					// Generate Sample Data
					_, err = f.GenerateSampleData(ss.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Add and Ensure annotations to Target
					err = f.AddAutoBackupAnnotations(annotations(bb.Name), ss)
					Expect(err).NotTo(HaveOccurred())

					// ensure Repository and BackupConfiguration
					backupConfig, err := f.VerifyAutoBackupConfigured(ss.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Take an Instant Backup of the Sample Data
					backupSession, err := f.TakeInstantBackup(backupConfig.ObjectMeta, v1beta1.BackupInvokerRef{
						Name: backupConfig.Name,
						Kind: v1beta1.ResourceKindBackupConfiguration,
					})
					Expect(err).NotTo(HaveOccurred())

					By("Verifying that BackupSession has failed")
					completedBS, err := f.StashClient.StashV1beta1().BackupSessions(backupSession.Namespace).Get(context.TODO(), backupSession.Name, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					Expect(completedBS.Status.Phase).Should(Equal(v1beta1.BackupSessionFailed))
				})
			})

			Context("Add inappropriate annotation to Target", func() {
				It("Should fail auto-backup for adding inappropriate BackupBlueprint annotation in StatefulSet", func() {
					// Create BackupBlueprint
					_, err := f.CreateBackupBlueprintForWorkload(framework.WorkloadBackupBlueprint)
					Expect(err).NotTo(HaveOccurred())

					// Deploy a StatefulSet
					ss, err := f.DeployStatefulSet(framework.SourceStatefulSet, int32(1), framework.SourceVolume)
					Expect(err).NotTo(HaveOccurred())

					// Generate Sample Data
					_, err = f.GenerateSampleData(ss.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Add and Ensure annotations to Target
					err = f.AddAutoBackupAnnotations(annotations(framework.WrongBackupBlueprintName), ss)
					Expect(err).NotTo(HaveOccurred())

					// AutoBackup Resource creation failed
					f.EventuallyEvent(ss.ObjectMeta, apis.KindStatefulSet).Should(matcher.HaveEvent(eventer.EventReasonAutoBackupResourcesCreationFailed))
				})
				It("should fail BackupSession for adding inappropriate TargetPath/MountPath StatefulSet", func() {
					// Create BackupBlueprint
					bb, err := f.CreateBackupBlueprintForWorkload(framework.WorkloadBackupBlueprint)
					Expect(err).NotTo(HaveOccurred())

					// Deploy a StatefulSet
					ss, err := f.DeployStatefulSet(framework.SourceStatefulSet, int32(1), framework.SourceVolume)
					Expect(err).NotTo(HaveOccurred())

					// Generate Sample Data
					_, err = f.GenerateSampleData(ss.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Add and Ensure annotations to Target
					anno := annotations(bb.Name)
					anno[v1beta1.KeyTargetPaths] = framework.WrongTargetPath
					err = f.AddAutoBackupAnnotations(anno, ss)
					Expect(err).NotTo(HaveOccurred())

					// ensure Repository and BackupConfiguration
					backupConfig, err := f.VerifyAutoBackupConfigured(ss.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Take an Instant Backup of the Sample Data
					backupSession, err := f.TakeInstantBackup(backupConfig.ObjectMeta, v1beta1.BackupInvokerRef{
						Name: backupConfig.Name,
						Kind: v1beta1.ResourceKindBackupConfiguration,
					})
					Expect(err).NotTo(HaveOccurred())

					By("Verifying that BackupSession has failed")
					completedBS, err := f.StashClient.StashV1beta1().BackupSessions(backupSession.Namespace).Get(context.TODO(), backupSession.Name, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					Expect(completedBS.Status.Phase).Should(Equal(v1beta1.BackupSessionFailed))
				})
			})
		})
	})
})
