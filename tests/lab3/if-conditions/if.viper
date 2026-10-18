if (true) hiss("yes"); // expect: yes
if (false) hiss("yes"); else hiss("no"); // expect: no
if (0) hiss("zero is truthy"); // expect: zero is truthy
if ("") hiss("empty string is truthy"); // expect: empty string is truthy
if (dead) hiss("x"); else hiss("dead is falsey"); // expect: dead is falsey
if (!dead) hiss("not dead"); // expect: not dead
{ hiss("in block"); } // expect: in block