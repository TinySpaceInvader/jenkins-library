package npm

import (
	"github.com/SAP/jenkins-library/pkg/mock"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

type npmMockUtilsBundle struct {
	*mock.FilesMock
	execRunner *mock.ExecMockRunner
}

func (u *npmMockUtilsBundle) GetExecRunner() ExecRunner {
	return u.execRunner
}

func newNpmMockUtilsBundle() npmMockUtilsBundle {
	utils := npmMockUtilsBundle{FilesMock: &mock.FilesMock{}, execRunner: &mock.ExecMockRunner{}}
	return utils
}

func TestNpm(t *testing.T) {
	t.Run("find package.json files with one package.json", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"name\": \"Test\" }"))

		options := ExecutorOptions{}

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}

		packageJSONFiles := exec.FindPackageJSONFiles()

		assert.Equal(t, []string{"package.json"}, packageJSONFiles)
	})

	t.Run("find package.json files with two package.json and filtered package.json", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{}"))
		utils.AddFile(filepath.Join("src", "package.json"), []byte("{}"))
		utils.AddFile(filepath.Join("node_modules", "package.json"), []byte("{}")) // is filtered out
		options := ExecutorOptions{}

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}

		packageJSONFiles := exec.FindPackageJSONFiles()

		assert.Equal(t, []string{"package.json", filepath.Join("src", "package.json")}, packageJSONFiles)
	})

	t.Run("find package.json files with script", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))
		utils.AddFile(filepath.Join("src", "package.json"), []byte("{ \"name\": \"test\" }"))
		utils.AddFile(filepath.Join("test", "package.json"), []byte("{ \"scripts\": { \"test\": \"exit 0\" } }"))

		options := ExecutorOptions{}

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}

		packageJSONFilesWithScript, err := exec.FindPackageJSONFilesWithScript([]string{"package.json", filepath.Join("src", "package.json"), filepath.Join("test", "package.json")}, "ci-lint")

		if assert.NoError(t, err) {
			assert.Equal(t, []string{"package.json"}, packageJSONFilesWithScript)
		}
	})

	t.Run("Install deps for package.json with package-lock.json", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))
		utils.AddFile("package-lock.json", []byte("{}"))

		options := ExecutorOptions{}
		options.DefaultNpmRegistry = "foo.bar"

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}
		err := exec.install("package.json")

		if assert.NoError(t, err) {
			if assert.Equal(t, 3, len(utils.execRunner.Calls)) {
				assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"ci"}}, utils.execRunner.Calls[2])
			}
		}
	})

	t.Run("Install deps for package.json without package-lock.json", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))

		options := ExecutorOptions{}
		options.DefaultNpmRegistry = "foo.bar"

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}
		err := exec.install("package.json")

		if assert.NoError(t, err) {
			if assert.Equal(t, 3, len(utils.execRunner.Calls)) {
				assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"install"}}, utils.execRunner.Calls[2])
			}
		}
	})

	t.Run("Install deps for package.json with yarn.lock", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))
		utils.AddFile("yarn.lock", []byte("{}"))

		options := ExecutorOptions{}
		options.DefaultNpmRegistry = "foo.bar"

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}
		err := exec.install("package.json")

		if assert.NoError(t, err) {
			if assert.Equal(t, 3, len(utils.execRunner.Calls)) {
				assert.Equal(t, mock.ExecCall{Exec: "yarn", Params: []string{"install", "--frozen-lockfile"}}, utils.execRunner.Calls[2])
			}
		}
	})

	t.Run("Install all deps", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))
		utils.AddFile("package-lock.json", []byte("{}"))
		utils.AddFile(filepath.Join("src", "package.json"), []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))
		utils.AddFile(filepath.Join("src", "package-lock.json"), []byte("{}"))

		options := ExecutorOptions{}
		options.DefaultNpmRegistry = "foo.bar"

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}
		err := exec.InstallAllDependencies([]string{"package.json", filepath.Join("src", "package.json")})

		if assert.NoError(t, err) {
			if assert.Equal(t, 6, len(utils.execRunner.Calls)) {
				assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"ci"}}, utils.execRunner.Calls[2])
				assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"ci"}}, utils.execRunner.Calls[5])
			}
		}
	})

	t.Run("check if yarn.lock and package-lock exist", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))
		utils.AddFile("yarn.lock", []byte("{}"))
		utils.AddFile("package-lock.json", []byte("{}"))

		options := ExecutorOptions{}

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}
		packageLock, yarnLock, err := exec.checkIfLockFilesExist()

		if assert.NoError(t, err) {
			assert.True(t, packageLock)
			assert.True(t, yarnLock)
		}
	})

	t.Run("check that yarn.lock and package-lock do not exist", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))

		options := ExecutorOptions{}
		options.SapNpmRegistry = "foo.sap"

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}
		packageLock, yarnLock, err := exec.checkIfLockFilesExist()

		if assert.NoError(t, err) {
			assert.False(t, packageLock)
			assert.False(t, yarnLock)
		}
	})

	t.Run("check Execute script", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))

		options := ExecutorOptions{}

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}
		err := exec.executeScript("package.json", "ci-lint", []string{"--silent"}, []string{"--tag", "tag1"})

		if assert.NoError(t, err) {
			if assert.Equal(t, 3, len(utils.execRunner.Calls)) {
				assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"run", "ci-lint", "--silent", "--", "--tag", "tag1"}}, utils.execRunner.Calls[2])
			}
		}
	})

	t.Run("check Execute all scripts", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))
		utils.AddFile(filepath.Join("src", "package.json"), []byte("{\"scripts\": { \"ci-build\": \"exit 0\" } }"))

		options := ExecutorOptions{}
		runScripts := []string{"ci-lint", "ci-build"}

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}
		err := exec.RunScriptsInAllPackages(runScripts, nil, nil, false)

		if assert.NoError(t, err) {
			if assert.Equal(t, 6, len(utils.execRunner.Calls)) {
				assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"run", "ci-lint"}}, utils.execRunner.Calls[2])
				assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"run", "ci-build"}}, utils.execRunner.Calls[5])
			}
		}
	})

	t.Run("check set npm registry", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))
		utils.AddFile(filepath.Join("src", "package.json"), []byte("{\"scripts\": { \"ci-build\": \"exit 0\" } }"))
		utils.execRunner = &mock.ExecMockRunner{StdoutReturn: map[string]string{"npm config get registry": "undefined"}}
		options := ExecutorOptions{}
		options.DefaultNpmRegistry = "https://example.org/npm"

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}
		err := exec.SetNpmRegistries()

		if assert.NoError(t, err) {
			if assert.Equal(t, 3, len(utils.execRunner.Calls)) {
				assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"config", "get", "registry"}}, utils.execRunner.Calls[0])
				assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"config", "set", "registry", exec.Options.DefaultNpmRegistry}}, utils.execRunner.Calls[1])
			}
		}
	})

	t.Run("check set npm registry", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"ci-lint\": \"exit 0\" } }"))
		utils.AddFile(filepath.Join("src", "package.json"), []byte("{\"scripts\": { \"ci-build\": \"exit 0\" } }"))
		utils.execRunner = &mock.ExecMockRunner{StdoutReturn: map[string]string{"npm config get @sap:registry": "undefined"}}
		options := ExecutorOptions{}
		options.SapNpmRegistry = "https://example.sap/npm"

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}
		err := exec.SetNpmRegistries()

		if assert.NoError(t, err) {
			if assert.Equal(t, 3, len(utils.execRunner.Calls)) {
				assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"config", "get", "@sap:registry"}}, utils.execRunner.Calls[1])
				assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"config", "set", "@sap:registry", exec.Options.SapNpmRegistry}}, utils.execRunner.Calls[2])
			}
		}
	})

	t.Run("Call run-scripts with virtual frame buffer", func(t *testing.T) {
		utils := newNpmMockUtilsBundle()
		utils.AddFile("package.json", []byte("{\"scripts\": { \"foo\": \"\" } }"))

		options := ExecutorOptions{}

		exec := &Execute{
			Utils:   &utils,
			Options: options,
		}
		err := exec.RunScriptsInAllPackages([]string{"foo"}, nil, nil, true)

		assert.Contains(t, utils.execRunner.Env, "DISPLAY=:99")
		assert.NoError(t, err)
		if assert.Len(t, utils.execRunner.Calls, 4) {
			xvfbCall := utils.execRunner.Calls[0]
			assert.Equal(t, "Xvfb", xvfbCall.Exec)
			assert.Equal(t, []string{"-ac", ":99", "-screen", "0", "1280x1024x16"}, xvfbCall.Params)
			assert.True(t, xvfbCall.Async)
			assert.True(t, xvfbCall.Execution.Killed)

			assert.Equal(t, mock.ExecCall{Exec: "npm", Params: []string{"run", "foo"}}, utils.execRunner.Calls[3])
		}
	})
}
