fang a = "outer";
{
  fang a = a;
  hiss(a);   // expect: outer
}