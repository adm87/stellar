package hierarchy

import (
	"errors"
	"testing"

	"github.com/adm87/stellar/errs"
	"github.com/yohamta/donburi"
)

var benchmarkEntitySink donburi.Entity
var benchmarkCountSink int

func TestSetParentAttachDetach(t *testing.T) {
	world := donburi.NewWorld()
	parent := world.Create(HierarchyComponent)
	child := world.Create(HierarchyComponent)
	parentEntry := world.Entry(parent)
	childEntry := world.Entry(child)

	if err := SetParent(world, childEntry, parentEntry); err != nil {
		t.Fatalf("SetParent attach failed: %v", err)
	}

	childModel := HierarchyComponent.Get(childEntry)
	parentModel := HierarchyComponent.Get(parentEntry)

	if childModel.parent != parent {
		t.Fatalf("child parent = %v, want %v", childModel.parent, parent)
	}
	if parentModel.firstChild != child {
		t.Fatalf("parent first child = %v, want %v", parentModel.firstChild, child)
	}

	if err := SetParent(world, childEntry, nil); err != nil {
		t.Fatalf("SetParent detach failed: %v", err)
	}

	if childModel.parent != donburi.Null {
		t.Fatalf("child parent = %v, want null", childModel.parent)
	}
	if childModel.nextSibling != donburi.Null {
		t.Fatalf("child next sibling = %v, want null", childModel.nextSibling)
	}
	if childModel.prevSibling != donburi.Null {
		t.Fatalf("child prev sibling = %v, want null", childModel.prevSibling)
	}
	if parentModel.firstChild != donburi.Null {
		t.Fatalf("parent first child = %v, want null", parentModel.firstChild)
	}
}

func TestSetParentSiblingLinks(t *testing.T) {
	world := donburi.NewWorld()
	parent := world.Create(HierarchyComponent)
	childA := world.Create(HierarchyComponent)
	childB := world.Create(HierarchyComponent)
	parentEntry := world.Entry(parent)
	childAEntry := world.Entry(childA)
	childBEntry := world.Entry(childB)

	if err := SetParent(world, childAEntry, parentEntry); err != nil {
		t.Fatalf("SetParent childA failed: %v", err)
	}
	if err := SetParent(world, childBEntry, parentEntry); err != nil {
		t.Fatalf("SetParent childB failed: %v", err)
	}

	parentModel := HierarchyComponent.Get(parentEntry)
	childAModel := HierarchyComponent.Get(childAEntry)
	childBModel := HierarchyComponent.Get(childBEntry)

	if parentModel.firstChild != childB {
		t.Fatalf("parent first child = %v, want %v", parentModel.firstChild, childB)
	}
	if childBModel.nextSibling != childA {
		t.Fatalf("childB next sibling = %v, want %v", childBModel.nextSibling, childA)
	}
	if childAModel.prevSibling != childB {
		t.Fatalf("childA prev sibling = %v, want %v", childAModel.prevSibling, childB)
	}

	if err := SetParent(world, childBEntry, nil); err != nil {
		t.Fatalf("detach childB failed: %v", err)
	}

	if parentModel.firstChild != childA {
		t.Fatalf("parent first child = %v, want %v", parentModel.firstChild, childA)
	}
	if childAModel.prevSibling != donburi.Null {
		t.Fatalf("childA prev sibling = %v, want null", childAModel.prevSibling)
	}
}

func TestSetParentRejectsSelfAndCycles(t *testing.T) {
	world := donburi.NewWorld()
	a := world.Create(HierarchyComponent)
	b := world.Create(HierarchyComponent)
	c := world.Create(HierarchyComponent)
	aEntry := world.Entry(a)
	bEntry := world.Entry(b)
	cEntry := world.Entry(c)

	if err := SetParent(world, bEntry, aEntry); err != nil {
		t.Fatalf("SetParent(b, a) failed: %v", err)
	}
	if err := SetParent(world, cEntry, bEntry); err != nil {
		t.Fatalf("SetParent(c, b) failed: %v", err)
	}

	if err := SetParent(world, aEntry, aEntry); err == nil {
		t.Fatal("expected self-parent to fail")
	} else {
		var opErr errs.InvalidOperation
		if !errors.As(err, &opErr) {
			t.Fatalf("expected InvalidOperation, got %T", err)
		}
	}

	if err := SetParent(world, aEntry, cEntry); err == nil {
		t.Fatal("expected cycle prevention to fail")
	} else {
		var opErr errs.InvalidOperation
		if !errors.As(err, &opErr) {
			t.Fatalf("expected InvalidOperation, got %T", err)
		}
	}
}

func TestGetParent(t *testing.T) {
	world := donburi.NewWorld()
	parent := world.Create(HierarchyComponent)
	child := world.Create(HierarchyComponent)
	parentEntry := world.Entry(parent)
	childEntry := world.Entry(child)

	if got := GetParent(world, childEntry); got != donburi.Null {
		t.Fatalf("initial parent = %v, want %v", got, donburi.Null)
	}

	if err := SetParent(world, childEntry, parentEntry); err != nil {
		t.Fatalf("attach failed: %v", err)
	}

	if got := GetParent(world, childEntry); got != parent {
		t.Fatalf("attached parent = %v, want %v", got, parent)
	}

	if err := SetParent(world, childEntry, nil); err != nil {
		t.Fatalf("detach failed: %v", err)
	}

	if got := GetParent(world, childEntry); got != donburi.Null {
		t.Fatalf("detached parent = %v, want %v", got, donburi.Null)
	}
}

func TestEachChildIteratesInInsertionHeadOrderAndSupportsBreak(t *testing.T) {
	world := donburi.NewWorld()
	parent := world.Create(HierarchyComponent)
	childA := world.Create(HierarchyComponent)
	childB := world.Create(HierarchyComponent)
	childC := world.Create(HierarchyComponent)

	parentEntry := world.Entry(parent)
	childAEntry := world.Entry(childA)
	childBEntry := world.Entry(childB)
	childCEntry := world.Entry(childC)

	if err := SetParent(world, childAEntry, parentEntry); err != nil {
		t.Fatalf("attach A failed: %v", err)
	}
	if err := SetParent(world, childBEntry, parentEntry); err != nil {
		t.Fatalf("attach B failed: %v", err)
	}
	if err := SetParent(world, childCEntry, parentEntry); err != nil {
		t.Fatalf("attach C failed: %v", err)
	}

	visited := make([]donburi.Entity, 0, 3)
	EachChild(world, parentEntry, func(entry *donburi.Entry) bool {
		visited = append(visited, entry.Entity())
		return true
	})

	if len(visited) != 3 {
		t.Fatalf("visited %d children, want 3", len(visited))
	}

	if visited[0] != childC || visited[1] != childB || visited[2] != childA {
		t.Fatalf("unexpected order: got [%v %v %v], want [%v %v %v]", visited[0], visited[1], visited[2], childC, childB, childA)
	}

	visited = visited[:0]
	EachChild(world, parentEntry, func(entry *donburi.Entry) bool {
		visited = append(visited, entry.Entity())
		return false
	})

	if len(visited) != 1 {
		t.Fatalf("early-break visited %d children, want 1", len(visited))
	}
	if visited[0] != childC {
		t.Fatalf("early-break first child = %v, want %v", visited[0], childC)
	}
}

func BenchmarkSetParentAttachDetach(b *testing.B) {
	world := donburi.NewWorld()
	parent := world.Create(HierarchyComponent)
	child := world.Create(HierarchyComponent)
	parentEntry := world.Entry(parent)
	childEntry := world.Entry(child)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := SetParent(world, childEntry, parentEntry); err != nil {
			b.Fatalf("attach failed: %v", err)
		}
		if err := SetParent(world, childEntry, nil); err != nil {
			b.Fatalf("detach failed: %v", err)
		}
	}
}

func BenchmarkSetParentReparent(b *testing.B) {
	world := donburi.NewWorld()
	parentA := world.Create(HierarchyComponent)
	parentB := world.Create(HierarchyComponent)
	child := world.Create(HierarchyComponent)
	parentAEntry := world.Entry(parentA)
	parentBEntry := world.Entry(parentB)
	childEntry := world.Entry(child)

	if err := SetParent(world, childEntry, parentAEntry); err != nil {
		b.Fatalf("initial attach failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := SetParent(world, childEntry, parentBEntry); err != nil {
			b.Fatalf("reparent to B failed: %v", err)
		}
		if err := SetParent(world, childEntry, parentAEntry); err != nil {
			b.Fatalf("reparent to A failed: %v", err)
		}
	}
}

func BenchmarkGetParent(b *testing.B) {
	world := donburi.NewWorld()
	parent := world.Create(HierarchyComponent)
	child := world.Create(HierarchyComponent)
	parentEntry := world.Entry(parent)
	childEntry := world.Entry(child)

	if err := SetParent(world, childEntry, parentEntry); err != nil {
		b.Fatalf("initial attach failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkEntitySink = GetParent(world, childEntry)
	}
}

func BenchmarkEachChildFullTraversal(b *testing.B) {
	world := donburi.NewWorld()
	parent := world.Create(HierarchyComponent)
	parentEntry := world.Entry(parent)

	for range 32 {
		child := world.Create(HierarchyComponent)
		if err := SetParent(world, world.Entry(child), parentEntry); err != nil {
			b.Fatalf("attach failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		EachChild(world, parentEntry, func(entry *donburi.Entry) bool {
			benchmarkEntitySink = entry.Entity()
			count++
			return true
		})
		benchmarkCountSink = count
	}
}

func BenchmarkEachChildBreakFirst(b *testing.B) {
	world := donburi.NewWorld()
	parent := world.Create(HierarchyComponent)
	parentEntry := world.Entry(parent)

	for range 32 {
		child := world.Create(HierarchyComponent)
		if err := SetParent(world, world.Entry(child), parentEntry); err != nil {
			b.Fatalf("attach failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		EachChild(world, parentEntry, func(entry *donburi.Entry) bool {
			benchmarkEntitySink = entry.Entity()
			count++
			return false
		})
		benchmarkCountSink = count
	}
}
