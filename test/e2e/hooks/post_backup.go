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

package hooks

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"stash.appscode.dev/apimachinery/apis"
	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	"stash.appscode.dev/stash/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	pfutil "kmodules.xyz/client-go/tools/portforward"
	probev1 "kmodules.xyz/prober/api/v1"
)

var _ = Describe("PostBackup Hook", func() {

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

	Context("HTTPGetAction", func() {
		Context("Sidecar Model", func() {
			Context("Success Test", func() {
				Context("Host and Port specified", func() {
					It("should execute hook successfully", func() {
						// Deploy a StatefulSet with prober client. Here, we are using a StatefulSet because we need a stable address
						// for pod where http request will be sent.
						statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
						Expect(err).NotTo(HaveOccurred())

						// Generate Sample Data
						_, err = f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
						Expect(err).NotTo(HaveOccurred())

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup workload Backup
						backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									HTTPGet: &core.HTTPGetAction{
										Scheme: "HTTP",
										Host:   fmt.Sprintf("%s-0.%s.%s.svc", statefulset.Name, statefulset.Name, f.Namespace()),
										Path:   "/success",
										Port:   intstr.FromInt(framework.HttpPort),
									},
								},
							}
						})
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

				Context("Host and Port from Pod", func() {
					It("should execute hook successfully", func() {
						// Deploy a StatefulSet.
						statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
						Expect(err).NotTo(HaveOccurred())

						// Generate Sample Data
						_, err = f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
						Expect(err).NotTo(HaveOccurred())

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup backup
						backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									HTTPGet: &core.HTTPGetAction{
										Scheme: "HTTP",
										Path:   "/success",
										Port:   intstr.FromString(framework.HttpPortName),
									},
									ContainerName: framework.ProberDemoPodPrefix,
								},
							}
						})
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
			})

			Context("Failure Test", func() {
				It("should take a backup even when the postBackup hook failed", func() {
					// Deploy a StatefulSet.
					statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
					Expect(err).NotTo(HaveOccurred())

					// Generate Sample Data
					_, err = f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Setup a Minio Repository
					repo, err := f.SetupMinioRepository()
					Expect(err).NotTo(HaveOccurred())
					f.AppendToCleanupList(repo)

					// Setup Backup
					backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
						bc.Spec.Hooks = &v1beta1.BackupHooks{
							PostBackup: &probev1.Handler{
								HTTPGet: &core.HTTPGetAction{
									Scheme: "HTTP",
									Path:   "/fail",
									Port:   intstr.FromString(framework.HttpPortName),
								},
								ContainerName: framework.ProberDemoPodPrefix,
							},
						}
					})
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

					By("Verifying that a backup has been taken")
					items, err := f.BrowseMinioRepository(repo)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(items).ShouldNot(BeEmpty())
				})
			})
		})
	})

	Context("HTTPPostAction", func() {
		Context("Sidecar Model", func() {
			Context("Success Test", func() {
				Context("Host and Port specified", func() {
					It("should execute hook successfully", func() {
						// Deploy a StatefulSet with prober client. Here, we are using a StatefulSet because we need a stable address
						// for pod where http request will be sent.
						statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
						Expect(err).NotTo(HaveOccurred())

						// Generate Sample Data
						_, err = f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
						Expect(err).NotTo(HaveOccurred())

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup workload Backup
						backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									HTTPPost: &probev1.HTTPPostAction{
										Scheme: "HTTP",
										Host:   fmt.Sprintf("%s-0.%s.%s.svc", statefulset.Name, statefulset.Name, f.Namespace()),
										Path:   "/post-demo",
										Port:   intstr.FromInt(framework.HttpPort),
									},
								},
							}
						})
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

				Context("Host and Port from Pod", func() {
					It("should execute hook successfully", func() {
						// Deploy a StatefulSet.
						statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
						Expect(err).NotTo(HaveOccurred())

						// Generate Sample Data
						_, err = f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
						Expect(err).NotTo(HaveOccurred())

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup Backup
						backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									HTTPPost: &probev1.HTTPPostAction{
										Scheme: "HTTP",
										Path:   "/post-demo",
										Port:   intstr.FromString(framework.HttpPortName),
									},
									ContainerName: framework.ProberDemoPodPrefix,
								},
							}
						})
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

				Context("Json Data in Request Body", func() {
					It("server should echo the 'expectedCode' and 'expectedResponse' passed in the json body", func() {
						// Deploy a StatefulSet.
						statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
						Expect(err).NotTo(HaveOccurred())

						// Generate Sample Data
						_, err = f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
						Expect(err).NotTo(HaveOccurred())

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup Backup
						backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									HTTPPost: &probev1.HTTPPostAction{
										Scheme: "HTTP",
										Path:   "/post-demo",
										Port:   intstr.FromString(framework.HttpPortName),
										Body:   `{"expectedCode":"200","expectedResponse":"success"}`,
									},
									ContainerName: framework.ProberDemoPodPrefix,
								},
							}
						})
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

				Context("Form Data in Request Body", func() {
					It("server should echo the 'expectedCode' and 'expectedResponse' passed as form data", func() {
						// Deploy a StatefulSet.
						statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
						Expect(err).NotTo(HaveOccurred())

						// Generate Sample Data
						_, err = f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
						Expect(err).NotTo(HaveOccurred())

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup Backup
						backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									HTTPPost: &probev1.HTTPPostAction{
										Scheme: "HTTP",
										Path:   "/post-demo",
										Port:   intstr.FromString(framework.HttpPortName),
										Form: []probev1.FormEntry{
											{
												Key:    "expectedResponse",
												Values: []string{"success"},
											},
											{
												Key:    "expectedCode",
												Values: []string{"202"},
											},
										},
									},
									ContainerName: framework.ProberDemoPodPrefix,
								},
							}
						})
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
			})

			Context("Failure Test", func() {
				It("should take a backup even when the postBackup hook failed", func() {
					// Deploy a StatefulSet.
					statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
					Expect(err).NotTo(HaveOccurred())

					// Generate Sample Data
					_, err = f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Setup a Minio Repository
					repo, err := f.SetupMinioRepository()
					Expect(err).NotTo(HaveOccurred())
					f.AppendToCleanupList(repo)

					// Setup Backup
					backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
						bc.Spec.Hooks = &v1beta1.BackupHooks{
							PostBackup: &probev1.Handler{
								HTTPPost: &probev1.HTTPPostAction{
									Scheme: "HTTP",
									Path:   "/post-demo",
									Port:   intstr.FromString(framework.HttpPortName),
									Form: []probev1.FormEntry{
										{
											Key:    "expectedResponse",
											Values: []string{"fail"},
										},
										{
											Key:    "expectedCode",
											Values: []string{"403"},
										},
									},
								},
								ContainerName: framework.ProberDemoPodPrefix,
							},
						}
					})
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

					By("Verifying that a backup has been taken")
					items, err := f.BrowseMinioRepository(repo)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(items).ShouldNot(BeEmpty())
				})
			})
		})
	})

	Context("TCPSocketAction", func() {
		Context("Sidecar Model", func() {
			Context("Success Test", func() {
				Context("Host and Port specified", func() {
					It("should execute hook successfully", func() {
						// Deploy a StatefulSet with prober client. Here, we are using a StatefulSet because we need a stable address
						// for pod where http request will be sent.
						statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
						Expect(err).NotTo(HaveOccurred())

						// Generate Sample Data
						_, err = f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
						Expect(err).NotTo(HaveOccurred())

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup workload Backup
						backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									TCPSocket: &core.TCPSocketAction{
										Host: fmt.Sprintf("%s-0.%s.%s.svc", statefulset.Name, statefulset.Name, f.Namespace()),
										Port: intstr.FromInt(framework.TcpPort),
									},
								},
							}
						})
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

				Context("Host and Port from Pod", func() {
					It("should execute hook successfully", func() {
						// Deploy a StatefulSet.
						statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
						Expect(err).NotTo(HaveOccurred())

						// Generate Sample Data
						_, err = f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
						Expect(err).NotTo(HaveOccurred())

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup Backup
						backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									TCPSocket: &core.TCPSocketAction{
										Port: intstr.FromString(framework.TcpPortName),
									},
									ContainerName: framework.ProberDemoPodPrefix,
								},
							}
						})
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
			})

			Context("Failure Test", func() {
				It("should take a backup even when the postBackup hook failed", func() {
					// Deploy a StatefulSet.
					statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
					Expect(err).NotTo(HaveOccurred())

					// Generate Sample Data
					_, err = f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Setup a Minio Repository
					repo, err := f.SetupMinioRepository()
					Expect(err).NotTo(HaveOccurred())
					f.AppendToCleanupList(repo)

					// Setup Backup
					backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
						bc.Spec.Hooks = &v1beta1.BackupHooks{
							PostBackup: &probev1.Handler{
								TCPSocket: &core.TCPSocketAction{
									Port: intstr.FromInt(9091),
								},
								ContainerName: framework.ProberDemoPodPrefix,
							},
						}
					})
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

					By("Verifying that a backup has been taken")
					items, err := f.BrowseMinioRepository(repo)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(items).ShouldNot(BeEmpty())
				})
			})
		})
	})

	Context("ExecAction", func() {
		Context("Sidecar Model", func() {
			Context("Success Test", func() {
				It("should cleanup the sample data in postBackup hook", func() {
					// Deploy a StatefulSet.
					statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
					Expect(err).NotTo(HaveOccurred())

					// Read data at empty state
					emptyData, err := f.ReadSampleDataFromFromWorkload(statefulset.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Generate Sample Data
					sampleData, err := f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())
					Expect(sampleData).ShouldNot(BeEquivalentTo(emptyData))

					// Setup a Minio Repository
					repo, err := f.SetupMinioRepository()
					Expect(err).NotTo(HaveOccurred())
					f.AppendToCleanupList(repo)

					// Setup Backup
					backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
						bc.Spec.Hooks = &v1beta1.BackupHooks{
							PostBackup: &probev1.Handler{
								Exec: &core.ExecAction{
									Command: []string{"/bin/sh", "-c", fmt.Sprintf("rm -rf %s/*", framework.TestSourceDataMountPath)},
								},
								ContainerName: framework.ProberDemoPodPrefix,
							},
						}
					})
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

					By("Verifying that data has been removed in postBackup hook")
					newData, err := f.ReadSampleDataFromFromWorkload(statefulset.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())
					Expect(newData).Should(BeEquivalentTo(emptyData))
				})

				It("should execute postBackup hook even when the backup process failed", func() {
					// Deploy a StatefulSet.
					statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
					Expect(err).NotTo(HaveOccurred())

					// Read data at empty state
					emptyData, err := f.ReadSampleDataFromFromWorkload(statefulset.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Generate Sample Data
					sampleData, err := f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())
					Expect(sampleData).ShouldNot(BeEquivalentTo(emptyData))

					// Setup a Minio Repository
					repo, err := f.SetupMinioRepository()
					Expect(err).NotTo(HaveOccurred())
					f.AppendToCleanupList(repo)

					// Setup Backup
					// Use invalid retentionPolicy so that the backup process fail in cleanup step
					// Remove old data in postBackup hook
					backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
						bc.Spec.Hooks = &v1beta1.BackupHooks{
							PostBackup: &probev1.Handler{
								Exec: &core.ExecAction{
									Command: []string{"/bin/sh", "-c", fmt.Sprintf("rm -rf %s/*", framework.TestSourceDataMountPath)},
								},
								ContainerName: framework.ProberDemoPodPrefix,
							},
						}
						bc.Spec.RetentionPolicy.KeepLast = 0 // invalid retention value to force backup process fail on cleanup step
					})
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

					By("Verifying that data has been removed in postBackup hook")
					newData, err := f.ReadSampleDataFromFromWorkload(statefulset.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())
					Expect(newData).Should(BeEquivalentTo(emptyData))
				})
			})

			Context("Failure Test", func() {
				It("should take a backup even when the postBackup hook failed", func() {
					// Deploy a StatefulSet.
					statefulset, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
					Expect(err).NotTo(HaveOccurred())

					// Read data at empty state
					emptyData, err := f.ReadSampleDataFromFromWorkload(statefulset.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())

					// Generate Sample Data
					sampleData, err := f.GenerateSampleData(statefulset.ObjectMeta, apis.KindStatefulSet)
					Expect(err).NotTo(HaveOccurred())
					Expect(sampleData).ShouldNot(BeEquivalentTo(emptyData))

					// Setup a Minio Repository
					repo, err := f.SetupMinioRepository()
					Expect(err).NotTo(HaveOccurred())
					f.AppendToCleanupList(repo)

					// Setup Backup
					// Return non-zero exit code so that the postBackup hook fail
					backupConfig, err := f.SetupWorkloadBackup(statefulset.ObjectMeta, repo, apis.KindStatefulSet, func(bc *v1beta1.BackupConfiguration) {
						bc.Spec.Hooks = &v1beta1.BackupHooks{
							PostBackup: &probev1.Handler{
								Exec: &core.ExecAction{
									Command: []string{"/bin/sh", "-c", "exit 1"},
								},
								ContainerName: framework.ProberDemoPodPrefix,
							},
						}
					})
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

					By("Verifying that a backup has been taken")
					items, err := f.BrowseMinioRepository(repo)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(items).ShouldNot(BeEmpty())
				})
			})
		})

		Context("Job Model", func() {
			Context("PVC", func() {
				Context("Success Cases", func() {
					It("should cleanup the sample data in postBackup hook", func() {
						// Create new PVC
						pvc, err := f.CreateNewPVC(fmt.Sprintf("%s-%s", framework.SourceVolume, f.App()))
						Expect(err).NotTo(HaveOccurred())

						// Deploy a Pod
						pod, err := f.DeployPod(pvc.Name)
						Expect(err).NotTo(HaveOccurred())

						// Read data at empty state
						emptyData, err := f.ReadSampleDataFromFromWorkload(pod.ObjectMeta, apis.KindPod)
						Expect(err).NotTo(HaveOccurred())

						// Generate Sample Data
						sampleData, err := f.GenerateSampleData(pod.ObjectMeta, apis.KindPod)
						Expect(err).NotTo(HaveOccurred())
						Expect(sampleData).NotTo(BeEquivalentTo(emptyData))

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup PVC Backup
						// Remove old data in postBackup hook
						backupConfig, err := f.SetupPVCBackup(pvc, repo, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									Exec: &core.ExecAction{
										Command: []string{"/bin/sh", "-c", fmt.Sprintf("rm -rf %s/*", apis.StashDefaultMountPath)},
									},
									ContainerName: apis.PostTaskHook,
								},
							}
						})
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

						By("Verifying that data has been removed in postBackup hook")
						newData, err := f.ReadSampleDataFromFromWorkload(pod.ObjectMeta, apis.KindPod)
						Expect(err).NotTo(HaveOccurred())
						Expect(newData).Should(BeEquivalentTo(emptyData))
					})

					It("should execute postBackup hook even when the backup process failed", func() {
						// Create new PVC
						pvc, err := f.CreateNewPVC(fmt.Sprintf("%s-%s", framework.SourceVolume, f.App()))
						Expect(err).NotTo(HaveOccurred())

						// Deploy a Pod
						pod, err := f.DeployPod(pvc.Name)
						Expect(err).NotTo(HaveOccurred())

						// Read data at empty state
						emptyData, err := f.ReadSampleDataFromFromWorkload(pod.ObjectMeta, apis.KindPod)
						Expect(err).NotTo(HaveOccurred())

						// Generate Sample Data
						sampleData, err := f.GenerateSampleData(pod.ObjectMeta, apis.KindPod)
						Expect(err).NotTo(HaveOccurred())
						Expect(sampleData).NotTo(BeEquivalentTo(emptyData))

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup PVC Backup
						// Use invalid retentionPolicy so that the backup process fail in cleanup step
						// Remove old data in postBackup hook
						backupConfig, err := f.SetupPVCBackup(pvc, repo, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									Exec: &core.ExecAction{
										Command: []string{"/bin/sh", "-c", fmt.Sprintf("rm -rf %s/*", apis.StashDefaultMountPath)},
									},
									ContainerName: apis.PostTaskHook,
								},
							}
							bc.Spec.RetentionPolicy.KeepLast = 0 // invalid retention value to force backup process fail on cleanup step
						})
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

						By("Verifying that data has been removed in postBackup hook")
						newData, err := f.ReadSampleDataFromFromWorkload(pod.ObjectMeta, apis.KindPod)
						Expect(err).NotTo(HaveOccurred())
						Expect(newData).Should(BeEquivalentTo(emptyData))
					})
				})

				Context("Failure Cases", func() {
					It("should take backup even when the postBackup hook failed", func() {
						// Create new PVC
						pvc, err := f.CreateNewPVC(fmt.Sprintf("%s-%s", framework.SourceVolume, f.App()))
						Expect(err).NotTo(HaveOccurred())

						// Deploy a Pod
						pod, err := f.DeployPod(pvc.Name)
						Expect(err).NotTo(HaveOccurred())

						// Read data at empty state
						emptyData, err := f.ReadSampleDataFromFromWorkload(pod.ObjectMeta, apis.KindPod)
						Expect(err).NotTo(HaveOccurred())

						// Generate Sample Data
						sampleData, err := f.GenerateSampleData(pod.ObjectMeta, apis.KindPod)
						Expect(err).NotTo(HaveOccurred())
						Expect(sampleData).NotTo(BeEquivalentTo(emptyData))

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup PVC Backup
						// Return non-zero exit code from postBackup hook so that it fail
						backupConfig, err := f.SetupPVCBackup(pvc, repo, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									Exec: &core.ExecAction{
										Command: []string{"/bin/sh", "-c", "exit 1"},
									},
									ContainerName: apis.PostTaskHook,
								},
							}
						})
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

						By("Verifying that a backup has been taken")
						items, err := f.BrowseMinioRepository(repo)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(items).ShouldNot(BeEmpty())
					})
				})
			})

			Context("MySQL", func() {
				const (
					sampleTable = "StashDemo"
				)

				BeforeEach(func() {
					// Skip test if respective Functions and Tasks are not installed.
					if !f.MySQLAddonInstalled() {
						Skip("MySQL addon is not installed")
					}
				})

				Context("Success Test", func() {
					It("should make the database writable that was made read-only in preBackup hook", func() {
						// Deploy MySQL database and respective service,secret,PVC and AppBinding.
						By("Deploying MySQL Server")
						dpl, appBinding, err := f.DeployMySQLDatabase()
						Expect(err).NotTo(HaveOccurred())

						By("Port forwarding MySQL pod")
						pod, err := f.GetPod(dpl.ObjectMeta)
						Expect(err).NotTo(HaveOccurred())
						tunnel := pfutil.NewTunnel(f.KubeClient.CoreV1().RESTClient(), f.ClientConfig, pod.Namespace, pod.Name, framework.MySQLServingPortNumber)
						defer tunnel.Close()
						err = tunnel.ForwardPort()
						Expect(err).NotTo(HaveOccurred())

						By("Connecting with MySQL Server")
						connstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", framework.SuperUser, f.App(), framework.LocalHostIP, tunnel.Local)
						db, err := sql.Open("mysql", connstr)
						Expect(err).NotTo(HaveOccurred())
						defer db.Close()
						db.SetConnMaxLifetime(time.Second * 10)
						err = f.EventuallyConnectWithMySQLServer(db)
						Expect(err).NotTo(HaveOccurred())

						By("Creating Sample Table")
						err = f.CreateTable(db, sampleTable)
						Expect(err).NotTo(HaveOccurred())

						By("Verifying that the sample table has been created")
						tables, err := f.ListTables(db)
						Expect(err).NotTo(HaveOccurred())
						Expect(tables.Has(sampleTable)).Should(BeTrue())

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup Database Backup
						// Here, we are going to make the database read only in preBackup hook.
						// Then, we will make the database writable in postBackup hook.
						backupConfig, err := f.SetupDatabaseBackup(appBinding, repo, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PreBackup: &probev1.Handler{
									Exec: &core.ExecAction{
										Command: []string{"/bin/sh", "-c",
											`mysql -u root --password=$MYSQL_ROOT_PASSWORD -e "SET GLOBAL super_read_only = ON;"`},
									},
									ContainerName: framework.MySQLContainerName,
								},
								PostBackup: &probev1.Handler{
									Exec: &core.ExecAction{
										Command: []string{"/bin/sh", "-c",
											`mysql -u root --password=$MYSQL_ROOT_PASSWORD -e "SET GLOBAL super_read_only = OFF;"`},
									},
									ContainerName: framework.MySQLContainerName,
								},
							}
						})
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

						By("Verifying that the database is writable")
						err = f.CreateTable(db, "readOnlyTest")
						Expect(err).NotTo(HaveOccurred())
					})

					It("should execute postBackup hook even when the backup process failed", func() {
						// Deploy MySQL database and respective service,secret,PVC and AppBinding.
						By("Deploying MySQL Server")
						dpl, appBinding, err := f.DeployMySQLDatabase()
						Expect(err).NotTo(HaveOccurred())

						By("Port forwarding MySQL pod")
						pod, err := f.GetPod(dpl.ObjectMeta)
						Expect(err).NotTo(HaveOccurred())
						tunnel := pfutil.NewTunnel(f.KubeClient.CoreV1().RESTClient(), f.ClientConfig, pod.Namespace, pod.Name, framework.MySQLServingPortNumber)
						defer tunnel.Close()
						err = tunnel.ForwardPort()
						Expect(err).NotTo(HaveOccurred())

						By("Connecting with MySQL Server")
						connstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", framework.SuperUser, f.App(), framework.LocalHostIP, tunnel.Local)
						db, err := sql.Open("mysql", connstr)
						Expect(err).NotTo(HaveOccurred())
						defer db.Close()
						db.SetConnMaxLifetime(time.Second * 10)
						err = f.EventuallyConnectWithMySQLServer(db)
						Expect(err).NotTo(HaveOccurred())

						By("Creating Sample Table")
						err = f.CreateTable(db, sampleTable)
						Expect(err).NotTo(HaveOccurred())

						By("Verifying that the sample table has been created")
						tables, err := f.ListTables(db)
						Expect(err).NotTo(HaveOccurred())
						Expect(tables.Has(sampleTable)).Should(BeTrue())

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup Database Backup
						// Here, we are going to make the database read only in preBackup hook.
						// Then, we will make the database writable in postBackup hook.
						// Use invalid retentionPolicy so that the backup process fail in cleanup step
						backupConfig, err := f.SetupDatabaseBackup(appBinding, repo, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PreBackup: &probev1.Handler{
									Exec: &core.ExecAction{
										Command: []string{"/bin/sh", "-c",
											`mysql -u root --password=$MYSQL_ROOT_PASSWORD -e "SET GLOBAL super_read_only = ON;"`},
									},
									ContainerName: framework.MySQLContainerName,
								},
								PostBackup: &probev1.Handler{
									Exec: &core.ExecAction{
										Command: []string{"/bin/sh", "-c",
											`mysql -u root --password=$MYSQL_ROOT_PASSWORD -e "SET GLOBAL super_read_only = OFF;"`},
									},
									ContainerName: framework.MySQLContainerName,
								},
							}
							bc.Spec.RetentionPolicy.KeepLast = 0 // invalid retention value to force backup process fail on cleanup step
						})
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

						By("Verifying that the postBackup hook has been executed")
						err = f.CreateTable(db, "readOnlyTest")
						Expect(err).NotTo(HaveOccurred())
					})
				})

				Context("Failure Test", func() {
					It("should take backup even when postBackup hook failed", func() {
						// Deploy MySQL database and respective service,secret,PVC and AppBinding.
						By("Deploying MySQL Server")
						dpl, appBinding, err := f.DeployMySQLDatabase()
						Expect(err).NotTo(HaveOccurred())

						By("Port forwarding MySQL pod")
						pod, err := f.GetPod(dpl.ObjectMeta)
						Expect(err).NotTo(HaveOccurred())
						tunnel := pfutil.NewTunnel(f.KubeClient.CoreV1().RESTClient(), f.ClientConfig, pod.Namespace, pod.Name, framework.MySQLServingPortNumber)
						defer tunnel.Close()
						err = tunnel.ForwardPort()
						Expect(err).NotTo(HaveOccurred())

						By("Connecting with MySQL Server")
						connstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", framework.SuperUser, f.App(), framework.LocalHostIP, tunnel.Local)
						db, err := sql.Open("mysql", connstr)
						Expect(err).NotTo(HaveOccurred())
						defer db.Close()
						db.SetConnMaxLifetime(time.Second * 10)
						err = f.EventuallyConnectWithMySQLServer(db)
						Expect(err).NotTo(HaveOccurred())

						By("Creating Sample Table")
						err = f.CreateTable(db, sampleTable)
						Expect(err).NotTo(HaveOccurred())

						By("Verifying that the sample table has been created")
						tables, err := f.ListTables(db)
						Expect(err).NotTo(HaveOccurred())
						Expect(tables.Has(sampleTable)).Should(BeTrue())

						// Setup a Minio Repository
						repo, err := f.SetupMinioRepository()
						Expect(err).NotTo(HaveOccurred())
						f.AppendToCleanupList(repo)

						// Setup Database Backup
						// Return non-zero exit code from the postBackup hook so that it fail.
						backupConfig, err := f.SetupDatabaseBackup(appBinding, repo, func(bc *v1beta1.BackupConfiguration) {
							bc.Spec.Hooks = &v1beta1.BackupHooks{
								PostBackup: &probev1.Handler{
									Exec: &core.ExecAction{
										Command: []string{"/bin/sh", "-c", "exit 1"},
									},
									ContainerName: framework.MySQLContainerName,
								},
							}
						})
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

						By("Verifying that a backup has been taken")
						items, err := f.BrowseMinioRepository(repo)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(items).ShouldNot(BeEmpty())
					})
				})
			})
		})
	})

	Context("Batch Backup", func() {
		Context("HTTPGetAction", func() {
			It("should execute global and local hooks successfully", func() {
				// Here, we are going to deploy two different StatefulSet with probe client.
				// Then, we are going to backup those StatefulSets using BatchBackup.
				// Each individual StatefulSet will have a hook for them.
				// The BackupBatch object will have a global hook.
				var members []v1beta1.BackupConfigurationTemplateSpec

				// Deploy the first StatefulSet.
				ss1, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix + "-1")
				Expect(err).NotTo(HaveOccurred())
				// Generate Sample Data
				_, err = f.GenerateSampleData(ss1.ObjectMeta, apis.KindStatefulSet)
				Expect(err).NotTo(HaveOccurred())
				// We will execute HTTPGetAction in the first StatefulSet
				members = append(members, v1beta1.BackupConfigurationTemplateSpec{
					Hooks: &v1beta1.BackupHooks{
						PostBackup: &probev1.Handler{
							HTTPGet: &core.HTTPGetAction{
								Scheme: "HTTP",
								Host:   fmt.Sprintf("%s-0.%s.%s.svc", ss1.Name, ss1.Name, f.Namespace()),
								Path:   "/success",
								Port:   intstr.FromInt(framework.HttpPort),
							},
						},
					},
					Target: &v1beta1.BackupTarget{
						Ref: v1beta1.TargetRef{
							APIVersion: appsv1.SchemeGroupVersion.String(),
							Kind:       apis.KindStatefulSet,
							Name:       ss1.Name,
						},
						Paths: []string{
							framework.TestSourceDataMountPath,
						},
						VolumeMounts: []core.VolumeMount{
							{
								Name:      framework.SourceVolume,
								MountPath: framework.TestSourceDataMountPath,
							},
						},
					},
				})

				// Deploy second StatefulSet
				ss2, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix + "-2")
				Expect(err).NotTo(HaveOccurred())
				// Generate Sample Data
				_, err = f.GenerateSampleData(ss2.ObjectMeta, apis.KindStatefulSet)
				Expect(err).NotTo(HaveOccurred())
				// We will execute HTTPPostAction in the second StatefulSet
				members = append(members, v1beta1.BackupConfigurationTemplateSpec{
					Hooks: &v1beta1.BackupHooks{
						PostBackup: &probev1.Handler{
							HTTPPost: &probev1.HTTPPostAction{
								Scheme: "HTTP",
								Host:   fmt.Sprintf("%s-0.%s.%s.svc", ss2.Name, ss2.Name, f.Namespace()),
								Path:   "/post-demo",
								Port:   intstr.FromInt(framework.HttpPort),
							},
						},
					},
					Target: &v1beta1.BackupTarget{
						Ref: v1beta1.TargetRef{
							APIVersion: appsv1.SchemeGroupVersion.String(),
							Kind:       apis.KindStatefulSet,
							Name:       ss2.Name,
						},
						Paths: []string{
							framework.TestSourceDataMountPath,
						},
						VolumeMounts: []core.VolumeMount{
							{
								Name:      framework.SourceVolume,
								MountPath: framework.TestSourceDataMountPath,
							},
						},
					},
				})

				// Setup a Minio Repository
				repo, err := f.SetupMinioRepository()
				Expect(err).NotTo(HaveOccurred())
				f.AppendToCleanupList(repo)

				// Setup Batch Backup
				backupBatch, err := f.SetupBatchBackup(repo, func(in *v1beta1.BackupBatch) {
					in.Spec.Members = members
					in.Spec.Hooks = &v1beta1.BackupHooks{
						// Just simply execute a http probe in the first StatefulSet.
						// Although it does not represent the actual use case, but it probes that the global are working.
						PostBackup: &probev1.Handler{
							HTTPGet: &core.HTTPGetAction{
								Scheme: "HTTP",
								Host:   fmt.Sprintf("%s-0.%s.%s.svc", ss1.Name, ss1.Name, f.Namespace()),
								Path:   "/success",
								Port:   intstr.FromInt(framework.HttpPort),
							},
						},
					}
				})
				Expect(err).NotTo(HaveOccurred())

				// Take an Instant Backup the Sample Data
				backupSession, err := f.TakeInstantBackup(backupBatch.ObjectMeta, v1beta1.BackupInvokerRef{
					Name: backupBatch.Name,
					Kind: v1beta1.ResourceKindBackupBatch,
				})
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that BackupSession has succeeded")
				completedBS, err := f.StashClient.StashV1beta1().BackupSessions(backupSession.Namespace).Get(context.TODO(), backupSession.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(completedBS.Status.Phase).Should(Equal(v1beta1.BackupSessionSucceeded))
			})
		})

		Context("ExecAction", func() {
			It("should execute global and local hooks successfully", func() {
				// Here, we are going to deploy two different types of workload.
				// First workload is a StatefulSet with probe client. We will execute a simple http hook there. The other workload is a database.
				// We will make the database readonly in local hook. We will execute a simple exec action in global hook.
				var members []v1beta1.BackupConfigurationTemplateSpec

				// Deploy the first StatefulSet.
				ss1, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
				Expect(err).NotTo(HaveOccurred())
				// Generate Sample Data
				_, err = f.GenerateSampleData(ss1.ObjectMeta, apis.KindStatefulSet)
				Expect(err).NotTo(HaveOccurred())
				// We will execute HTTPGetAction in the first StatefulSet
				members = append(members, v1beta1.BackupConfigurationTemplateSpec{
					Hooks: &v1beta1.BackupHooks{
						PostBackup: &probev1.Handler{
							Exec: &core.ExecAction{
								Command: []string{"/bin/sh", "-c", `exit $EXIT_CODE_SUCCESS`},
							},
						},
					},
					Target: &v1beta1.BackupTarget{
						Ref: v1beta1.TargetRef{
							APIVersion: appsv1.SchemeGroupVersion.String(),
							Kind:       apis.KindStatefulSet,
							Name:       ss1.Name,
						},
						Paths: []string{
							framework.TestSourceDataMountPath,
						},
						VolumeMounts: []core.VolumeMount{
							{
								Name:      framework.SourceVolume,
								MountPath: framework.TestSourceDataMountPath,
							},
						},
					},
				})

				// Deploy MySQL database and respective service,secret,PVC and AppBinding.
				By("Deploying MySQL Server")
				dpl, appBinding, err := f.DeployMySQLDatabase()
				Expect(err).NotTo(HaveOccurred())

				By("Port forwarding MySQL pod")
				pod, err := f.GetPod(dpl.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())
				tunnel := pfutil.NewTunnel(f.KubeClient.CoreV1().RESTClient(), f.ClientConfig, pod.Namespace, pod.Name, framework.MySQLServingPortNumber)
				defer tunnel.Close()
				err = tunnel.ForwardPort()
				Expect(err).NotTo(HaveOccurred())

				By("Connecting with MySQL Server")
				connstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", framework.SuperUser, f.App(), framework.LocalHostIP, tunnel.Local)
				db, err := sql.Open("mysql", connstr)
				Expect(err).NotTo(HaveOccurred())
				defer db.Close()
				db.SetConnMaxLifetime(time.Second * 10)
				err = f.EventuallyConnectWithMySQLServer(db)
				Expect(err).NotTo(HaveOccurred())

				// add the database as member of batch backup
				members = append(members, v1beta1.BackupConfigurationTemplateSpec{
					// make the database readonly in preBackup hook then make the database writable again in postBackup hook.
					Hooks: &v1beta1.BackupHooks{
						PreBackup: &probev1.Handler{
							Exec: &core.ExecAction{
								Command: []string{"/bin/sh", "-c",
									`mysql -u root --password=$MYSQL_ROOT_PASSWORD -e "SET GLOBAL super_read_only = ON;"`},
							},
							ContainerName: framework.MySQLContainerName,
						},
						PostBackup: &probev1.Handler{
							Exec: &core.ExecAction{
								Command: []string{"/bin/sh", "-c",
									`mysql -u root --password=$MYSQL_ROOT_PASSWORD -e "SET GLOBAL super_read_only = OFF;"`},
							},
							ContainerName: framework.MySQLContainerName,
						},
					},
					Task: v1beta1.TaskRef{
						Name: framework.MySQLBackupTask,
					},
					Target: &v1beta1.BackupTarget{
						Ref: v1beta1.TargetRef{
							APIVersion: core.SchemeGroupVersion.String(),
							Kind:       apis.KindAppBinding,
							Name:       appBinding.Name,
						},
					},
				})

				// Setup a Minio Repository
				repo, err := f.SetupMinioRepository()
				Expect(err).NotTo(HaveOccurred())
				f.AppendToCleanupList(repo)

				// Setup Batch Backup
				backupBatch, err := f.SetupBatchBackup(repo, func(in *v1beta1.BackupBatch) {
					in.Spec.Members = members
					// Execute a simple exec hook in the global hook. This hook will be executed inside Stash operator.
					// Currently, we don't have any known command that really make sense to execute inside operator.
					// So, we are using the simplest command to test the global hook. :P
					in.Spec.Hooks = &v1beta1.BackupHooks{
						PostBackup: &probev1.Handler{
							Exec: &core.ExecAction{
								Command: []string{"/bin/sh", "-c", "exit 0"},
							},
							ContainerName: "operator",
						},
					}
				})
				Expect(err).NotTo(HaveOccurred())

				// Take an Instant Backup the Sample Data
				backupSession, err := f.TakeInstantBackup(backupBatch.ObjectMeta, v1beta1.BackupInvokerRef{
					Name: backupBatch.Name,
					Kind: v1beta1.ResourceKindBackupBatch,
				})
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that BackupSession has succeeded")
				completedBS, err := f.StashClient.StashV1beta1().BackupSessions(backupSession.Namespace).Get(context.TODO(), backupSession.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(completedBS.Status.Phase).Should(Equal(v1beta1.BackupSessionSucceeded))

				By("Verifying that the database is writable")
				err = f.CreateTable(db, "readOnlyTest")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Different Situations", func() {
			It("should backup even when the global postBackup hook failed", func() {
				// Here, we are going to deploy two different types of workload.
				// First workload is a StatefulSet with probe client. We will execute a simple http hook there. The other workload is a database.
				// We will make the database readonly in local hook.
				// We will execute a simple exec action in global hook. This global hook will fail.
				var members []v1beta1.BackupConfigurationTemplateSpec

				// Deploy the first StatefulSet.
				ss1, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
				Expect(err).NotTo(HaveOccurred())
				// Generate Sample Data
				_, err = f.GenerateSampleData(ss1.ObjectMeta, apis.KindStatefulSet)
				Expect(err).NotTo(HaveOccurred())
				// We will execute HTTPGetAction in the first StatefulSet
				members = append(members, v1beta1.BackupConfigurationTemplateSpec{
					Hooks: &v1beta1.BackupHooks{
						PostBackup: &probev1.Handler{
							HTTPGet: &core.HTTPGetAction{
								Scheme: "HTTP",
								Host:   fmt.Sprintf("%s-0.%s.%s.svc", ss1.Name, ss1.Name, f.Namespace()),
								Path:   "/success",
								Port:   intstr.FromInt(framework.HttpPort),
							},
						},
					},
					Target: &v1beta1.BackupTarget{
						Ref: v1beta1.TargetRef{
							APIVersion: appsv1.SchemeGroupVersion.String(),
							Kind:       apis.KindStatefulSet,
							Name:       ss1.Name,
						},
						Paths: []string{
							framework.TestSourceDataMountPath,
						},
						VolumeMounts: []core.VolumeMount{
							{
								Name:      framework.SourceVolume,
								MountPath: framework.TestSourceDataMountPath,
							},
						},
					},
				})

				// Deploy MySQL database and respective service,secret,PVC and AppBinding.
				By("Deploying MySQL Server")
				dpl, appBinding, err := f.DeployMySQLDatabase()
				Expect(err).NotTo(HaveOccurred())

				By("Port forwarding MySQL pod")
				pod, err := f.GetPod(dpl.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())
				tunnel := pfutil.NewTunnel(f.KubeClient.CoreV1().RESTClient(), f.ClientConfig, pod.Namespace, pod.Name, framework.MySQLServingPortNumber)
				defer tunnel.Close()
				err = tunnel.ForwardPort()
				Expect(err).NotTo(HaveOccurred())

				By("Connecting with MySQL Server")
				connstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", framework.SuperUser, f.App(), framework.LocalHostIP, tunnel.Local)
				db, err := sql.Open("mysql", connstr)
				Expect(err).NotTo(HaveOccurred())
				defer db.Close()
				db.SetConnMaxLifetime(time.Second * 10)
				err = f.EventuallyConnectWithMySQLServer(db)
				Expect(err).NotTo(HaveOccurred())

				// add the database as member of batch backup
				members = append(members, v1beta1.BackupConfigurationTemplateSpec{
					// make the database readonly in preBackup hook.
					Hooks: &v1beta1.BackupHooks{
						PreBackup: &probev1.Handler{
							Exec: &core.ExecAction{
								Command: []string{"/bin/sh", "-c",
									`mysql -u root --password=$MYSQL_ROOT_PASSWORD -e "SET GLOBAL super_read_only = ON;"`},
							},
							ContainerName: framework.MySQLContainerName,
						},
					},
					Task: v1beta1.TaskRef{
						Name: framework.MySQLBackupTask,
					},
					Target: &v1beta1.BackupTarget{
						Ref: v1beta1.TargetRef{
							APIVersion: core.SchemeGroupVersion.String(),
							Kind:       apis.KindAppBinding,
							Name:       appBinding.Name,
						},
					},
				})

				// Setup a Minio Repository
				repo, err := f.SetupMinioRepository()
				Expect(err).NotTo(HaveOccurred())
				f.AppendToCleanupList(repo)

				// Setup Batch Backup
				backupBatch, err := f.SetupBatchBackup(repo, func(in *v1beta1.BackupBatch) {
					in.Spec.Members = members
					// intentionally fail the global postBackup hook
					in.Spec.Hooks = &v1beta1.BackupHooks{
						PostBackup: &probev1.Handler{
							Exec: &core.ExecAction{
								Command: []string{"/bin/sh", "-c", "exit 1"},
							},
							ContainerName: "operator",
						},
					}
				})
				Expect(err).NotTo(HaveOccurred())

				// Take an Instant Backup the Sample Data
				backupSession, err := f.TakeInstantBackup(backupBatch.ObjectMeta, v1beta1.BackupInvokerRef{
					Name: backupBatch.Name,
					Kind: v1beta1.ResourceKindBackupBatch,
				})
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that BackupSession has failed")
				completedBS, err := f.StashClient.StashV1beta1().BackupSessions(backupSession.Namespace).Get(context.TODO(), backupSession.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(completedBS.Status.Phase).Should(Equal(v1beta1.BackupSessionFailed))

				By("Verifying that the database is read-only") // this ensure that the local preBackup hook has been executed.
				err = f.CreateTable(db, "readOnlyTest")
				Expect(err).To(HaveOccurred())

				By("Verifying that backup has been taken in the backend")
				items, err := f.BrowseMinioRepository(repo)
				Expect(err).NotTo(HaveOccurred())
				Expect(items).ShouldNot(BeEmpty())

				By("Verifying that the global postBackup hook is the reason for BackupSession failure")
				backupSession, err = f.GetBackupSession(backupSession.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())

				allMembersSucceeded := true
				for _, target := range backupSession.Status.Targets {
					if target.Phase != v1beta1.TargetBackupSucceeded {
						allMembersSucceeded = false
					}
				}
				Expect(allMembersSucceeded).To(BeTrue())
			})

			It("should execute global postBackup hook even when a member failed", func() {
				// Here, we are going to deploy two different types of workload. We won't use any hook to those workloads.
				// The first workload is a StatefulSet with probe client. The another is a database.
				// We will intentionally fail the StatefulSet backup. As a result  the overall BackupSession will fail.
				// We will  execute a simple exec action in global hook.
				// Again, we will intentionally fail the global hook because if it does not fail we won't be able to detect whether it was executed or not.
				var members []v1beta1.BackupConfigurationTemplateSpec

				// Deploy the first StatefulSet.
				ss1, err := f.DeployStatefulSetWithProbeClient(framework.ProberDemoPodPrefix)
				Expect(err).NotTo(HaveOccurred())
				// Generate Sample Data
				_, err = f.GenerateSampleData(ss1.ObjectMeta, apis.KindStatefulSet)
				Expect(err).NotTo(HaveOccurred())
				// We will execute HTTPGetAction in the first StatefulSet
				members = append(members, v1beta1.BackupConfigurationTemplateSpec{
					Target: &v1beta1.BackupTarget{
						Ref: v1beta1.TargetRef{
							APIVersion: appsv1.SchemeGroupVersion.String(),
							Kind:       apis.KindStatefulSet,
							Name:       ss1.Name,
						},
						Paths: []string{
							"/invalid/path", // This path does not exist. Hence, the backup will fail.
						},
						VolumeMounts: []core.VolumeMount{
							{
								Name:      framework.SourceVolume,
								MountPath: framework.TestSourceDataMountPath,
							},
						},
					},
				})

				// Deploy MySQL database and respective service,secret,PVC and AppBinding.
				By("Deploying MySQL Server")
				dpl, appBinding, err := f.DeployMySQLDatabase()
				Expect(err).NotTo(HaveOccurred())

				By("Port forwarding MySQL pod")
				pod, err := f.GetPod(dpl.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())
				tunnel := pfutil.NewTunnel(f.KubeClient.CoreV1().RESTClient(), f.ClientConfig, pod.Namespace, pod.Name, framework.MySQLServingPortNumber)
				defer tunnel.Close()
				err = tunnel.ForwardPort()
				Expect(err).NotTo(HaveOccurred())

				By("Connecting with MySQL Server")
				connstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", framework.SuperUser, f.App(), framework.LocalHostIP, tunnel.Local)
				db, err := sql.Open("mysql", connstr)
				Expect(err).NotTo(HaveOccurred())
				defer db.Close()
				db.SetConnMaxLifetime(time.Second * 10)
				err = f.EventuallyConnectWithMySQLServer(db)
				Expect(err).NotTo(HaveOccurred())

				// add the database as member of batch backup
				members = append(members, v1beta1.BackupConfigurationTemplateSpec{
					Task: v1beta1.TaskRef{
						Name: framework.MySQLBackupTask,
					},
					Target: &v1beta1.BackupTarget{
						Ref: v1beta1.TargetRef{
							APIVersion: core.SchemeGroupVersion.String(),
							Kind:       apis.KindAppBinding,
							Name:       appBinding.Name,
						},
					},
				})

				// Setup a Minio Repository
				repo, err := f.SetupMinioRepository()
				Expect(err).NotTo(HaveOccurred())
				f.AppendToCleanupList(repo)

				// Setup Batch Backup
				backupBatch, err := f.SetupBatchBackup(repo, func(in *v1beta1.BackupBatch) {
					in.Spec.Members = members
					// intentionally fail the global postBackup hook
					in.Spec.Hooks = &v1beta1.BackupHooks{
						PostBackup: &probev1.Handler{
							Exec: &core.ExecAction{
								Command: []string{"/bin/sh", "-c", "exit 1"},
							},
							ContainerName: "operator",
						},
					}
				})
				Expect(err).NotTo(HaveOccurred())

				// Take an Instant Backup the Sample Data
				backupSession, err := f.TakeInstantBackup(backupBatch.ObjectMeta, v1beta1.BackupInvokerRef{
					Name: backupBatch.Name,
					Kind: v1beta1.ResourceKindBackupBatch,
				})
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that BackupSession has failed")
				completedBS, err := f.StashClient.StashV1beta1().BackupSessions(backupSession.Namespace).Get(context.TODO(), backupSession.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(completedBS.Status.Phase).Should(Equal(v1beta1.BackupSessionFailed))

				By("Verifying that at least one member has failed to take backup")
				backupSession, err = f.GetBackupSession(backupSession.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())

				memberFailed := false
				for _, target := range backupSession.Status.Targets {
					if target.Phase == v1beta1.TargetBackupFailed {
						memberFailed = true
					}
				}
				Expect(memberFailed).To(BeTrue())

				By("Verifying that the global postBackup hook has been executed")
				hookFailed, err := f.HookFailed(v1beta1.ResourceKindBackupSession, backupSession.ObjectMeta, "exec")
				Expect(err).NotTo(HaveOccurred())
				Expect(hookFailed).Should(BeTrue())
			})
		})
	})
})
