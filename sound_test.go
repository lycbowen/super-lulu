package main

import "testing"

func TestSynthChirpProducesStereoPCM(t *testing.T) {
	data := synthChirp(0.05, 220, 440, 0.5, waveSine)
	if len(data) == 0 {
		t.Fatal("expected sound data")
	}
	if len(data)%4 != 0 {
		t.Fatalf("expected 16-bit stereo PCM length, got %d", len(data))
	}
}
