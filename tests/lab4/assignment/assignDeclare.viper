fang x = "outer";
{
  fang x = "inner";
  x = "changed";
  hiss(x);   // expect: changed
}
hiss(x);     // expect: outer
{
  x = "changed";
}
hiss(x);     // expect: changed