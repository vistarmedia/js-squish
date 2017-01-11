package jssquish

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robertkrimen/otto/parser"
)

var _ = Describe("Require Visitor", func() {

	var visitor *RequireVisitor

	parse := func(src string) {
		program, err := parser.ParseFile(nil, "", src, 0)
		Expect(err).ToNot(HaveOccurred())

		visitor = NewRequireVisitor()
		for _, stmt := range program.Body {
			err = WalkNode(visitor, stmt)
			Expect(err).ToNot(HaveOccurred())
		}
	}

	Describe("an import call", func() {
		BeforeEach(func() {
			parse(`require("cool-mans")`)
		})

		It("should be found", func() {
			found := visitor.Requires()
			Expect(found).To(And(
				HaveLen(1),
				ContainElement("cool-mans"),
			))
		})
	})

	Describe("a variable assignment", func() {

		BeforeEach(func() {
			parse(`var time = require('wristwatch')`)
		})

		It("should be found", func() {
			found := visitor.Requires()
			Expect(found).To(And(
				HaveLen(1),
				ContainElement("wristwatch"),
			))
		})

	})

	Describe("module export", func() {

		BeforeEach(func() {
			parse(`module.exports = require('cool.thing.dude');`)
		})

		It("should be found", func() {
			found := visitor.Requires()
			Expect(found).To(And(
				HaveLen(1),
				ContainElement("cool.thing.dude"),
			))
		})

	})

	// TODO: This should probably be rewritten to do the resolution and not
	// require both to be included
	Describe("try/catch require", func() {

		BeforeEach(func() {
			parse(`
				var lib;
				try {
					lib = require('primary');
				} catch(e) {
					lib = require('secondary');
				}
			`)
		})

		It("should be found", func() {
			found := visitor.Requires()
			Expect(found).To(And(
				HaveLen(2),
				ContainElement("primary"),
				ContainElement("secondary"),
			))
		})

	})

	Describe("Object literal", func() {

		BeforeEach(func() {
			parse(`
        var coolThings = {
          "a": require('pants')()
        };
			`)
		})

		It("should be found", func() {
			found := visitor.Requires()
			Expect(found).To(And(
				HaveLen(1),
				ContainElement("pants"),
			))
		})
	})

})
