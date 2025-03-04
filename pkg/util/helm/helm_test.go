package helm

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

var (
	NormalError   = errors.New("normal error")
	NotFoundError = errors.New("release name not found")
)

var mockedRelease = release.Release{Name: "test"}
var helmParam = &HelmParam{
	Repo{
		Name: "helm",
		URL:  "test1",
	},
	Chart{
		ReleaseName: "helm:v1.0.0",
		Timeout:     "1m",
		Wait:        true,
	},
}

type DefaultMockClient struct {
	helmclient.Client
}

func (c *DefaultMockClient) AddOrUpdateChartRepo(enrty repo.Entry) error {
	return nil
}

func (c *DefaultMockClient) UninstallReleaseByName(name string) error {
	return NotFoundError
}

func (c *DefaultMockClient) InstallOrUpgradeChart(ctx context.Context, spec *helmclient.ChartSpec) (*release.Release, error) {
	return &mockedRelease, nil
}

type DefaultMockClient2 struct {
	DefaultMockClient
}

func (c *DefaultMockClient2) UninstallReleaseByName(name string) error {
	return NormalError
}

type DefaultMockClient3 struct {
	DefaultMockClient
}

func (c *DefaultMockClient3) UninstallReleaseByName(name string) error {
	return nil
}

type DefaultMockClient4 struct {
	DefaultMockClient
}

func (c *DefaultMockClient4) AddOrUpdateChartRepo(enrty repo.Entry) error {
	return NormalError
}

func TestNewHelm(t *testing.T) {
	got, err := NewHelm(helmParam, WithClient(&DefaultMockClient{}))
	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	if got == nil {
		t.Errorf("got: %v must not be nil\n", got)
	}

	got, err = NewHelm(helmParam, WithClient(&DefaultMockClient4{}))
	if err != NormalError {
		t.Errorf("error: %v must be %v\n", err, NormalError)
	}
	if got != nil {
		t.Errorf("got: %v must be nil\n", got)
	}
}

func TestNewHelmWithOption(t *testing.T) {
	entry := &repo.Entry{
		Name:                  helmParam.Repo.Name,
		URL:                   helmParam.Repo.URL,
		Username:              "",
		Password:              "",
		CertFile:              "",
		KeyFile:               "",
		CAFile:                "",
		InsecureSkipTLSverify: false,
		PassCredentialsAll:    false,
	}
	atomic := true
	if !helmParam.Chart.Wait {
		atomic = false
	}
	tmout, err := time.ParseDuration(helmParam.Chart.Timeout)
	if err != nil {
		t.Log(err)
	}
	spec := &helmclient.ChartSpec{
		ReleaseName:      helmParam.Chart.ReleaseName,
		ChartName:        helmParam.Chart.ChartName,
		Namespace:        helmParam.Chart.Namespace,
		ValuesYaml:       helmParam.Chart.ValuesYaml,
		Version:          helmParam.Chart.Version,
		CreateNamespace:  false,
		DisableHooks:     false,
		Replace:          true,
		Wait:             helmParam.Chart.Wait,
		DependencyUpdate: false,
		Timeout:          tmout,
		GenerateName:     false,
		NameTemplate:     "",
		Atomic:           atomic,
		SkipCRDs:         false,
		UpgradeCRDs:      helmParam.Chart.UpgradeCRDs,
		SubNotes:         false,
		Force:            false,
		ResetValues:      false,
		ReuseValues:      false,
		Recreate:         false,
		MaxHistory:       0,
		CleanupOnFail:    false,
		DryRun:           false,
	}

	mockClient := &DefaultMockClient{}

	got, err := NewHelm(helmParam, WithClient(mockClient))
	if err != nil {
		t.Errorf("err: %v\n", err)
	}

	want := &Helm{
		Entry:     entry,
		ChartSpec: spec,
		Client:    mockClient,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewHelm() = \n%+v\n, want \n%+v\n", got, want)
	}
}

func TestInstallOrUpgradeChart(t *testing.T) {
	atomic := true
	if !helmParam.Chart.Wait {
		atomic = false
	}
	tmout, err := time.ParseDuration(helmParam.Chart.Timeout)
	if err != nil {
		t.Log(err)
	}
	chartSpec := &helmclient.ChartSpec{
		ReleaseName:      helmParam.Chart.ReleaseName,
		ChartName:        helmParam.Chart.ChartName,
		Namespace:        helmParam.Chart.Namespace,
		ValuesYaml:       helmParam.Chart.ValuesYaml,
		Version:          helmParam.Chart.Version,
		CreateNamespace:  false,
		DisableHooks:     false,
		Replace:          true,
		Wait:             helmParam.Chart.Wait,
		DependencyUpdate: false,
		Timeout:          tmout,
		GenerateName:     false,
		NameTemplate:     "",
		Atomic:           atomic,
		SkipCRDs:         false,
		UpgradeCRDs:      helmParam.Chart.UpgradeCRDs,
		SubNotes:         false,
		Force:            false,
		ResetValues:      false,
		ReuseValues:      false,
		Recreate:         false,
		MaxHistory:       0,
		CleanupOnFail:    false,
		DryRun:           false,
	}

	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	// mockClient := mockhelmclient.NewMockClient(ctrl)
	// if mockClient == nil {
	// 	t.Fail()
	// }
	h, err := NewHelm(helmParam, WithChartSpec(chartSpec), WithClient(&DefaultMockClient{}))

	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	// mockClient.EXPECT().InstallOrUpgradeChart(context.TODO(), chartSpec).Return(&mockedRelease, nil)

	err = h.InstallOrUpgradeChart()
	if err != nil {
		t.Error(err)
	}
}

func TestAddOrUpdateChartRepo(t *testing.T) {
	entry := &repo.Entry{
		Name:                  helmParam.Repo.Name,
		URL:                   helmParam.Repo.URL,
		Username:              "",
		Password:              "",
		CertFile:              "",
		KeyFile:               "",
		CAFile:                "",
		InsecureSkipTLSverify: false,
		PassCredentialsAll:    false,
	}
	atomic := true
	if !helmParam.Chart.Wait {
		atomic = false
	}
	tmout, err := time.ParseDuration(helmParam.Chart.Timeout)
	if err != nil {
		t.Log(err)
	}
	chartSpec := &helmclient.ChartSpec{
		ReleaseName:      helmParam.Chart.ReleaseName,
		ChartName:        helmParam.Chart.ChartName,
		Namespace:        helmParam.Chart.Namespace,
		ValuesYaml:       helmParam.Chart.ValuesYaml,
		Version:          helmParam.Chart.Version,
		CreateNamespace:  false,
		DisableHooks:     false,
		Replace:          true,
		Wait:             helmParam.Chart.Wait,
		DependencyUpdate: false,
		Timeout:          tmout,
		GenerateName:     false,
		NameTemplate:     "",
		Atomic:           atomic,
		SkipCRDs:         false,
		UpgradeCRDs:      helmParam.Chart.UpgradeCRDs,
		SubNotes:         false,
		Force:            false,
		ResetValues:      false,
		ReuseValues:      false,
		Recreate:         false,
		MaxHistory:       0,
		CleanupOnFail:    false,
		DryRun:           false,
	}

	h, err := NewHelm(helmParam, WithEntry(entry), WithChartSpec(chartSpec), WithClient(&DefaultMockClient{}))
	if err != nil {
		t.Errorf("error: %v\n", err)
	}

	err = h.AddOrUpdateChartRepo(*entry)
	if err != nil {
		t.Error(err)
	}
}

func TestHelm_UninstallHelmChartRelease(t *testing.T) {
	atomic := true
	if !helmParam.Chart.Wait {
		atomic = false
	}
	tmout, err := time.ParseDuration(helmParam.Chart.Timeout)
	if err != nil {
		t.Log(err)
	}
	chartSpec := &helmclient.ChartSpec{
		ReleaseName:      helmParam.Chart.ReleaseName,
		ChartName:        helmParam.Chart.ChartName,
		Namespace:        helmParam.Chart.Namespace,
		ValuesYaml:       helmParam.Chart.ValuesYaml,
		Version:          helmParam.Chart.Version,
		CreateNamespace:  false,
		DisableHooks:     false,
		Replace:          true,
		Wait:             helmParam.Chart.Wait,
		DependencyUpdate: false,
		Timeout:          tmout,
		GenerateName:     false,
		NameTemplate:     "",
		Atomic:           atomic,
		SkipCRDs:         false,
		UpgradeCRDs:      helmParam.Chart.UpgradeCRDs,
		SubNotes:         false,
		Force:            false,
		ResetValues:      false,
		ReuseValues:      false,
		Recreate:         false,
		MaxHistory:       0,
		CleanupOnFail:    false,
		DryRun:           false,
	}
	// base
	h, err := NewHelm(helmParam, WithChartSpec(chartSpec), WithClient(&DefaultMockClient3{}))
	if err != nil {
		t.Errorf("error: %v\n", err)
	}

	err = h.UninstallHelmChartRelease()
	if err != nil {
		t.Error(err)
	}

	// mock error not found
	h, err = NewHelm(helmParam, WithChartSpec(chartSpec), WithClient(&DefaultMockClient{}))
	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	err = h.UninstallHelmChartRelease()
	if err != nil {
		t.Error(err)
	}

	// mock error
	h, err = NewHelm(helmParam, WithChartSpec(chartSpec), WithClient(&DefaultMockClient2{}))
	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	err = h.UninstallHelmChartRelease()
	if err == nil {
		t.Error("error not found")
	}
	if err != NormalError {
		t.Errorf("got: %+v\n, want %+v\n", err, NormalError)
	}
}
