# Color IDE

Source code text editor



## Build on Mac

```
brew install sdl2 sdl2_image sdl2_ttf sdl2_mixer
```

Install Go.

Then:
```
make
```


## Notes

If compiled on Mac, you may see this warning:
```
ld: warning: ignoring duplicate libraries: '-lSDL2', '-lSDL2_ttf'
```
This does not affect the application, but you can add to your `.zprofile`:
```
export CGO_LDFLAGS="-Wl,-no_warn_duplicate_libraries"
```
