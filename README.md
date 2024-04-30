# lopix

[![Go Reference](https://pkg.go.dev/badge/github.com/tinne26/lopix.svg)](https://pkg.go.dev/github.com/tinne26/lopix)

**(VERY WIP, VERY EARLY STAGES, VERY UNFINISHED)**

A package to assist the development of simple [Ebitengine](https://ebitengine.org) pixel art games.

This package allows you to implement your game using a logical resolution and ignoring `Game.Layout()` completely. Scaling is managed internally by `lopix` with a pixel art aware scaling algorithm.

## Context

There are multiple models for pixel art games:
- The simplest model consists in aligning both the graphical elements and the "game camera" to a logical pixel grid. In many cases, these games don't even have explicit camera systems, or use fixed screens.
- The next model aligns graphical elements to a logical pixel grid, but gives the camera the freedom to move more smoothly (in the high resolution plane instead of the logical one). Pixel art RPGs would be a common genre for this.
- The final and more modern approach revolves around floating point positioning for both graphical elements and the camera. Different strategies are possible when it comes to projecting from the "logical space" to the final high resolution canvas. Most pixel art games use this approach nowadays... at least when they aren't doing outright 3D already.

My original intent was to create a library that would help with all the pixel-art game models... but it soon became obvious that the workflows are different enough to warrant different packages for each one. This `lopix` package takes care of the first type described.

## Code example

https://github.com/tinne26/lopix-examples

