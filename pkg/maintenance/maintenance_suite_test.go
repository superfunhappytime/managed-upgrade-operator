package maintenance

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMaintenance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Maintenance Suite")
}
