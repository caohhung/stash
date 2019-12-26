/*
Copyright The Stash Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"fmt"

	"stash.appscode.dev/stash/apis"
	"stash.appscode.dev/stash/apis/stash/v1alpha1"
	"stash.appscode.dev/stash/apis/stash/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Invocation) PersistentVolumeClaim(name string) *core.PersistentVolumeClaim {
	return &core.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: f.namespace,
		},
		Spec: core.PersistentVolumeClaimSpec{
			AccessModes: []core.PersistentVolumeAccessMode{
				core.ReadWriteOnce,
			},
			StorageClassName: &f.StorageClass,
			Resources: core.ResourceRequirements{
				Requests: core.ResourceList{
					core.ResourceStorage: resource.MustParse("10Mi"),
				},
			},
		},
	}
}

func (f *Framework) CreatePersistentVolumeClaim(pvc *core.PersistentVolumeClaim) (*core.PersistentVolumeClaim, error) {
	return f.KubeClient.CoreV1().PersistentVolumeClaims(pvc.Namespace).Create(pvc)
}

func (f *Invocation) DeletePersistentVolumeClaim(meta metav1.ObjectMeta) error {
	err := f.KubeClient.CoreV1().PersistentVolumeClaims(meta.Namespace).Delete(meta.Name, deleteInForeground())
	if err != nil && !kerr.IsNotFound(err) {
		return err
	}
	return nil
}

func (f *Invocation) CreateNewPVC(name string) (*core.PersistentVolumeClaim, error) {
	// Generate PVC definition
	pvc := f.PersistentVolumeClaim(name)

	By(fmt.Sprintf("Creating PVC: %s/%s", pvc.Namespace, pvc.Name))
	createdPVC, err := f.CreatePersistentVolumeClaim(pvc)
	if err != nil {
		return nil, err
	}
	f.AppendToCleanupList(createdPVC)

	return createdPVC, nil
}

func (f *Invocation) SetupPVCBackup(pvc *core.PersistentVolumeClaim, repo *v1alpha1.Repository, transformFuncs ...func(bc *v1beta1.BackupConfiguration)) (*v1beta1.BackupConfiguration, error) {
	// Generate desired BackupConfiguration definition
	backupConfig := f.GetBackupConfiguration(repo.Name, func(bc *v1beta1.BackupConfiguration) {
		bc.Spec.Target = &v1beta1.BackupTarget{
			Ref: GetTargetRef(pvc.Name, apis.KindPersistentVolumeClaim),
		}
		bc.Spec.Task.Name = TaskPVCBackup
	})

	// transformFuncs provides a array of functions that made test specific change on the BackupConfiguration
	// apply these test specific changes
	for _, fn := range transformFuncs {
		fn(backupConfig)
	}

	By("Creating BackupConfiguration: " + backupConfig.Name)
	createdBC, err := f.StashClient.StashV1beta1().BackupConfigurations(backupConfig.Namespace).Create(backupConfig)
	f.AppendToCleanupList(createdBC)

	By("Verifying that backup triggering CronJob has been created")
	f.EventuallyCronJobCreated(backupConfig.ObjectMeta).Should(BeTrue())

	return createdBC, err
}

func (f *Invocation) SetupRestoreProcessForPVC(pvc *core.PersistentVolumeClaim, repo *v1alpha1.Repository, transformFuncs ...func(restore *v1beta1.RestoreSession)) (*v1beta1.RestoreSession, error) {
	// Generate desired RestoreSession definition
	By("Creating RestoreSession")
	restoreSession := f.GetRestoreSession(repo.Name, func(restore *v1beta1.RestoreSession) {
		restore.Spec.Target = &v1beta1.RestoreTarget{
			Ref: GetTargetRef(pvc.Name, apis.KindPersistentVolumeClaim),
		}
		restore.Spec.Rules = []v1beta1.Rule{
			{
				Snapshots: []string{"latest"},
			},
		}
		restore.Spec.Task.Name = TaskPVCRestore
	})

	// transformFuncs provides a array of functions that made test specific change on the RestoreSession
	// apply these test specific changes.
	for _, fn := range transformFuncs {
		fn(restoreSession)
	}

	err := f.CreateRestoreSession(restoreSession)
	f.AppendToCleanupList(restoreSession)

	By("Waiting for restore process to complete")
	f.EventuallyRestoreProcessCompleted(restoreSession.ObjectMeta).Should(BeTrue())

	return restoreSession, err
}
