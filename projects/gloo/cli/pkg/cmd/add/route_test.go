package add_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
)

var _ = Describe("Upstream", func() {

	BeforeEach(func() {
		helpers.MemoryResourceClient = &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
	})

	BeforeEach(func() {
		err := testutils.Glooctl("create upstream static default-petstore-8080 --static-hosts jsonplaceholder.typicode.com:80")
		Expect(err).NotTo(HaveOccurred())

	})

	It("should create static upstream", func() {
		err := testutils.Glooctl("add route --path-exact /sample-route-1 --dest-name default-petstore-8080 --prefix-rewrite /api/pets")
		Expect(err).NotTo(HaveOccurred())

		up, err := helpers.MustVirtualServiceClient().Read("gloo-system", "default", clients.ReadOpts{})
		Expect(up.Metadata.Name).To(Equal("default"))
	})
})
