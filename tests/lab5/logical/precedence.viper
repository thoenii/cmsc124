// and binds tighter than or
hiss(true or false and false);    // expect: true
hiss((true or false) and false);  // expect: false
hiss(false and true or true);     // expect: true

// comparison and equality bind tighter than and/or
hiss(1 < 2 and 2 < 3);            // expect: true
hiss(1 == 2 or 2 == 2);           // expect: true

// assignment binds looser than or
fang b = false or "set";
hiss(b); // expect: set