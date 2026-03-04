package hierarchy

import (
	"errors"
	"testing"

	"github.com/adm87/stellar/errs"
	"github.com/yohamta/donburi"
)

func TestSetParentAttachDetach(t *testing.T) {
	world := donburi.NewWorld()
	parent := world.Create(Component)
	child := world.Create(Component)

	if err := SetParent(world, child, parent); err != nil {
		t.Fatalf("SetParent attach failed: %v", err)
	}

	childModel := Component.Get(world.Entry(child))
	parentModel := Component.Get(world.Entry(parent))

	if childModel.Parent != parent {
		t.Fatalf("child parent = %v, want %v", childModel.Parent, parent)
	}
	if parentModel.FirstChild != child {
		t.Fatalf("parent first child = %v, want %v", parentModel.FirstChild, child)
	}

	if err := SetParent(world, child, donburi.Null); err != nil {
		t.Fatalf("SetParent detach failed: %v", err)
	}

	if childModel.Parent != donburi.Null {
		t.Fatalf("child parent = %v, want null", childModel.Parent)
	}
	if childModel.NextSibling != donburi.Null {
		t.Fatalf("child next sibling = %v, want null", childModel.NextSibling)
	}
	if childModel.PrevSibling != donburi.Null {
		t.Fatalf("child prev sibling = %v, want null", childModel.PrevSibling)
	}
	if parentModel.FirstChild != donburi.Null {
		t.Fatalf("parent first child = %v, want null", parentModel.FirstChild)
	}
}

func TestSetParentSiblingLinks(t *testing.T) {
	world := donburi.NewWorld()
	parent := world.Create(Component)
	childA := world.Create(Component)
	childB := world.Create(Component)

	if err := SetParent(world, childA, parent); err != nil {
		t.Fatalf("SetParent childA failed: %v", err)
	}
	if err := SetParent(world, childB, parent); err != nil {
		t.Fatalf("SetParent childB failed: %v", err)
	}

	parentModel := Component.Get(world.Entry(parent))
	childAModel := Component.Get(world.Entry(childA))
	childBModel := Component.Get(world.Entry(childB))

	if parentModel.FirstChild != childB {
		t.Fatalf("parent first child = %v, want %v", parentModel.FirstChild, childB)
	}
	if childBModel.NextSibling != childA {
		t.Fatalf("childB next sibling = %v, want %v", childBModel.NextSibling, childA)
	}
	if childAModel.PrevSibling != childB {
		t.Fatalf("childA prev sibling = %v, want %v", childAModel.PrevSibling, childB)
	}

	if err := SetParent(world, childB, donburi.Null); err != nil {
		t.Fatalf("detach childB failed: %v", err)
	}

	if parentModel.FirstChild != childA {
		t.Fatalf("parent first child = %v, want %v", parentModel.FirstChild, childA)
	}
	if childAModel.PrevSibling != donburi.Null {
		t.Fatalf("childA prev sibling = %v, want null", childAModel.PrevSibling)
	}
}

func TestSetParentRejectsSelfAndCycles(t *testing.T) {
	world := donburi.NewWorld()
	a := world.Create(Component)
	b := world.Create(Component)
	c := world.Create(Component)

	if err := SetParent(world, b, a); err != nil {
		t.Fatalf("SetParent(b, a) failed: %v", err)
	}
	if err := SetParent(world, c, b); err != nil {
		t.Fatalf("SetParent(c, b) failed: %v", err)
	}

	if err := SetParent(world, a, a); err == nil {
		t.Fatal("expected self-parent to fail")
	} else {
		var opErr errs.InvalidOperation
		if !errors.As(err, &opErr) {
			t.Fatalf("expected InvalidOperation, got %T", err)
		}
	}

	if err := SetParent(world, a, c); err == nil {
		t.Fatal("expected cycle prevention to fail")
	} else {
		var opErr errs.InvalidOperation
		if !errors.As(err, &opErr) {
			t.Fatalf("expected InvalidOperation, got %T", err)
		}
	}
}

func BenchmarkSetParentAttachDetach(b *testing.B) {
	world := donburi.NewWorld()
	parent := world.Create(Component)
	child := world.Create(Component)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := SetParent(world, child, parent); err != nil {
			b.Fatalf("attach failed: %v", err)
		}
		if err := SetParent(world, child, donburi.Null); err != nil {
			b.Fatalf("detach failed: %v", err)
		}
	}
}

func BenchmarkSetParentReparent(b *testing.B) {
	world := donburi.NewWorld()
	parentA := world.Create(Component)
	parentB := world.Create(Component)
	child := world.Create(Component)

	if err := SetParent(world, child, parentA); err != nil {
		b.Fatalf("initial attach failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := SetParent(world, child, parentB); err != nil {
			b.Fatalf("reparent to B failed: %v", err)
		}
		if err := SetParent(world, child, parentA); err != nil {
			b.Fatalf("reparent to A failed: %v", err)
		}
	}
}
