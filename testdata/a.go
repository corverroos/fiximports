package testdata

import (
	"testing"
	"github.com/sebdah/goldie/v2"
	
	_ "fmt"
	"github.com/drand/kyber-bls12381"
)

func TestFoo(t *testing.T) {
	goldie.New(t)
	bls.NewGroupG1()
}
