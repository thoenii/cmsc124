fang i = 0;
slither (i < 3) {
  hiss(i); // expect: 0
  // expect: 1
  // expect: 2
  i = i + 1;
}