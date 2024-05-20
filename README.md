# lopix

[![Go Reference](https://pkg.go.dev/badge/github.com/tinne26/lopix.svg)](https://pkg.go.dev/github.com/tinne26/lopix)

A package to assist the development of very simple [Ebitengine](https://ebitengine.org) pixel art games.

This package allows you to implement your game using a logical resolution and ignoring `Game.Layout()` completely. Scaling is managed internally by `lopix` with pixel art aware scaling algorithms.

## Context

There are multiple models for pixel art games:
- The simplest model consists in aligning both the graphical elements and the game view to a logical pixel grid. In most cases, these games don't have a notion of "camera" and use fixed screens. This is what `lopix` implements.
- The next model aligns graphical elements to a logical pixel grid, but gives the camera the freedom to move smoothly through non-integer coordinates. Pixel art RPGs would be a common genre for this. This is what [`mipix`](https://github.com/tinne26/mipix) implements.
- The final and more modern approach revolves around floating point positioning for both graphical elements and the camera. Different strategies are possible when it comes to projecting from the "logical space" to the final high resolution canvas. Most pixel art games use this approach nowadays... at least when they aren't doing outright 3D already.

My original intent was to create a library that would help with all the pixel-art game models... but it soon became obvious that the workflows are different enough to warrant different packages for each one.

All in all, it must be explicitly said that lopix implements a very limited and restrictive model; the package might have more educative value than practical value. Don't be afraid to look into [`mipix`](https://github.com/tinne26/mipix) or roll your own solutions if you need to.

## Code example

https://github.com/tinne26/lopix-examples

