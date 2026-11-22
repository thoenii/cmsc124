// natives_game.go connects Viper to package game. World holds a swappable
// pointer to the current game; the natives are thin wrappers that validate
// arguments and call one game method. They read World.G at call time so
// setup() can replace the game.
package evaluator
