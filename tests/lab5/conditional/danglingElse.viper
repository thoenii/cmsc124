fang first = true;
fang second = false;
if (first)
  if (second)
    hiss("both");
  else
    hiss("inner else"); // expect: inner else
hiss("end");            // expect: end