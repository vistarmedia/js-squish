package jssquish

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resolver", func() {

	var (
		repo     *MemRepository
		resolver *Resolver
	)

	BeforeEach(func() {
		repo = NewMemRepository(map[string]string{
			"project-one/a.js":   "",
			"project-one/b.js":   "",
			"project-one/c.json": "",

			"project-two/a.js":   "",
			"project-two/b.js":   "",
			"project-two/c.json": "",

			"project-three/a/index.js": "",
			"project-three/b.js":       "",

			"project-four/four-main.js": "",
			"project-four/index.js":     "",
			"project-four/package.json": `{"main": "./four-main.js"}`,
		})

		resolver = NewResolver(repo)
	})

	It("should not resolve a file that doesn't exist", func() {
		_, err := resolver.Resolve("missing", ".")
		Expect(err).To(HaveOccurred())
		Expect(repo.opened).To(BeEmpty())
	})

	It("should resolve a simple file", func() {
		fq, err := resolver.Resolve("project-one/a.js", ".")
		Expect(err).ToNot(HaveOccurred())
		Expect(fq).To(Equal("project-one/a.js"))
	})

	It("should resolve a relative file", func() {
		fq, err := resolver.Resolve("./a.js", "project-one")
		Expect(err).ToNot(HaveOccurred())
		Expect(fq).To(Equal("project-one/a.js"))

		fq, err = resolver.Resolve("./a.js", "project-two")
		Expect(err).ToNot(HaveOccurred())
		Expect(fq).To(Equal("project-two/a.js"))
	})

	It("should walk up a directory", func() {
		fq, err := resolver.Resolve("../project-one/b.js", "project-two")
		Expect(err).ToNot(HaveOccurred())
		Expect(fq).To(Equal("project-one/b.js"))
	})

	It("should assume a .js extension", func() {
		fq, err := resolver.Resolve("project-one/b", "project-two")
		Expect(err).ToNot(HaveOccurred())
		Expect(fq).To(Equal("project-one/b.js"))

		Expect(repo.checked).To(HaveLen(2))
		Expect(repo.checked[0]).To(Equal("project-one/b"))
		Expect(repo.checked[1]).To(Equal("project-one/b.js"))
	})

	It("should assume a .json extension", func() {
		fq, err := resolver.Resolve("project-one/c", ".")
		Expect(err).ToNot(HaveOccurred())
		Expect(fq).To(Equal("project-one/c.json"))

		Expect(repo.checked).To(HaveLen(3))
		Expect(repo.checked[0]).To(Equal("project-one/c"))
		Expect(repo.checked[1]).To(Equal("project-one/c.js"))
		Expect(repo.checked[2]).To(Equal("project-one/c.json"))
	})

	It("should resolve a directory with an index.js", func() {
		fq, err := resolver.Resolve("project-three/a", ".")
		Expect(err).ToNot(HaveOccurred())
		Expect(fq).To(Equal("project-three/a/index.js"))

		Expect(repo.checked).To(HaveLen(5))
		Expect(repo.checked[0]).To(Equal("project-three/a"))
		Expect(repo.checked[1]).To(Equal("project-three/a.js"))
		Expect(repo.checked[2]).To(Equal("project-three/a.json"))
		Expect(repo.checked[3]).To(Equal("project-three/a/package.json"))
		Expect(repo.checked[4]).To(Equal("project-three/a/index.js"))
	})

	It("should resolve a directory package.json specifying a 'main'", func() {
		fq, err := resolver.Resolve("project-four", ".")
		Expect(err).ToNot(HaveOccurred())
		Expect(fq).To(Equal("project-four/four-main.js"))
	})

	It("should cache simple lookups", func() {
		resolver.Resolve("project-three/a", ".")
		Expect(repo.checked).To(HaveLen(5))
		Expect(repo.opened).To(HaveLen(0))

		resolver.Resolve("./a", "project-three")
		Expect(repo.checked).To(HaveLen(5))
		Expect(repo.opened).To(HaveLen(0))
	})
})
