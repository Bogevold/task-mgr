package task_test

import (
	"testing"
	"time"

	"github.com/bogevold/task-mgr/internal/task"
)

func TestSaveAndGet(t *testing.T) {
	store := task.NewInMemoryStore()
	// 1. Lag en Task
	title := "t1"
	var tsk = task.Task{Title: title, CreatedAt: time.Now()}
	// 2. Lagre den med Save
	saved, err := store.Save(tsk)
	if err != nil {
		t.Fatalf("forventet ingen feil, fikk: %v", err)
	}
	// 3. Hent den ut med Get på ID
	gotTask, err := store.Get(saved.ID)
	// 4. Sjekk at Title stemmer
	if gotTask.Title != saved.Title {
		t.Fatalf("forventet identisk tittel fikk: lagret=%v; get=%v", saved.Title, gotTask.Title)
	}
}

func TestDelete(t *testing.T) {
	store := task.NewInMemoryStore()
	// Lag en task som kan slettes
	saved, err := store.Save(task.Task{Title: "Slett meg", CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("forventet ingen feil ved oppretting, fikk: %v", err)
	}
	// Slett task
	err = store.Delete(saved.ID)
	if err != nil {
		t.Fatalf("forventet ingen feil ved sletting, fikk: %v", err)
	}
	gotTask, err := store.Get(saved.ID)
	if err == nil {
		t.Fatalf("store inneholder fortsatt task ID=%v", gotTask.ID)
	}
}

func TestUpdate(t *testing.T) {
	store := task.NewInMemoryStore()
	// Lag en task som kan oppdateres
	saved, err := store.Save(task.Task{Title: "Oppdater meg", CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("forventet ingen feil ved oppretting, fikk: %v", err)
	}
	saved.Title = "Oppdatert"
	upd, err := store.Update(saved)
	if err != nil {
		t.Fatalf("forventet ikke feil ved oppdatering, fikk: %v", err)
	}
	if upd.Title != saved.Title {
		t.Fatalf("Title avviker %v <-> %v", upd.Title, saved.Title)
	}
}
