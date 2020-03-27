package main

import (
	"testing"
)

type fakeKiller struct {
	killed bool
}

func (killer *fakeKiller) Kill(message string) {
	killer.killed = true
}

func Test_validateCli_When_cs_Is_Empty_String_It_Kills_App(t *testing.T) {
	// arrange
	fakeKiller := &fakeKiller{killed: false}
	cs, dbName, scripts := "", "dbName", "scripts"

	// act
	validateCli(fakeKiller, cs, dbName, scripts)

	// assert
	if !fakeKiller.killed {
		t.Error("Expected Kill() to have been called")
	}
}

func Test_validateCli_When_dbName_Is_Empty_String_It_Kills_App(t *testing.T) {
	// arrange
	fakeKiller := &fakeKiller{killed: false}
	cs, dbName, scripts := "cs", "", "scripts"

	// act
	validateCli(fakeKiller, cs, dbName, scripts)

	// assert
	if !fakeKiller.killed {
		t.Error("Expected Kill() to have been called")
	}
}

func Test_validateCli_When_scripts_Is_Empty_String_It_Kills_App(t *testing.T) {
	// arrange
	fakeKiller := &fakeKiller{killed: false}
	cs, dbName, scripts := "cs", "dbName", ""

	// act
	validateCli(fakeKiller, cs, dbName, scripts)

	// assert
	if !fakeKiller.killed {
		t.Error("Expected Kill() to have been called")
	}
}
