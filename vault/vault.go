package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/vcs"
	"github.com/fatih/color"
	v1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func UnSetVaultUrl(cluster kinds.ProviderCluster) error {

	if cluster.Vault.Url != "" {
		token := os.Getenv("VAULT_TOKEN")
		// Set the Vault sys/auth paht -----------------------------------------------------
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/sys/auth/%s", cluster.Vault.Url, cluster.Vault.Path), nil)
		if err != nil {
			return err
		}
		req.Header.Set("X-Vault-Token", token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Set the Vault auth paht -----------------------------------------------------
		req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/v1/auth/%s/config", cluster.Vault.Url, cluster.Vault.Path), nil)
		if err != nil {
			return err
		}
		req.Header.Set("X-Vault-Token", token)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		color.Green(fmt.Sprintf("Cluster %s unset on Vault Cluster", cluster.ShortName))
	}
	return nil
}

func UnSetVaultGit(opt configuration.Options, cluster kinds.ProviderCluster) error {

	if cluster.Vault.Repo.Url != "" {
		var extension = filepath.Ext(cluster.Vault.Repo.Url)
		var remoteVCSRepoName = cluster.Vault.Repo.Url[0 : len(cluster.Vault.Repo.Url)-len(extension)]
		splitStrings := strings.Split(remoteVCSRepoName, "/")
		remoteVCSRepoName = splitStrings[len(splitStrings)-1]

		vaultTempDir, err := ioutil.TempDir("/tmp", "vault_")
		if err != nil {
			return err
		}

		color.Green(fmt.Sprintf("Fetching deployment %s into %s\n", remoteVCSRepoName, vaultTempDir))
		gvcs := new(vcs.GitVCS)
		_, err = vcs.Fetch(gvcs, vaultTempDir, cluster.Vault.Repo.Url, opt.SSHKeyPath, cluster.Vault.Repo.Branch)
		if err != nil {
			return err
		}

		vaultSetConfigPath := path.Join(vaultTempDir, cluster.Vault.Repo.Path, "set/auth", cluster.Vault.Path, "config.json")
		err = os.Remove(vaultSetConfigPath)
		if err != nil {
			return err
		}

		vaultUnSetConfigPath := path.Join(vaultTempDir, cluster.Vault.Repo.Path, "unset/auth", cluster.Vault.Path, "config.json")
		dirName := filepath.Dir(vaultUnSetConfigPath)
		if _, serr := os.Stat(dirName); serr != nil {
			merr := os.MkdirAll(dirName, os.ModePerm)
			if merr != nil {
				return err
			}
		}
		emptyFile, err := os.Create(vaultUnSetConfigPath)
		if err != nil {
			return err
		}
		emptyFile.Close()

		vaultSetSysConfigPath := path.Join(vaultTempDir, cluster.Vault.Repo.Path, "set/sys/auth", fmt.Sprintf("%s.%s", cluster.Vault.Path, "json"))
		err = os.Remove(vaultSetSysConfigPath)
		if err != nil {
			return err
		}

		vaultUnSetSysConfigPath := path.Join(vaultTempDir, cluster.Vault.Repo.Path, "unset/sys/auth", fmt.Sprintf("%s.%s", cluster.Vault.Path, "json"))
		dirName = filepath.Dir(vaultUnSetSysConfigPath)
		if _, serr := os.Stat(dirName); serr != nil {
			merr := os.MkdirAll(dirName, os.ModePerm)
			if merr != nil {
				return err
			}
		}
		emptyFile, err = os.Create(vaultUnSetSysConfigPath)
		if err != nil {
			return err
		}
		emptyFile.Close()

		err = vcs.Add(gvcs, vaultTempDir, cluster.Vault.Repo.Url, opt.SSHKeyPath, []string{
			path.Join(cluster.Vault.Repo.Path, "set/sys/auth", fmt.Sprintf("%s.%s", cluster.Vault.Path, "json")),
			path.Join(cluster.Vault.Repo.Path, "set/auth", cluster.Vault.Path, "config.json"),
			path.Join(cluster.Vault.Repo.Path, "unset/sys/auth", fmt.Sprintf("%s.%s", cluster.Vault.Path, "json")),
			path.Join(cluster.Vault.Repo.Path, "unset/auth", cluster.Vault.Path, "config.json"),
		})
		if err != nil {
			return err
		}

		err = vcs.Commit(gvcs, vaultTempDir, cluster.Vault.Repo.Url, opt.SSHKeyPath, fmt.Sprintf("Remove kubernetes secret setting for %s on %s project", cluster.ShortName, cluster.Project))
		if err != nil {
			return err
		}

		err = vcs.Push(gvcs, vaultTempDir, cluster.Vault.Repo.Url, opt.SSHKeyPath)
		if err != nil {
			return err
		}

		defer os.RemoveAll(vaultTempDir)
		color.Green(fmt.Sprintf("Cluster %s unset on git repo %s on branch %s on the folder %s/unset.", cluster.ShortName, cluster.Vault.Repo.Url, cluster.Vault.Repo.Branch, cluster.Vault.Repo.Path))
	}
	return nil
}

func SetVaultUrl(cluster kinds.ProviderCluster, vaultSysAuthJsonContent []byte, vaultAuthJsonContent []byte) error {
	if cluster.Vault.Url != "" {

		token := os.Getenv("VAULT_TOKEN")
		// Set the Vault sys/auth paht -----------------------------------------------------
		data := bytes.NewReader(vaultSysAuthJsonContent)
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/sys/auth/%s", cluster.Vault.Url, cluster.Vault.Path), data)
		if err != nil {
			return err
		}
		req.Header.Set("X-Vault-Token", token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Set the Vault auth paht -----------------------------------------------------
		data = bytes.NewReader(vaultAuthJsonContent)
		req, err = http.NewRequest("POST", fmt.Sprintf("%s/v1/auth/%s/config", cluster.Vault.Url, cluster.Vault.Path), data)
		if err != nil {
			return err
		}
		req.Header.Set("X-Vault-Token", token)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		color.Green(fmt.Sprintf("Cluster %s set on Vault Cluster with path: %s", cluster.ShortName, cluster.Vault.Path))
	}
	return nil
}

func SetVaultConfiguration(opt configuration.Options, cluster kinds.ProviderCluster) error {

	if cluster.Vault.Path != "" {

		// Creates the clientset
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

		k8sconfig, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
		clientset, err := kubernetes.NewForConfig(k8sconfig)
		if err != nil {
			return err
		}

		vaultK8SServiceAccount, err := SetK8SVaultServiceAccount(clientset, "vault-tokenreview", "kube-system")
		if err != nil {
			return err
		}
		err = SetK8SVaultClusterRoleBinding(clientset, "vault-tokenreview-binding", "vault-tokenreview", "kube-system")
		if err != nil {
			return err
		}

		var secret = &v1.Secret{}
		var vaultAuthJsonContent = []byte{}
		var vaultSysAuthJsonContent = []byte{}
		for _, ref := range vaultK8SServiceAccount.Secrets {
			secret, err = clientset.Core().Secrets("kube-system").Get(ref.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			content := map[string]string{
				"kubernetes_host":    fmt.Sprintf("https://%s", cluster.Endpoint),
				"kubernetes_ca_cert": fmt.Sprintf("%s", secret.Data["ca.crt"]),
				"token_reviewer_jwt": fmt.Sprintf("%s", secret.Data["token"]),
			}
			vaultAuthJsonContent, _ = json.Marshal(content)
			if err != nil {
				return err
			}

			content = map[string]string{
				"path":        fmt.Sprintf("%s", cluster.Vault.Path),
				"description": fmt.Sprintf("%s", cluster.Vault.Description),
				"type":        "kubernetes",
			}
			vaultSysAuthJsonContent, _ = json.Marshal(content)
			if err != nil {
				return err
			}
		}

		if err := SetVaultUrl(cluster, vaultSysAuthJsonContent, vaultAuthJsonContent); err != nil {
			return err
		}

		if err := SetVaultGit(opt, cluster, vaultSysAuthJsonContent, vaultAuthJsonContent); err != nil {
			return err
		}
	}
	return nil
}

func SetVaultGit(opt configuration.Options, cluster kinds.ProviderCluster, vaultSysAuthJsonContent []byte, vaultAuthJsonContent []byte) error {

	if cluster.Vault.Repo.Url != "" {
		var extension = filepath.Ext(cluster.Vault.Repo.Url)
		var remoteVCSRepoName = cluster.Vault.Repo.Url[0 : len(cluster.Vault.Repo.Url)-len(extension)]
		splitStrings := strings.Split(remoteVCSRepoName, "/")
		remoteVCSRepoName = splitStrings[len(splitStrings)-1]

		vaultTempDir, err := ioutil.TempDir("/tmp", "vault_")
		if err != nil {
			return err
		}

		color.Green(fmt.Sprintf("Fetching deployment %s into %s\n", remoteVCSRepoName, vaultTempDir))
		gvcs := new(vcs.GitVCS)
		_, err = vcs.Fetch(gvcs, vaultTempDir, cluster.Vault.Repo.Url, opt.SSHKeyPath, cluster.Vault.Repo.Branch)
		if err != nil {
			return err
		}

		vaultConfigPath := path.Join(vaultTempDir, cluster.Vault.Repo.Path, "set/auth", cluster.Vault.Path, "config.json")
		dirName := filepath.Dir(vaultConfigPath)
		if _, serr := os.Stat(dirName); serr != nil {
			merr := os.MkdirAll(dirName, os.ModePerm)
			if merr != nil {
				return err
			}
		}
		if err := ioutil.WriteFile(vaultConfigPath, vaultAuthJsonContent, 0644); err != nil {
			return err
		}

		vaultSysConfigPath := path.Join(vaultTempDir, cluster.Vault.Repo.Path, "set/sys/auth", fmt.Sprintf("%s.%s", cluster.Vault.Path, "json"))
		dirName = filepath.Dir(vaultSysConfigPath)
		if _, serr := os.Stat(dirName); serr != nil {
			merr := os.MkdirAll(dirName, os.ModePerm)
			if merr != nil {
				return err
			}
		}
		err = ioutil.WriteFile(vaultSysConfigPath, vaultSysAuthJsonContent, 0644)
		if err != nil {
			return err
		}

		err = vcs.Add(gvcs, vaultTempDir, cluster.Vault.Repo.Url, opt.SSHKeyPath, []string{
			path.Join(cluster.Vault.Repo.Path, "set/sys/auth", fmt.Sprintf("%s.%s", cluster.Vault.Path, "json")),
			path.Join(cluster.Vault.Repo.Path, "set/auth", cluster.Vault.Path, "config.json"),
		})
		if err != nil {
			return err
		}

		err = vcs.Commit(gvcs, vaultTempDir, cluster.Vault.Repo.Url, opt.SSHKeyPath, fmt.Sprintf("Set kubernetes secret setting for %s on %s project", cluster.ShortName, cluster.Project))
		if err != nil {
			return err
		}

		err = vcs.Push(gvcs, vaultTempDir, cluster.Vault.Repo.Url, opt.SSHKeyPath)
		if err != nil {
			return err
		}

		defer os.RemoveAll(vaultTempDir)
		color.Green(fmt.Sprintf("Cluster %s set on git repo %s on branch %s on the folder %s/set.", cluster.ShortName, cluster.Vault.Repo.Url, cluster.Vault.Repo.Branch, cluster.Vault.Repo.Path))
	}

	return nil
}

// Set K8S Service Account -----------------------------------------------------
func SetK8SVaultServiceAccount(clientset *kubernetes.Clientset, name string, namespace string) (*v1.ServiceAccount, error) {
	_, err := clientset.Core().ServiceAccounts(namespace).Get(name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		color.Green(fmt.Sprintf("Set %s ServiceAccount in %s Namespace", name, namespace))
		K8SServiceAccount := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}
		_, err = clientset.CoreV1().ServiceAccounts(namespace).Create(K8SServiceAccount)
		if err != nil {
			return nil, err
		}
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		return nil, statusError
	} else if err != nil {
		return nil, err
	} else {
		color.Green(fmt.Sprintf("Found %s ServiceAccount\n", "vault-tokenreview"))
	}
	return clientset.Core().ServiceAccounts(namespace).Get(name, metav1.GetOptions{})
}

// Set K8S ClusterRoleBinding -----------------------------------------------------
func SetK8SVaultClusterRoleBinding(clientset *kubernetes.Clientset, name string, serviceaccount string, namespace string) error {
	_, err := clientset.RbacV1beta1().ClusterRoleBindings().Get(name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		color.Green(fmt.Sprintf("Set %s ClusterRoleBinding", name))
		K8SClusterRoleBinding := &rbacv1beta1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Subjects: []rbacv1beta1.Subject{
				{Kind: "ServiceAccount", Name: serviceaccount, Namespace: namespace},
			},
			RoleRef: rbacv1beta1.RoleRef{Kind: "ClusterRole", Name: "system:auth-delegator"},
		}
		_, err = clientset.RbacV1beta1().ClusterRoleBindings().Create(K8SClusterRoleBinding)
		if err != nil {
			return err
		}

	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		return statusError
	} else if err != nil {
		return err
	} else {
		color.Green(fmt.Sprintf("Found %s ClusterRoleBinding\n", name))
	}
	return nil
}
