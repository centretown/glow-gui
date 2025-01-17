package ui

// func test_editor_setup(t *testing.T) (effect control.Control, e *LayerEditor, err error) {
// 	app := test.NewApp()
// 	w := app.NewWindow("Editor")
// 	effect = control.Neweffect(fileio.NewStore(app.Preferences()))
// 	toolBar := NewSharedTools(effect)
// 	e = NewLayerEditor(effect, w, toolBar)
// 	return
// }

// // func test_apply_button_disabled(t *testing.T, e *LayerEditor, expected bool) {
// // 	b := e.applyButton.Disabled()
// // 	if b != expected {
// // 		t.Fatalf("apply button got %v expected %v", b, expected)
// // 	}
// // }

// // func test_revert_button_disabled(t *testing.T, e *LayerEditor, expected bool) {
// // 	b := e.revertButton.Disabled()
// // 	if b != expected {
// // 		t.Fatalf("revert button got %v expected %v", b, expected)
// // 	}
// // }

// func test_dirty(t *testing.T, e *LayerEditor, expected bool) {
// 	b := isDirty(e)
// 	if b != expected {
// 		t.Fatalf("dirty got %v expected %v", b, expected)
// 	}
// }

// func test_layer_editor_init(t *testing.T) (e *LayerEditor) {
// 	_, e, err := test_editor_setup(t)
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}

// 	// test_apply_button_disabled(t, e, true)
// 	// test_revert_button_disabled(t, e, true)
// 	test_dirty(t, e, false)

// 	return
// }

// func testScan(t *testing.T, e *LayerEditor, expected int) {
// 	scan, _ := e.fields.Scan.Get()
// 	if scan != expected {
// 		t.Fatalf("expected %d value %d", expected, scan)
// 	}
// }

// func testHue(t *testing.T, e *LayerEditor, expected int) {
// 	shift, _ := e.fields.HueShift.Get()
// 	if shift != expected {
// 		t.Fatalf("expected %d value %d", expected, shift)
// 	}
// }

// func testRate(t *testing.T, e *LayerEditor, expected int) {
// 	rate, _ := e.fields.Rate.Get()
// 	if rate != expected {
// 		t.Fatalf("expected %d value %d", expected, rate)
// 	}
// }

// func isDirty(le *LayerEditor) bool {
// 	return le.effect.HasChanged()
// }

// func TestLayerEditor(t *testing.T) {
// 	e := test_layer_editor_init(t)

// 	if isDirty(e) {
// 		t.Fatal("Dirty")
// 	}

// 	testScan(t, e, e.scanBounds.OffVal)
// 	testHue(t, e, e.hueBounds.OffVal)
// 	testRate(t, e, e.rateBounds.OffVal)

// 	e.checkScan.SetChecked(true)
// 	testScan(t, e, e.scanBounds.OnVal)
// 	e.checkScan.SetChecked(false)
// 	testScan(t, e, e.scanBounds.OffVal)

// 	e.checkHue.SetChecked(true)
// 	testHue(t, e, e.hueBounds.OnVal)
// 	e.checkHue.SetChecked(false)
// 	testRate(t, e, e.hueBounds.OffVal)

// 	e.checkRate.SetChecked(true)
// 	testRate(t, e, e.rateBounds.OnVal)
// 	e.checkRate.SetChecked(false)
// 	testRate(t, e, e.rateBounds.OffVal)
// }
