fang a = true;
fang b = false;
if (a) {
  if (b) { hiss("ab"); }
  else { hiss("a only"); } // expect: a only
}