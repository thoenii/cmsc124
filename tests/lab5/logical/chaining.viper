// chains are left-associative and stop at the first decider
hiss(false or false or 3);    // expect: 3
hiss(1 and 2 and 3);          // expect: 3
hiss(1 and dead and 3);       // expect: dead
hiss(dead or false or dead);  // expect: dead