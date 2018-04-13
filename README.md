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
2. Install go dependencies

    ```
    go get github.com/go-gl/gl/v4.1-core/gl
    go get github.com/go-gl/mathgl/mgl32
    go get -v github.com/veandco/go-sdl2/{sdl,mix,img,ttf}
    ```
3. Install SDL2 via brew: `brew install sdl2{,_image,_ttf,_mixer} pkg-config`
4. run `make`
5. run `./bin/app`
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