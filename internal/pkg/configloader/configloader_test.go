package configloader_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
)

var _ = Describe("Dependency", func() {
	Context("singe dep", func() {
		tools := []configloader.Tool{
			{InstanceID: "argocd", Name: "argocd"},
			{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{"argocd.argocd"}},
		}
		config := configloader.Config{
			Tools: tools,
		}
		errors := config.ValidateDependency()
		Expect(len(errors)).To(Equal(0))
	})

	Context("dep not exist", func() {
		tools := []configloader.Tool{
			{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{"argocd.argocd"}},
		}
		config := configloader.Config{
			Tools: tools,
		}
		errors := config.ValidateDependency()
		Expect(len(errors)).To(Equal(1))
	})

	Context("multi-dep", func() {
		tools := []configloader.Tool{
			{InstanceID: "argocd", Name: "argocd"},
			{InstanceID: "repo", Name: "github"},
			{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{"argocd.argocd", "github.repo"}},
		}
		config := configloader.Config{
			Tools: tools,
		}
		errors := config.ValidateDependency()
		Expect(len(errors)).To(Equal(0))
	})

	Context("empty dep", func() {
		tools := []configloader.Tool{
			{InstanceID: "argocd", Name: "argocd"},
			{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{}},
		}
		config := configloader.Config{
			Tools: tools,
		}
		errors := config.ValidateDependency()
		Expect(len(errors)).To(Equal(0))
	})
})

var _ = Describe("Tool Validation", func() {
	It("should return empty error array if tools all valid", func() {
		tools := []configloader.Tool{
			{Name: "test_tool", InstanceID: "0", DependsOn: []string{}},
		}
		config := configloader.Config{
			Tools: tools,
		}
		errors := config.Validate()
		Expect(errors).Should(BeEmpty())
	})
	It("should return error if tool not valid", func() {
		tools := []configloader.Tool{
			{Name: "", InstanceID: "", DependsOn: []string{}},
		}
		config := configloader.Config{
			Tools: tools,
		}
		errors := config.Validate()
		Expect(errors).ShouldNot(BeEmpty())
	})

})
