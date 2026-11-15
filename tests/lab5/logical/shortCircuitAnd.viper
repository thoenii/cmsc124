// and skips the right side when the left is falsey
hiss(false and (1 / 0)); // expect: false
hiss(dead and (1 / 0));  // expect: dead

// side effect proves the right side never ran
fang x = 0;
false and (x = 1);
hiss(x); // expect: 0

// left is truthy, so the right side does run
true and (x = 1);
hiss(x); // expect: 1