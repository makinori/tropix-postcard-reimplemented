# Tropix Postcard Reimplemented

-   https://tropixgame.com
-   https://tropix.hotmilk.space ([github repo](https://github.com/makinori/tropixgame.com))

Fixes postcards in Tropix 1 and 2. Written in Go and cross compiles easily.

<!-- Can also send to multiple emails if you seperate in-game with `^`. -->

-   Copy `config.example.json` to `config.json` and fill out
-   Run `build.sh` or `build.bat` to get `Postcard.exe`
-   Place in `C:\Program Files (x86)\Tropix\`

Test by running:

`go run main.go _ _ _ email@example.com Name assets/tropixTitle.jpg assets/tropixTitle.jpg`

The last two arguments are usually temporary paths to the front and back postcard images that the game generates.

![Example](example.jpg)
