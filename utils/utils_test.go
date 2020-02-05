package utils

import "testing"

func TestUtil(t *testing.T) {
	t.Run("defaultGOPATH", func(t *testing.T) {
		str := defaultGOPATH()
		if str == "" {
			t.Fatalf("Failed")
		}
	})

	t.Run("GetGOPATHs", func(t *testing.T) {
		str := GetGOPATHs()
		if len(str) == 0 {
			t.Fatalf("Failed")
		}
	})

	t.Run("CheckEnv", func(t *testing.T) {
		_, _, err := CheckEnv("hello")
		if err != nil {
			t.Fatal(err)
		}
	})

}
