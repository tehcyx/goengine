# goengine
golang voxel engine

Look into this: https://github.com/raedatoui/learn-opengl-golang

# How to contribute

- Fork the project, add your code & create a pull request.
- Create an [issue](https://github.com/tehcyx/goengine/issues)

### Pre requisites
...

## Compile & run
...

## Cross compile macOS to Windows
1. Install [Homebrew](https://brew.sh/)
2. Install MinGW through Homebrew via `brew install mingw-w64`
3. Download the SDL2 development package for MinGW [here](http://libsdl.org/download-2.0.php) (and the others like SDL_image, SDL_mixer, etc.. [here](https://www.libsdl.org/projects/) if you use them).
4. Extract the SDL2 development package and copy the `x86_64-w64-mingw` folder inside recursively to the system's MinGW `x86_64-w64-mingw32` folder. You may also do the same for the `i686-w64-mingw32` folder. The path to MinGW may be slightly different but the command should look something like `cp -r x86_64-w64-mingw32 /usr/local/Cellar/mingw-w64/5.0.3_2/toolchain-x86_64`.
Now you can start cross-compiling your Go program by running `env CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows CGO_LDFLAGS="-L/usr/local/Cellar/mingw-w64/5.0.3/toolchain-x86_64/x86_64-w64-mingw32/lib -lSDL2" CGO_CFLAGS="-I/usr/local/Cellar/mingw-w64/5.0.3/toolchain-x86_64/x86_64-w64-mingw32/include -D_REENTRANT" go build -x main.go`. You can change some of the parameters if you'd like to. In this example, it should produce a main.exe executable file.
5. Before running the program, you need to put SDL2.dll from the [SDL2 runtime package](http://libsdl.org/download-2.0.php) (For others like SDL_image, SDL_mixer, etc.., look for them [here](https://www.libsdl.org/projects/)) for Windows in the same folder as your executable.
6. Now you should be able to run the program using Wine or Windows!

# Future plans (in no particular order)

- Use fonts for UI:
    - http://www.dafont.com/5x5-square.font
    - http://www.dafont.com/5x5-rounded.font
- Implement proper mesh rendering
- Implement chunk handling
- Implement terrain generation
    - Restrict terrain to be more of a desert style
- Character models
    - Customization
- Items
    - stat system