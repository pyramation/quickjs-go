package quickjs

import (
	"testing"
)

func TestNewWithoutConstructor(t *testing.T) {
	rt := NewRuntime()
	defer rt.Close()
	ctx := rt.NewContext()
	defer ctx.Close()

	_, err := ctx.Eval(`
		class TestClass {
			constructor() {
				this.constructorCalled = true;
				this.value = 42;
			}
			
			getConstructorCalled() {
				return this.constructorCalled;
			}
		}
	`)
	if err != nil {
		t.Fatalf("Failed to define test class: %v", err)
	}

	testClass := ctx.Globals().Get("TestClass")
	defer testClass.Free()

	if !testClass.IsConstructor() {
		t.Fatal("TestClass should be a constructor")
	}

	normalInstance := testClass.New()
	defer normalInstance.Free()
	
	constructorCalledNormal := normalInstance.Get("constructorCalled")
	defer constructorCalledNormal.Free()
	
	if !constructorCalledNormal.ToBool() {
		t.Error("Normal constructor call should set constructorCalled to true")
	}

	noConstructorInstance := testClass.NewWithoutConstructor()
	defer noConstructorInstance.Free()
	
	constructorCalledSkipped := noConstructorInstance.Get("constructorCalled")
	defer constructorCalledSkipped.Free()
	
	if !constructorCalledSkipped.IsUndefined() {
		t.Error("NewWithoutConstructor should not call constructor, constructorCalled should be undefined")
	}

	prototypeMethod := noConstructorInstance.Get("getConstructorCalled")
	defer prototypeMethod.Free()
	
	if !prototypeMethod.IsFunction() {
		t.Error("Object created without constructor should still have prototype methods")
	}
}

func TestNewWithoutConstructorError(t *testing.T) {
	rt := NewRuntime()
	defer rt.Close()
	ctx := rt.NewContext()
	defer ctx.Close()

	notConstructor := ctx.String("not a constructor")
	defer notConstructor.Free()

	result := notConstructor.NewWithoutConstructor()
	defer result.Free()

	if !result.IsException() {
		t.Error("NewWithoutConstructor should return error for non-constructor")
	}
}
