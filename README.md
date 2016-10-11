# godgt

An experimental utility for interacting with DGT's Electronic Chess
Boards, written in Go. The development platform is Linux; on other
platforms, YMMV.

There are two executables, one (`rawdump`) for dumping out more or
less raw position updates from the board, and another (`dgtd`) for
attempting to turn these into actual moves.

The latter task is not as straightforward as it might seem. The reason
for this is that the board does not "understand" chess; it just
reports low-level piece move events.

For example, the move "exd5" might be received as three events, as
follows:

1. EMPTY to e4 (White picks up his pawn on e4)
2. EMPTY to d5 (White removes the captured black pawn on d5)
3. WHITE PAWN to d5 (White places his pawn on d5)

Worse still, because of the sequential nature of the scanning of the
board, these might be received out of order. Out of order events can
also occur due to player behaviour; for example, removing a captured
piece before moving the capturing piece.

However, a long slow slide of a piece along a rank, file or diagonal
can produce a "storm" of updates as several squares along the path
detect the temporary presence of the piece over the square.

Even worse is the situation when a player drops one or several pieces,
or makes and retracts an invalid move.

The code to turn the events into moves is quite sophisticated, but
it's not particularly robust; it can easily get itself into a
situation whereby it considers any additional input from the board to
be invalid.

For this reason, if you want to simply log moves from the board your
best best is currently to run the rawdump utility and reconstruct the
moves manually after the fact.

## Installation

```
go get github.com/kgigitdev/godgt
```

## Dependencies

```
go get github.com/jacobsa/go-serial/serial
go get github.com/jessevdk/go-flags
go get github.com/malbrecht/chess
```

## Information

Luckily, the DGT boards don't need any additional driver under Linux:
plugging one in should automatically load the FTDI USB-serial driver
(ftdi_sio and usbserial) on a reasonably up-to-date Linux
installation. You should also see a serial device like the following
appear:

```
$ ls -la /dev/ttyUSB*
crw-rw---- 1 root dialout 188, 0 Apr 11 18:52 /dev/ttyUSB0
```

In order to make this readable by you, you might need to add yourself
to the group owning the device, using a command like the following,
after which you will need to log out and back in.

```
sudo usermod -a -G dialout ${USER}
```

## Compiling and running

```
cd rawdump
go build
./rawdump --pngs
```

## Output

Log output looks like:

```
2016/10/11 22:39:05 BOARD: r1bqkbnr/pp1p1ppp/2p5/n3p3/P3P3/2PP1N2/RP3PPP/1NBQKB1R w K a1 0 0
2016/10/11 22:39:05 r bqkbnr
2016/10/11 22:39:05 pp p ppp
2016/10/11 22:39:05   p     
2016/10/11 22:39:05 n   p   
2016/10/11 22:39:05 P   P   
2016/10/11 22:39:05   PP N  
2016/10/11 22:39:05 RP   PPP
2016/10/11 22:39:05  NBQKB R
```

The `--pngs` option will also cause one `.png` file of the position to
be created for each field update:

```
$ ls *.png
boardupdate-0001.png  boardupdate-0008.png  boardupdate-0015.png
boardupdate-0002.png  boardupdate-0010.png  boardupdate-0016.png
boardupdate-0004.png  boardupdate-0012.png
boardupdate-0006.png  boardupdate-0014.png
```

Note that there will be many *more* `.png` files than there are moves in
the game. However, it's *much* easier to reconstruct the game from these
`.png` files than from the FEN strings.

## Embedded Assets

The `assets/` directory contains `.png` image files for the
PNG writer, as well as `.svg` files.

These were created as described in the file `assets/images/README`

These are not used directly at runtime; instead, they are bundled into
a static filesystem, in `assets.go`, using the excellent `enc`
utility:

```
go install github.com/mjibson/esc
esc -o assets.go -pkg godgt assets
```

You should not need this utility unless you want to recreate
`assets.go`.

## History

This is something like my fourth attempt to write this utility. The
first couple got bogged down by trying to write a general-purpose
chess library in Go, in the spirit of `python-chess`.

The third attempt gave up on that and used `malbrecht/chess` instead,
but also got bogged down attempting to be a general-purpose bridge
between the DGT board, xboard, FICS, and stockfish, as well as hosting
an HTTP server for viewing the game in progress.

This new, much less ambitious, attempt is cobbled together from pieces
of the previous failed attempts. Its code style is therefore rather
haphazard. Notably, most of the older salvaged code was written in a
"everything is a pointer to a struct" object-oriented-ish style; later
code is consciously experimenting with not doing that.

