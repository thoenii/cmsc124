fang x = "a";
{
  fang x = "b";
  {
    fang x = "c";
    hiss(x);   // expect: c
  }
  hiss(x);     // expect: b
}
hiss(x);       // expect: a