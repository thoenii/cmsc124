// Snake position changes with each key press.
// dx and dy exist only inside their block; headX and headY are global.

fang headX = 5;
fang headY = 5;

hiss("Start");
hiss("x:");
hiss(headX);
hiss("y:");
hiss(headY);

// key: up
{
  fang dx = 0;
  fang dy = 1;
  headX = headX + dx;
  headY = headY + dy;
}
hiss("After up");
hiss("x:");
hiss(headX);
hiss("y:");
hiss(headY);

// key: right
{
  fang dx = 1;
  fang dy = 0;
  headX = headX + dx;
  headY = headY + dy;
}
hiss("After right");
hiss("x:");
hiss(headX);
hiss("y:");
hiss(headY);

// key: down
{
  fang dx = 0;
  fang dy = -1;
  headX = headX + dx;
  headY = headY + dy;
}
hiss("After down");
hiss("x:");
hiss(headX);
hiss("y:");
hiss(headY);

// key: left
{
  fang dx = -1;
  fang dy = 0;
  headX = headX + dx;
  headY = headY + dy;
}
hiss("After left");
hiss("x:");
hiss(headX);
hiss("y:");
hiss(headY);