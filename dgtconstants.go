package godgt

/* Commands not resulting in returning messages:                            */

const DGT_SEND_RESET = 0x40

/* Puts the board in IDLE mode, cancelling any UPDATE mode.
 */

const DGT_TO_BUSMODE = 0x4a

/* This is an addition on the other single-board commands. This command is
 * recognised in single-board mode. The RS232 output goes in
 * pull-up mode and bus commands are immediatly recognised hereafter.
 * Note that when the board is in single-board mode, and eventually a bus
 * mode command is found, this command is not processed, but the board
 * switches to bus mode. The next (bus) command is processed regularly.
 */

const DGT_STARTBOOTLOADER = 0x4e

/* Makes a long jump to the FC00 boot loader code. Start FLIP now
 */

/* ------------------------------------------------------------------------ */
/* Commands resulting in returning message(s):                              */
/* ------------------------------------------------------------------------ */

const DGT_SEND_CLK = 0x41

/* Results in a DGT_MSG_BWTIME message
 */

const DGT_SEND_BRD = 0x42

/* Results in a DGT_MSG_BOARD_DUMP message
 */

const DGT_SEND_UPDATE = 0x43

/* Results in DGT_MSG_FIELD_UPDATE messages and DGT_MSG_BWTIME messages
 * as long as the board is in UPDATE mode  */

const DGT_SEND_UPDATE_BRD = 0x44

/* results in DGT_MSG_FIELD_UPDATE messages
 * as long as the board is in UPDATE_BOARD mode  */

const DGT_RETURN_SERIALNR = 0x45

/* Results in a DGT_MSG_SERIALNR message
 */

const DGT_RETURN_BUSADRES = 0x46

/* Results in a DGT_MSG_BUSADRES message
 */

const DGT_SEND_TRADEMARK = 0x47

/* Results in a DGT_MSG_TRADEMARK message
 */

const DGT_SEND_EE_MOVES = 0x49

/* Results in a DGT_MSG_EE_MOVES message
 */

const DGT_SEND_UPDATE_NICE = 0x4b

/* results in DGT_MSG_FIELD_UPDATE messages and DGT_MSG_BWTIME messages,
 * the latter only at time changes,
 * as long as the board is in UPDATE_NICE mode*/

const DGT_SEND_BATTERY_STATUS = 0x4c

/* New command for bluetooth board. Requests the
 * battery status from the board.
 */

const DGT_SEND_VERSION = 0x4d

/* Results in a DGT_MSG_VERSION message
 */

const DGT_SEND_BRD_50B = 0x50

/* Results in a DGT_MSG_BOARD_DUMP_50 message: only the black squares
 */

const DGT_SCAN_50B = 0x51

/* Sets the board in scanning only the black squares. This is written
 * in EEPROM
 */

const DGT_SEND_BRD_50W = 0x52

/* Results in a DGT_MSG_BOARD_DUMP_50 message: only the black squares
 */

const DGT_SCAN_50W = 0x53

/* Sets the board in scanning only the black squares. This is written in
 * EEPROM.
 */

const DGT_SCAN_100 = 0x54

/* Sets the board in scanning all squares. This is written in EEPROM
 */

const DGT_RETURN_LONG_SERIALNR = 0x55

/* Results in a DGT_LONG_SERIALNR message
 */

const DGT_SET_LEDS = 0x60

/* Only for the Revelation II to switch a LED pattern on. This is a command that
 * has three extra bytes with data:
 * byte 1 - DGT_SET_LEDS (= 0x60)
 * byte 2 - size (= 0x04)
 * byte 3 - the pattern to display (to be determined, for now: 0 - off, 1 - on)
 * byte 4 - the start field
 * byte 5 - the end field
 * byte 6 - end of message (= 0x00)
 * Start and end field have the range 0..63 where 0 is a8 and 63 is h1. This
 * is compliant with the DGT field coding of the board.
 * For the future it is foreseen that this message can have more fields which
 * can be controlled. This means also a different size of the message.
 */

/* ------------------------------------------------------------------------ */
/* Clock commands, returns ACK message if mode is in UPDATE or UPDATE_NICE  */
/* ------------------------------------------------------------------------ */
/*
 * USB/Serial: You cannot send multiple clock commands directly after each
 * other. After each command, you should wait for the ack before sending the
 * next one. The ack is usually returned within one or two seconds.
 *
 * BT2.10: Clock commands are processed faster and more reliable, but note:
 * - When a clock command is sent, it is possible that an empty bwtime message
 * is returned (full of zeroes). This is to be ignored.
 * - Within 1000ms after sending a clock command (and receiving its ack),
 * subsequent clock update requests may return this ack instead of the clock
 * update. This is because the clock is polled only 2-3 times a second.
 *
 *
 */

const DGT_CLOCK_MESSAGE = 0x2b

/* This message contains a command for the clock. There are clock commands
 * for showing text, displaying icons, setting beep, clearing display,
 * and for setting clock times. All these clock commands are wrapped
 * within the DGT_CLOCK_MESSAGE command.
 *
 * byte 1 - DGT_CLOCK_MESSAGE (= 0x2b)
 * byte 2 - the size
 * byte 3 - DGT_CMD_CLOCK_START_MESSAGE (= 0x03)
 * byte 4 - one of the 7 clock command id's
 * byte 5..n-1 - the content (can be empty)
 * byte n - DGT_CMD_CLOCK_END_MESSAGE (= 0x00)
 *
 * The clock messages can be sent in any board mode, but a SetNRun message is
 *   only processed if the clock is in mode 23.
 * The board responds with an DGT_MSG_ACK message if the board is in UPDATE or
 *   UPDATE_NICE mode. If the board is in IDLE or UPDATE_BRD mode, then no
 *   ack message is returned.
 */

const DGT_CMD_CLOCK_DISPLAY = 0x01

/*
 * This command can control the segments of six 7-segment characters,
 * two dots, two semicolons and the two '1' symbols.
 * These are arranged as follows: "1A:BC 1D:EF" and thus capable of
 * showing clock times or movetext of up to 6 chars (not counting the dots/'1').
 * "1A:BC 1D:EF" : the A..F are six characters with 7 segments, the '1' is a
 * single segment and the ':' has both a dot and semicolon segment.
 *
 * byte 1  - DGT_CMD_CLOCK_MESSAGE (= 0x2b)
 * byte 2  - Size (= 0x0b)
 * byte 3  - DGT_CMD_CLOCK_START_MESSAGE (= 0x03)
 * byte 4  - DGT_CMD_CLOCK_DISPLAY (= 0x01)
 * byte 5  - 'C' location segments. Bits: 0x01=top segment, 0x02=right top,
 *           0x04=right bottom, 0x08=bottom, 0x10=left bottom, 0x20=left top,
 *           0x40=center segment.
 * byte 6  - 'B' location segments. See 'C' location for the available bits.
 * byte 7  - 'A' location segments.
 * byte 8  - 'F' location segments.
 * byte 9  - 'E' location segments.
 * byte 10 - 'D' location segments.
 * byte 11 - icons: Bitmask for displaying dots and one's. 0x01=right dot,
 *           0x02=right semicolon, 0x04=right '1', 0x08=left dot,
 *           0x10=left semicolon, 0x20=left '1'.
 * byte 12 - 0x03 if beep, 0x01 if no beep
 * byte 13 - DGT_CMD_CLOCK_END_MESSAGE (= 0x00)
 */

const DGT_CMD_CLOCK_ICONS = 0x02

/*
 * Used to control the clock icons like flags etc.
 *
 * byte 1 -  DGT_CMD_CLOCK_MESSAGE (= 0x2b)
 * byte 2 -  Size (= 0x0b)
 * byte 3 -  DGT_CMD_CLOCK_START_MESSAGE (= 0x03)
 * byte 4 -  DGT_CMD_CLOCK_ICONS (= 0x02)
 * byte 5 -  Icon data byte 0
 *           bit 0 - left TIME symbol
 *           bit 1 - left FISCH symbol
 *           bit 2 - left DELAY symbol
 *           bit 3 - left HGLASS symbol
 *           bit 4 - left UPCNT symbol
 *           bit 5 - left BYO symbol
 *           bit 6 - left END symbol
 *           bit 7 - not in use
 * byte 6 -  Icon data byte 1
 *           bit 0 - right TIME symbol
 *           bit 1 - right FISCH symbol
 *           bit 2 - right DELAY symbol
 *           bit 3 - right HGLASS symbol
 *           bit 4 - right UPCNT symbol
 *           bit 5 - right BYO symbol
 *           bit 6 - right END symbol
 *           bit 7 - not in use
 * byte 7 -  Icon data byte 2
 *           bit 0 - left period '1' symbol
 *           bit 1 - left period '2' symbol
 *           bit 2 - left period '3' symbol
 *           bit 3 - left period '4' symbol
 *           bit 4 - left period '5' symbol
 *           bit 5 - left flag symbol
 *           bit 6 - not in use
 *           bit 7 - not in use
 * byte 8 -  Icon data byte 3
 *           bit 0 - right period '1' symbol
 *           bit 1 - right period '2' symbol
 *           bit 2 - right period '3' symbol
 *           bit 3 - right period '4' symbol
 *           bit 4 - right period '5' symbol
 *           bit 5 - right flag symbol
 *           bit 6 - not in use
 *           bit 7 - not in use
 * byte 9 -  Special (clr / condat)
 *           bit 0 - if "1" all icons are cleared, if "0" no switch off
 *           bit 1 - sound symbol
 *           bit 2 - black - white symbol
 *           bit 3 - white - black symbol
 *           bit 4 - BAT signal
 *           bit 5 - sound symbol (again!)
 *           bit 6 - display pool selec "0" icons clear after time change
 *                   "1" icons stay visible
 *           bit 7 -
 * byte 10 - Reserved (= 0x00)
 * byte 11 - Reserved (= 0x00)
 * byte 12 - Reserved (= 0x00)
 * byte 13 - DGT_CMD_CLOCK_END_MESSAGE (= 0x00)
 */

const DGT_CMD_CLOCK_END = 0x03

/* This command clears the message and brings the clock back to the
 * normal display (showing clock times).
 * byte 1  - DGT_CMD_CLOCK_MESSAGE (= 0x2b)
 * byte 2  - Size (= 0x03)
 * byte 3  - DGT_CMD_CLOCK_START_MESSAGE (= 0x03)
 * byte 4  - DGT_CMD_CLOCK_END (= 0x03)
 * byte 5  - DGT_CMD_CLOCK_END_MESSAGE (= 0x00)
 */

const DGT_CMD_CLOCK_BUTTON = 0x08

/*
 * Requests the current button pressed (if any).
 * byte 1  - DGT_CMD_CLOCK_MESSAGE (= 0x2b)
 * byte 2  - Size (= 0x03)
 * byte 3  - DGT_CMD_CLOCK_START_MESSAGE (= 0x03)
 * byte 4  - DGT_CMD_CLOCK_BUTTON (= 0x08)
 * byte 5  - DGT_CMD_CLOCK_END_MESSAGE (= 0x00)
 */

const DGT_CMD_CLOCK_VERSION = 0x09

/* This commands requests the clock version.
 * byte 1  - DGT_CMD_CLOCK_MESSAGE (= 0x2b)
 * byte 2  - Size (= 0x03)
 * byte 3  - DGT_CMD_CLOCK_START_MESSAGE (= 0x03)
 * byte 4  - DGT_CMD_CLOCK_VERSION (= 0x09)
 * byte 5  - DGT_CMD_CLOCK_END_MESSAGE (= 0x00)
 */

const DGT_CMD_CLOCK_SETNRUN = 0x0a

/* This commands controls the clock times and counting direction, when
* the clock is in mode 23. A clock can be paused or counting down. But
* counting up isn't supported on current DGT XL's (1.14 and lower) yet.
*
* byte 1  - DGT_CMD_CLOCK_MESSAGE (= 0x2b)
* byte 2  - Size (= 0x0a)
* byte 3  - DGT_CMD_CLOCK_START_MESSAGE (= 0x03)
* byte 4  - DGT_CMD_CLOCK_SETNRUN (= 0x0a)
* byte 5  - left hours (|0x10 if left counts up, but not supported)
* byte 6  - left minutes
* byte 7  - left seconds
* byte 8  - right hours (|0x10 if right counts up, but not supported)
* byte 9  - right minutes
* byte 10 - right seconds
* byte 11 - 0x01: left counts down. 0x02: right counts down.
              0x04: pause clock. 0x08: toggle player at lever change.
* byte 12 - DGT_CMD_CLOCK_END_MESSAGE (= 0x00)
*/

const DGT_CMD_CLOCK_BEEP = 0x0b

/*
 * This clock command turns the beep on, for a specified time (64ms * byte 5)
 * byte 1 - DGT_CMD_CLOCK_MESSAGE (= 0x2b)
 * byte 2 - Size (= 0x04)
 * byte 3 - DGT_CLOCK_START_MESSAGE (= 0x03)
 * byte 4 - CMD_CLOCK_BEEP (= 0x0b)
 * byte 5 - The time in multiplies of 64ms. E.g. a value of 16 plays the beep for 16*64=1024ms.
 * byte 6 - DGT_CMD_CLOCK_END_MESSAGE (= 0x00)
 */

const DGT_CMD_CLOCK_ASCII = 0x0c

/*
 * This clock commands sends a ASCII message to the clock that can be displayed only by
 * the DGT3000.
 * byte 1  - DGT_CMD_CLOCK_MESSAGE (= 0x2b)
 * byte 2  - Size (= 0x0c)
 * byte 3  - DGT_CLOCK_START_MESSAGE (= 0x03)
 * byte 4  - CMD_CLOCK_ASCII (= 0x0c)
 * byte 5  - Character 1
 * byte 6  - Character 2
 * byte 7  - Character 3
 * byte 8  - Character 4
 * byte 9  - Character 5
 * byte 10 - Character 6
 * byte 11 - Character 7
 * byte 12 - Character 8
 * byte 13 - Beep value 0 - 15, 0 is no beep, else else 62.5ms + (value / 16) * 1000ms
 * byte 14 - DGT_CMD_CLOCK_END_MESSAGE (= 0x00)
 */

/* ------------------------------------------------------------------------ */
/* DESCRIPTION OF THE MESSAGES FROM BOARD TO PC

A message consists of three header bytes:
MESSAGE ID             one byte, MSB (MESSAGE BIT) always 1
MSB of MESSAGE SIZE    one byte, MSB always 0, carrying D13 to D7 of the
					   total message length, including the 3 header byte
LSB of MESSAGE SIZE    one byte, MSB always 0, carrying  D6 to D0 of the
					   total message length, including the 3 header bytes
followed by the data:
0 to ((2 EXP 14) minus 3) data bytes, of which the MSB is always zero.
*/

/* ------------------------------------------------------------------------ */
/* DEFINITION OF THE BOARD-TO-PC MESSAGE ID CODES and message descriptions */

/* the Message ID is the logical OR of MESSAGE_BIT and ID code */
const MESSAGE_BIT = 0x80

const MESSAGE_MASK = 0x7F

/* ID codes: */
const DGT_NONE = 0x00
const DGT_BOARD_DUMP = 0x06
const DGT_BWTIME = 0x0d
const DGT_FIELD_UPDATE = 0x0e
const DGT_EE_MOVES = 0x0f
const DGT_BUSADRES = 0x10
const DGT_SERIALNR = 0x11
const DGT_TRADEMARK = 0x12
const DGT_VERSION = 0x13

/* Added for Draughts board  */
const DGT_BOARD_DUMP_50B = 0x14
const DGT_BOARD_DUMP_50W = 0x15

/* Added for Bluetooth board */
const DGT_BATTERY_STATUS = 0x20
const DGT_LONG_SERIALNR = 0x22

/* ------------------------------------------------------------------------ */
/* Macros for message length coding (to avoid MSB set to 1) */

// const BYTE    char

// const LLL_SEVEN(a) ((BYTE)(a & 0x7f))            /* 0000 0000 0111 1111 */
// const LLH_SEVEN(a) ((BYTE)((a & 0x3F80)>>7))   /* 0011 1111 1000 0000 */

/* ------------------------------------------------------------------------ */
/* DGT_MSG_BOARD_DUMP is the message that follows on a DGT_SEND_BOARD
 * command
 */
const DGT_MSG_BOARD_DUMP = (MESSAGE_BIT | DGT_BOARD_DUMP)
const DGT_SIZE_BOARD_DUMP = 67
const DGT_SIZE_BOARD_DUMP_DRAUGHTS = 103

/* message format:
 * byte 0: DGT_MSG_BOARD_DUMP
 * byte 1: LLH_SEVEN(DGT_SIZE_BOARD_DUMP) (=0 fixed)
 * byte 2: LLL_SEVEN(DGT_SIZE_BOARD_DUMP) (=67 fixed)
 * byte 3-66: Pieces on position 0-63
 *
 * Board fields are numbered from 0 to 63, row by row, in normal reading
 * sequence. When the connector is on the left hand, counting starts at
 * the top left square. The board itself does not rotate the numbering,
 * when black instead of white plays with the clock/connector on the left hand.
 * In non-rotated board use, the field numbering is as follows:
 *
 * Field A8 is numbered 0
 * Field B8 is numbered 1
 * Field C8 is numbered 2
 * ..
 * Field A7 is numbered 8
 * ..
 * Field H1 is numbered 63
 *
 * So the board always numbers the black edge field closest to the connector
 * as 57.
 */

/* Piece codes for chess pieces: */
const EMPTY = 0x00
const WPAWN = 0x01
const WROOK = 0x02
const WKNIGHT = 0x03
const WBISHOP = 0x04
const WKING = 0x05
const WQUEEN = 0x06
const BPAWN = 0x07
const BROOK = 0x08
const BKNIGHT = 0x09
const BBISHOP = 0x0a
const BKING = 0x0b
const BQUEEN = 0x0c
const PIECE1 = 0x0d /* Magic piece: Draw */
const PIECE2 = 0x0e /* Magic piece: White win */
const PIECE3 = 0x0f /* Magic piece: Black win */

/* For the draughts board */
const WDISK = 0x01
const BDISK = 0x04
const WCROWN = 0x07
const BCROWN = 0x0a

/* ------------------------------------------------------------------------ */
/* message format DGT_MSG_BOARD_DUMP_50B                                    */
const DGT_MSG_BOARD_DUMP_50B = (MESSAGE_BIT | DGT_BOARD_DUMP_50B)
const DGT_SIZE_BOARD_DUMP_50B = 53

/* byte 0: DGT_MSG_BOARD_DUMP_50B
 * byte 1: LLH_SEVEN(DGT_SIZE_BOARD_DUMP) (=0 fixed)
 * byte 2: LLL_SEVEN(DGT_SIZE_BOARD_DUMP) (=53 fixed)
 * byte 3-52: Pieces on position 0-50
 */

/* ------------------------------------------------------------------------ */
/* message format DGT_MSG_BOARD_DUMP_50W                                    */
const DGT_MSG_BOARD_DUMP_50W = (MESSAGE_BIT | DGT_BOARD_DUMP_50W)
const DGT_SIZE_BOARD_DUMP_50W = 53

/* byte 0: DGT_MSG_BOARD_DUMP_50w
 * byte 1: LLH_SEVEN(DGT_SIZE_BOARD_DUMP) (=0 fixed)
 * byte 2: LLL_SEVEN(DGT_SIZE_BOARD_DUMP) (=53 fixed)
 * byte 3-52: Pieces on position 0-50
 */

/* ------------------------------------------------------------------------ */
/* message format DGT_MSG_BWTIME                                            */
const DGT_MSG_BWTIME = (MESSAGE_BIT | DGT_BWTIME)
const DGT_SIZE_BWTIME = 10

/*
 * There are two possible distinct BwTime messages: 1) Clock Times, 2) Clock Ack.
 * The total size is always 10 bytes and the first byte is always DGT_MSG_BWTIME (=0x4d).
 * If the (4th byte & 0x0f) equals 0x0a, or if the (7th byte & 0x0f) equals 0x0a, then the
 * message is a Clock Ack message. Otherwise it is a Clock Times message.
 *
 * Clock Times:
 *
 * byte 0: DGT_MSG_BWTIME
 * byte 1: LLH_SEVEN(DGT_SIZE_BWTIME) (=0 fixed)
 * byte 2: LLL_SEVEN(DGT_SIZE_BWTIME) (=10 fixed)
 * byte 3:
 * D0-D3: Hours (units, 0-9 Binary coded) of right player.
 *  (If (byte 3 & 0x0f) is 0x0a, then the msg is a Clock Ack message instead)
 * D4: 1 = Flag fallen for right player, and clock blocked to zero
 *     0 = not the above situation
 * D5: 1 = Time per move indicator on for right player ( i.e. Bronstein, Fischer)
 *     0 = Time per move indicator off for right player
 * D6: 1 = Right players flag fallen and indicated on display,
 *           clock possibly still running (e.g. activation of next time period)
 *     0 = not the above situation
 * (D7 is MSB)
 * byte 4: Minutes (0-59, BCD coded)
 * byte 5: Seconds (0-59, BCD coded)
 * byte 6-8: the same for the left player
 * byte 9: Clock status byte: 7 bits
 * D0 (LSB): 1 = Clock running
 *           0 = Clock stopped by Start/Stop
 * D1: 1 = tumbler position high on right player (front view: / , right side high)
 *     0 = tumbler position high on left player (front view: \, left side high)
 * D2: 1 = Battery low indication on display
 *     0 = no battery low indication on display
 * D3: 1 = Right player's turn
 *     0 = not Right player's turn
 * D4: 1 = Left player's turn
 *     0 = not Left player's turn
 * D5: 1 = No clock connected; reading invalid
 *     0 = clock connected, reading valid
 * D6: not used (read as 0)
 * D7:  Always 0
 *
 *
 * Clock Ack:
 *
 * byte 0: DGT_MSG_BWTIME
 * byte 1: LLH_SEVEN(DGT_SIZE_BWTIME) (=0 fixed)
 * byte 2: LLL_SEVEN(DGT_SIZE_BWTIME) (=10 fixed)
 * byte 3-9: used to construct a 4-byte ack message in the following way:
 *
 * ack0 = ((byte4) & 0x7f) | ((byte6 << 3) & 0x80);
 * ack1 = ((byte5) & 0x7f) | ((byte6 << 2) & 0x80);
 * ack2 = ((byte7) & 0x7f) | ((byte3 << 3) & 0x80);
 * ack3 = ((byte8) & 0x7f) | ((byte3 << 2) & 0x80);
 *
 * So ack0 is constructed from bits 0-6 from byte4, and bit 4 from byte6 (which goes
 *   to ack0's 7th bit).
 *
 * ack0: 0x40 indicates an error, e.g. sending a CMD_CLOCK_SETNRUN while the
 *   clock is not in mode 23.
 * ack0 is 0x10 for normal ack messages.
 * ack1: if 8th bit set (i.e. (ack1 & 0x80)==0x80), then it is an auto-generated ack.
 *   Otherwise it is a response to a command.
 *
 * Auto-generated acks:
 *    ack1==0x81: ready
 *    ack1==0x88: Button pressed. ack3: Back=0x31; Plus=0x32; Run=0x33; Minus=0x34; OK=0x35.
 *    ack1==0x8a: mode 23? ack3 contains the mode?
 *    ack1==0x90: not in mode 23?
 *
 *  Response to command ack's:
 *  ack1 is the command that was ack'ed.
 *    ack1==0x01: Display ack.
 *    ack1==0x08: Buttons ack, but no button information is returned though.
 *    ack1==0x09: Version ack. ack2>>4 is main version, ack2&0x0f is sub version.
 *    ack1==0x0a: SetNRun ack.
 *    ack1==0x0b: Beep ack.
 *
 *
 */

/* ------------------------------------------------------------------------ */
/* message format DGT_MSG_FIELD_UPDATE:                                     */
const DGT_MSG_FIELD_UPDATE = (MESSAGE_BIT | DGT_FIELD_UPDATE)
const DGT_SIZE_FIELD_UPDATE = 5

/* byte 0: DGT_MSG_FIELD_UPDATE
 * byte 1: LLH_SEVEN(DGT_SIZE_FIELD_UPDATE) (=0 fixed)
 * byte 2: LLL_SEVEN(DGT_SIZE_FIELD_UPDATE) (=5 fixed)
 * byte 3: field number (0-63) which changed the piece code
 * byte 4: piece code including EMPTY, where a non-empty field became empty
 */

/* ------------------------------------------------------------------------ */
/* message format: DGT_MSG_TRADEMARK which returns a trade mark message     */
const DGT_MSG_TRADEMARK = (MESSAGE_BIT | DGT_TRADEMARK)

/* byte 0: DGT_MSG_TRADEMARK
 * byte 1: LLH_SEVEN(DGT_SIZE_TRADEMARK)
 * byte 2: LLL_SEVEN(DGT_SIZE_TRADEMARK)
 * byte 3-end: ASCII TRADEMARK MESSAGE, codes 0 to 0x3F
 * The value of DGT_SIZE_TRADEMARK is not known beforehand, and may be in the
 * range of 0 to 256
 * Current trade mark message: ...
 */

/* ------------------------------------------------------------------------ */
/* Message format DGT_MSG_BUSADRES return message with bus adres            */
const DGT_MSG_BUSADRES = (MESSAGE_BIT | DGT_BUSADRES)
const DGT_SIZE_BUSADRES = 5

/* byte 0: DGT_MSG_BUSADRES
 * byte 1: LLH_SEVEN(DGT_SIZE_BUSADRES)
 * byte 2: LLL_SEVEN(DGT_SIZE_BUSADRES)
 * byte 3,4: Busadres in 2 bytes of 7 bits hexadecimal value
 *           byte 3: 0bbb bbbb with bus adres MSB 7 bits
 *           byte 4: 0bbb bbbb with bus adres LSB 7 bits
 * The value of the 14-bit busadres is het hexadecimal representation
 * of the (decimal coded) serial number
 * i.e. When the serial number is "01025 1.0" the busadres will be
 *      byte 3: 0000 1000 (0x08)
 *      byte 4: 0000 0001 (0x01)
 */

/* ------------------------------------------------------------------------ */
/* Message format DGT_MSG_SERIALNR return message with bus adres            */
const DGT_MSG_SERIALNR = (MESSAGE_BIT | DGT_SERIALNR)
const DGT_SIZE_SERIALNR = 8

/* Returns 5 ASCII decimal serial number:
 * byte 0-5 serial number string, sixth byte is LSByte
 */

/* ------------------------------------------------------------------------ */
/* Message format DGT_LONG_SERIALNR                                         */
const DGT_MSG_LONG_SERIALNR = (MESSAGE_BIT | DGT_LONG_SERIALNR)
const DGT_SIZE_LONG_SERIALNR = 13

/* byte 0: DGT_LONG_SERIALNR
 * byte 1:  0
 * byte 2: DGT_SIZE_ LONG_SERIALNR
 * byte 3-12: 10 ASCII decimal serial number
 * The 10th character is Least Significant
 */

/* ------------------------------------------------------------------------ */
/* Message format DGT_MSG_VERSION return message with bus adres             */
const DGT_MSG_VERSION = (MESSAGE_BIT | DGT_VERSION)
const DGT_SIZE_VERSION = 5

/* byte 0: DGT_MSG_VERSION
 * byte 1: LLH_SEVEN(DGT_SIZE_VERSION)
 * byte 2: LLL_SEVEN(DGT_SIZE_VERSION)
 * byte 3,4: Version in 2 bytes of 7 bits hexadecimal value
 *           byte 3: 0bbb bbbb with main version number MSB 7 bits
 *           byte 4: 0bbb bbbb with sub version number LSB 7 bits
 * The value of the version is coded in binary
 * i.e. When the number is "1.02" the busadres will be
 *      byte 3: 0000 0001 (0x01)
 *      byte 4: 0000 0010 (0x02)
 */

/* ------------------------------------------------------------------------ */
/* Retrieve the battery status from the Bluetooth board                     */
const DGT_MSG_BATTERY_STATUS = (MESSAGE_BIT | DGT_BATTERY_STATUS)
const DGT_SIZE_BATTERY_STATUS = 7

/* byte  0: DGT_MSG_BATTERY_STATUS
 * byte  1: LLH_SEVEN(DGT_SIZE_BATTERY_STATUS)
 * byte  2: LLL_SEVEN(DGT_SIZE_BATTERY_STATUS)
 * byte  3: Current battery capacity left in %
 * byte  4: Running/Charging time left (hours)    (0x7f is N/A)
 * byte  5: Running/Charging time left (minutes)  (0x7f is N/A)
 * byte  6: On time (hours)
 * byte  7: On time (minutes)
 * byte  8: Standby time (days)
 * byte  9: Standby time (hours)
 * byte 10: Standby time (minutes)
 * byte 11: Status bits
 *         bit 0: Charge status
 *         bit 1: Discharge status
 *         bit 2: Not yet implemented
 *         bit 3: Not yet implemented
 *         bit 4: Not yet implemented
 *         bit 5: Not yet implemented
 *         bit 6: Not yet implemented
 *         bit 7: Not yet implemented
 */

/* ------------------------------------------------------------------------ */
/* DGT_SIZE_EE_MOVES is defined in dgt_ee1.h: current (0x2000-0x100+3)      */
const DGT_MSG_EE_MOVES = (MESSAGE_BIT | DGT_EE_MOVES)

/* Message format DGT_MSG_EE_MOVES, which is the contens of the storage array
 * byte 0: DGT_MSG_EE_MOVES
 * byte 1: LLH_SEVEN(DGT_SIZE_EE_MOVES)
 * byte 2: LLL_SEVEN(DGT_SIZE_EE_MOVES)
 * byte 3-end: field change storage stream: See defines below for contents
 *
 * The DGT_MSG_EE_MOVES message contains the contens of the storage,
 * starting with the oldest data, until the last written changes, and will
 * always end with EE_EOF
 */
