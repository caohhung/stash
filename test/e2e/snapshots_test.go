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

package e2e_test

//import (
//	"context"
//	"net"
//	"os/exec"
//	"strconv"
//
//	"stash.appscode.dev/apimachinery/apis"
//	api "stash.appscode.dev/apimachinery/apis/stash/v1alpha1"
//	"stash.appscode.dev/stash/test/e2e/framework"
//	. "stash.appscode.dev/stash/test/e2e/matcher"
//
//	"github.com/appscode/go/crypto/rand"
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	apps "k8s.io/api/apps/v1"
//	core "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	core_util "kmodules.xyz/client-go/core/v1"
//)
//
//var _ = XDescribe("Snapshots", func() {
//	var (
//		err        error
//		f          *framework.Invocation
//		restic     api.Restic
//		cred       core.Secret
//		deployment apps.Deployment
//		daemon     apps.DaemonSet
//		rc         core.ReplicationController
//		rs         apps.ReplicaSet
//		ss         apps.StatefulSet
//		svc        core.Service
//		workload   api.LocalTypedReference
//	)
//
//	BeforeEach(func() {
//		f = root.Invoke()
//	})
//	AfterEach(func() {
//		err = f.DeleteDaemonSet(daemon.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = f.DeleteDeployment(deployment.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = f.DeleteReplicationController(rc.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = f.DeleteReplicaSet(rs.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = f.DeleteService(svc.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = f.DeleteStatefulSet(ss.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = f.DeleteRestic(restic.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = f.DeleteSecret(cred.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		f.DeleteRepositories(f.DaemonSetRepos(&daemon))
//		f.DeleteRepositories(f.DeploymentRepos(&deployment))
//		f.DeleteRepositories(f.ReplicationControllerRepos(&rc))
//		f.DeleteRepositories(f.ReplicaSetRepos(&rs))
//		f.DeleteRepositories(f.StatefulSetRepos(&ss))
//
//		err := framework.WaitUntilDaemonSetDeleted(f.KubeClient, daemon.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = framework.WaitUntilRepositoriesDeleted(f.StashClient, f.DaemonSetRepos(&daemon))
//		Expect(err).NotTo(HaveOccurred())
//
//		err = framework.WaitUntilDeploymentDeleted(f.KubeClient, deployment.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = framework.WaitUntilRepositoriesDeleted(f.StashClient, f.DeploymentRepos(&deployment))
//		Expect(err).NotTo(HaveOccurred())
//
//		err = framework.WaitUntilReplicationControllerDeleted(f.KubeClient, rc.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = framework.WaitUntilRepositoriesDeleted(f.StashClient, f.ReplicationControllerRepos(&rc))
//		Expect(err).NotTo(HaveOccurred())
//
//		err = framework.WaitUntilReplicaSetDeleted(f.KubeClient, rs.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = framework.WaitUntilRepositoriesDeleted(f.StashClient, f.ReplicaSetRepos(&rs))
//		Expect(err).NotTo(HaveOccurred())
//
//		err = framework.WaitUntilStatefulSetDeleted(f.KubeClient, ss.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = framework.WaitUntilServiceDeleted(f.KubeClient, svc.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//		err = framework.WaitUntilRepositoriesDeleted(f.StashClient, f.StatefulSetRepos(&ss))
//		Expect(err).NotTo(HaveOccurred())
//
//		err = framework.WaitUntilSecretDeleted(f.KubeClient, cred.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//
//		err = framework.WaitUntilResticDeleted(f.StashClient, restic.ObjectMeta)
//		Expect(err).NotTo(HaveOccurred())
//
//		err = core_util.WaitUntillPodTerminatedByLabel(context.TODO(), f.KubeClient, ss.Namespace, f.AppLabel())
//		Expect(err).NotTo(HaveOccurred())
//	})
//
//	JustBeforeEach(func() {
//		if missing, _ := BeZero().Match(cred); missing {
//			Skip("Missing repository credential")
//		}
//		restic.Spec.Backend.StorageSecretName = cred.Name
//		pvc := f.PersistentVolumeClaim(rand.WithUniqSuffix("pvc"))
//		_, err := f.CreatePersistentVolumeClaim(pvc)
//		Expect(err).NotTo(HaveOccurred())
//		daemon = f.DaemonSet(framework.SourceDaemonSet, framework.SourceVolume)
//
//		deployment = f.Deployment(framework.SourceDeployment, pvc.Name, framework.SourceVolume)
//		rc = f.ReplicationController(framework.SourceReplicationController, pvc.Name, framework.SourceVolume)
//		rs = f.ReplicaSet(framework.SourceReplicaSet, pvc.Name, framework.SourceVolume)
//		ss = f.StatefulSet(framework.SourceStatefulSet, pvc.Name, framework.SourceVolume)
//		svc = f.HeadlessService(ss.Name)
//
//		// if a deployment's labels match to labels of replicaset, kubernetes make the deployment owner of the replicaset.
//		// avoid this adding extra label to deployment.
//		deployment.Labels["tier"] = "test"
//		deployment.Spec.Template.Labels["tier"] = "test"
//	})
//
//	var (
//		shouldCreateWorkloads = func() {
//			By("Creating DaemonSet " + daemon.Name)
//			_, err = f.CreateDaemonSet(daemon)
//			Expect(err).NotTo(HaveOccurred())
//
//			By("Creating Deployment " + deployment.Name)
//			_, err = f.CreateDeployment(deployment)
//			Expect(err).NotTo(HaveOccurred())
//
//			By("Creating ReplicationController " + rc.Name)
//			_, err = f.CreateReplicationController(rc)
//			Expect(err).NotTo(HaveOccurred())
//
//			By("Creating ReplicaSet " + rs.Name)
//			_, err = f.CreateReplicaSet(rs)
//			Expect(err).NotTo(HaveOccurred())
//
//			By("Creating service " + svc.Name)
//			_, err = f.CreateService(svc)
//			Expect(err).NotTo(HaveOccurred())
//
//			By("Creating StatefulSet " + ss.Name)
//			_, err = f.CreateStatefulSet(ss)
//			Expect(err).NotTo(HaveOccurred())
//		}
//
//		shouldHaveSidecar = func() {
//			By("Waiting for sidecar of DaemonSet")
//			f.EventuallyDaemonSet(daemon.ObjectMeta).Should(HaveSidecar(apis.StashContainer))
//
//			By("Waiting for sidecar of Deployment")
//			f.EventuallyDeployment(deployment.ObjectMeta).Should(HaveSidecar(apis.StashContainer))
//
//			By("Waiting for sidecar of ReplicationController")
//			f.EventuallyReplicationController(rc.ObjectMeta).Should(HaveSidecar(apis.StashContainer))
//
//			By("Waiting for sidecar of ReplicaSet")
//			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveSidecar(apis.StashContainer))
//
//			By("Waiting for sidecar of StatefulSet")
//			f.EventuallyStatefulSet(ss.ObjectMeta).Should(HaveSidecar(apis.StashContainer))
//		}
//
//		shouldHaveRepositoryCRD = func() {
//			By("Waiting for Repository CRD for DaemonSet")
//			f.EventuallyRepository(&daemon).ShouldNot(BeEmpty())
//
//			By("Waiting for Repository CRD for Deployment")
//			f.EventuallyRepository(&deployment).ShouldNot(BeEmpty())
//
//			By("Waiting for Repository CRD for ReplicationController")
//			f.EventuallyRepository(&rc).ShouldNot(BeEmpty())
//
//			By("Waiting for Repository CRD for ReplicaSet")
//			f.EventuallyRepository(&rs).ShouldNot(BeEmpty())
//
//			By("Waiting for Repository CRD for StatefulSet")
//			f.EventuallyRepository(&ss).Should(WithTransform(func(repoList []*api.Repository) int {
//				return len(repoList)
//			}, BeNumerically("==", int(*ss.Spec.Replicas))))
//		}
//
//		shouldBackupComplete = func() {
//			By("Waiting for backup to complete for DaemonsSet")
//			f.EventuallyRepository(&daemon).Should(WithTransform(f.BackupCountInRepositoriesStatus, BeNumerically(">=", 2)))
//
//			By("Waiting for backup to complete for Deployment")
//			f.EventuallyRepository(&deployment).Should(WithTransform(f.BackupCountInRepositoriesStatus, BeNumerically(">=", 2)))
//
//			By("Waiting for backup to complete for ReplicationController")
//			f.EventuallyRepository(&rc).Should(WithTransform(f.BackupCountInRepositoriesStatus, BeNumerically(">=", 2)))
//
//			By("Waiting for backup to complete for ReplicaSet")
//			f.EventuallyRepository(&rs).Should(WithTransform(f.BackupCountInRepositoriesStatus, BeNumerically(">=", 2)))
//
//			By("Waiting for backup to complete for StatefulSet")
//			f.EventuallyRepository(&ss).Should(WithTransform(f.BackupCountInRepositoriesStatus, BeNumerically(">=", 2)))
//		}
//
//		performOperationOnSnapshot = func() {
//			By("Creating repository Secret " + cred.Name)
//			_, err = f.CreateSecret(cred)
//			Expect(err).NotTo(HaveOccurred())
//
//			By("Creating restic")
//			err = f.CreateRestic(restic)
//			Expect(err).NotTo(HaveOccurred())
//
//			By("Creating workloads")
//			shouldCreateWorkloads()
//
//			By("Waiting for workloads to have sidecar")
//			shouldHaveSidecar()
//
//			By("Waiting for Repository CRD")
//			shouldHaveRepositoryCRD()
//
//			By("Waiting for backup to complete")
//			shouldBackupComplete()
//
//			By("Listing all snapshots")
//			_, err := f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).List(context.TODO(), metav1.ListOptions{})
//			Expect(err).NotTo(HaveOccurred())
//
//			By("Get a particular snapshot")
//			snapshots, err := f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).List(context.TODO(), metav1.ListOptions{LabelSelector: "workload-kind=Deployment"})
//			Expect(err).NotTo(HaveOccurred())
//			singleSnapshot, err := f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).Get(context.TODO(), snapshots.Items[len(snapshots.Items)-1].Name, metav1.GetOptions{})
//			Expect(err).NotTo(HaveOccurred())
//			Expect(singleSnapshot.Name).To(BeEquivalentTo(snapshots.Items[len(snapshots.Items)-1].Name))
//
//			By("Filter by workload kind")
//			snapshots, err = f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).List(context.TODO(), metav1.ListOptions{LabelSelector: "workload-kind=Deployment"})
//			Expect(err).NotTo(HaveOccurred())
//			Expect(snapshots).Should(HavePrefixInName("deployment"))
//
//			By("Filter by workload name")
//			snapshots, err = f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).List(context.TODO(), metav1.ListOptions{LabelSelector: "workload-name=" + deployment.Name})
//			Expect(err).NotTo(HaveOccurred())
//			workload.Kind = apis.KindDeployment
//			workload.Name = deployment.Name
//			Expect(snapshots).Should(HavePrefixInName(workload.GetRepositoryCRDName("", "")))
//
//			By("Filter by pod name")
//			snapshots, err = f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).List(context.TODO(), metav1.ListOptions{LabelSelector: "pod-name=" + ss.Name + "-0"})
//			Expect(err).NotTo(HaveOccurred())
//			workload.Kind = apis.KindStatefulSet
//			workload.Name = ss.Name
//			Expect(snapshots).Should(HavePrefixInName(workload.GetRepositoryCRDName(ss.Name+"-0", "")))
//
//			nodename := f.GetNodeName(daemon.ObjectMeta)
//			By("Filter by node name")
//			snapshots, err = f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).List(context.TODO(), metav1.ListOptions{LabelSelector: "node-name=" + nodename})
//			Expect(err).NotTo(HaveOccurred())
//			workload.Kind = apis.KindDaemonSet
//			workload.Name = daemon.Name
//			Expect(snapshots).Should(HavePrefixInName(workload.GetRepositoryCRDName("", nodename)))
//
//			workload.Kind = apis.KindDeployment
//			workload.Name = deployment.Name
//			reponame := workload.GetRepositoryCRDName("", "")
//
//			By("Filter by repository name")
//			snapshots, err = f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).List(context.TODO(), metav1.ListOptions{LabelSelector: "repository=" + reponame})
//			Expect(err).NotTo(HaveOccurred())
//			Expect(snapshots).Should(HavePrefixInName(reponame))
//
//			By("Filter by negated selector")
//			snapshots, err = f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).List(context.TODO(), metav1.ListOptions{LabelSelector: "repository!=" + reponame})
//			Expect(err).NotTo(HaveOccurred())
//			Expect(snapshots).ShouldNot(HavePrefixInName(reponame))
//
//			By("Filter by set based selector")
//			snapshots, err = f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).List(context.TODO(), metav1.ListOptions{LabelSelector: "repository in(" + reponame + ")"})
//			Expect(err).NotTo(HaveOccurred())
//			Expect(snapshots).Should(HavePrefixInName(reponame))
//
//			snapshots, err = f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).List(context.TODO(), metav1.ListOptions{LabelSelector: "workload-kind=Deployment"})
//			Expect(err).NotTo(HaveOccurred())
//			snapshotToDelete := snapshots.Items[len(snapshots.Items)-1].Name
//			By("Deleting snapshot " + snapshotToDelete)
//			policy := metav1.DeletePropagationForeground
//			err = f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).Delete(context.TODO(), snapshotToDelete, metav1.DeleteOptions{PropagationPolicy: &policy})
//			Expect(err).NotTo(HaveOccurred())
//
//			By("Checking deleted snapshot not exist")
//			_, err = f.StashClient.RepositoriesV1alpha1().Snapshots(f.Namespace()).Get(context.TODO(), snapshotToDelete, metav1.GetOptions{})
//			Expect(err).To(HaveOccurred())
//		}
//	)
//
//	Describe("Snapshots operations", func() {
//		Context(`"Local" backend`, func() {
//			AfterEach(func() {
//				err = framework.CleanupMinikubeHostPath()
//				Expect(err).NotTo(HaveOccurred())
//			})
//			BeforeEach(func() {
//				cred = f.SecretForLocalBackend()
//				restic = f.ResticForLocalBackend()
//			})
//			It(`should success to perform Snapshot's operations`, performOperationOnSnapshot)
//
//		})
//		Context(`"Minio" backend`, func() {
//			AfterEach(func() {
//				err = f.DeleteMinioServer()
//				Expect(err).NotTo(HaveOccurred())
//			})
//			BeforeEach(func() {
//				By("Checking if Restic installed in /bin directory")
//				cmd := exec.Command("restic")
//				err := cmd.Run()
//				if err != nil {
//					Skip("restic executable not found in /bin directory. Please install in /bin directory from: https://github.com/restic/restic/releases")
//				}
//
//				minikubeIP := net.IP{192, 168, 99, 100}
//
//				By("Creating Minio server without cacert")
//				_, err = f.CreateMinioServer(true, []net.IP{minikubeIP})
//				Expect(err).NotTo(HaveOccurred())
//
//				msvc, err := f.KubeClient.CoreV1().Services(f.Namespace()).Get(context.TODO(), "minio-service", metav1.GetOptions{})
//				Expect(err).NotTo(HaveOccurred())
//				minioServiceNodePort := strconv.Itoa(int(msvc.Spec.Ports[0].NodePort))
//
//				restic = f.ResticForMinioBackend("https://" + minikubeIP.String() + ":" + minioServiceNodePort)
//				cred = f.SecretForMinioBackend(true)
//			})
//			It(`should success to perform Snapshot's operations`, performOperationOnSnapshot)
//
//		})
//	})
//
//})
