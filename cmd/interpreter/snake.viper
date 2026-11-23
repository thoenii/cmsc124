setup(12, 8);
fang running = true;

slither (running) {
  render();
  hiss("w/a/s/d + Enter to steer, q to quit");
  fang key = input();

  if (key == "q") {
    running = false;
  } else {
    if (key == "w") { turn("up"); }
    if (key == "s") { turn("down"); }
    if (key == "a") { turn("left"); }
    if (key == "d") { turn("right"); }

    move();

    if (over()) {
      render();
      hiss("Game over! Final score: " + str(score()));
      running = false;
    }
  }
}
