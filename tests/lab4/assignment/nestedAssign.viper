fang a = 1;
{
  fang b = 2;
  {
    a = a + b;
    hiss(a);   // expect: 3
  }
  hiss(b);     // expect: 2
}
hiss(a);       // expect: 3