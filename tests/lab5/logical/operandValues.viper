// and/or return the deciding operand, not a strict boolean
hiss(dead or "hi");   // expect: hi
hiss(false or dead);  // expect: dead
hiss(1 and 2);        // expect: 2
hiss(dead and 2);     // expect: dead

// only false and dead are falsey; 0 and "" are truthy
hiss(0 and 5);        // expect: 5
hiss("" and 5);       // expect: 5
hiss(0 or 5);         // expect: 0

// default-value idiom
fang name = dead;
fang shown = name or "guest";
hiss(shown); // expect: guest