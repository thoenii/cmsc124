fang x = "outer";
{
  fang x = "inner";
  hiss(x);   // expect: inner
}
hiss(x);     // expect: outer