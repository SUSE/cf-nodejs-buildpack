package integration_test

import (
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack/cutlass"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("override yml", func() {
	var app *cutlass.App
	var buildpackName string
	AfterEach(func() {
		if buildpackName != "" {
			cutlass.DeleteBuildpack(buildpackName)
		}

		if app != nil {
			app.Destroy()
		}
		app = nil
	})

	BeforeEach(func() {
		if !ApiHasMultiBuildpack() {
			Skip("Multi buildpack support is required")
		}

		buildpackName = "override_yml_" + cutlass.RandStringRunes(5)
		Expect(cutlass.CreateOrUpdateBuildpack(buildpackName, filepath.Join(bpDir, "fixtures", "overrideyml_bp"), "")).To(Succeed())

		app = cutlass.New(filepath.Join(bpDir, "fixtures", "simple_app"))
		app.Buildpacks = []string{buildpackName + "_buildpack", "nodejs_buildpack"}
	})

	It("Forces node from override buildpack", func() {
		Expect(app.V3Push()).ToNot(Succeed())
		Eventually(func() error { return app.ConfirmBuildpack(buildpackVersion) }, "30s").Should(Succeed())
		Eventually(app.Stdout.String).Should(ContainSubstring("-----> OverrideYML Buildpack"))

		Eventually(app.Stdout.String).Should(ContainSubstring("-----> Installing node"))
		Eventually(app.Stdout.String).Should(MatchRegexp("Copy .*/node.tgz"))
		Eventually(app.Stdout.String).Should(ContainSubstring("Unable to install node: dependency sha256 mismatch: expected sha256"))
	})
})
