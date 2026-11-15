// or skips the right side when the left is truthy
hiss(true or (1 / 0)); // expect: true
hiss(1 or (1 / 0));    // expect: 1

fang x = 0;
true or (x = 1);
hiss(x); // expect: 0

// left is falsey, so the right side does run
false or (x = 1);
hiss(x); // expect: 1