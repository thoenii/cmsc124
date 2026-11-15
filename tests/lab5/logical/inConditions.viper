if (1 < 2 and 2 < 3) {
  hiss("both"); // expect: both
}

if (false or false) {
  hiss("no");
} else {
  hiss("neither"); // expect: neither
}

fang i = 0;
slither (i < 5 and i != 2) {
  hiss(i); // expect: 0
           // expect: 1
  i = i + 1;
}