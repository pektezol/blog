---
title: "making the best bitreader library in golang"
date: "2023-09-17"
author: "Arda Serdar Pektezol"
---
## What is this?

At [github.com/pektezol/bitreader](https://github.com/pektezol/bitreader), has the (self-proclaimed) best library in Golang that works as a bit-level reader. With its simplicity and ease-of-use, pektezol/bitreader is ready to be used widely in the sub-byte environments.

## Disclaimer

I want to start off by saying that I started this project almost a year ago when I wrote this post. In this whole year, I have immensely  improved my Golang and coding knowledge, know-how, and logic - compared to myself. At the moment, I still feel like a beginner, but it is always nice seeing your -not even so distant- work and cringe at it.

## Why?

The whole need for this library sprawled when I was working on my demo parser for the video game Portal 2. Essentially, what it is is that Portal 2 can record and output what is called "demos" using in-game commands, that record every single thing that is done in the game, during the whole period of its recording.

Being a binary file, it obviously needs to be parsed in order to extract every single information the demo file has to offer. Luckily, I don't have to get my hands dirty with reverse-engineering this complicated output. [nekz.me/dem](https://nekz.me/dem) by [@NeKzor](https://github.com/NeKzor) offers more than a general information about why, where, and how the information would be read, in order to receive meaningful output.

When first getting into this, it becomes clear that you, at least somewhat, need some utils functions to read data. Moreover, there are times where bit-packing occurs and lower than 8 number of bits need to be read into that value - which Golang does not natively support reading lower than a byte.

While, it is true that you can just use bitwise operations here and there to mitigate this, why not implement your own bit reader to be very flexible, comfortable, and heck, even release it as a library?

One last point to make, is that while yes; bitreader libraries do indeed exist in the Golang ecosystem (even a standard library one), essentially none of them had what I needed to accomplish (or I just didn't see it);

- It had to support little-endianness: Some supported only big-endian, some only little-endian, but I don't remember seeing one that supported both. I needed little-endian since the demo information was stored like that. Why not release a library that just supports both?
- It had to be simple: I didn't really want to work with complex stuff, just to read some bits. Sounds childish, but hey, what are you going to do about it.
- Some other stuff that I thought at the time, where I was really inexperienced with Golang and bitwise stuff. Probably wouldn't make sense to write now.

## How?

If the commit history is checked (please don't), it becomes clear that I actually don't know what I am doing at first. First commit dated at 2022-09-04, some weird stuff is happening with turning the byte data to string, then converting back to an int- I don't know, don't judge me. Essentially not even using any bitwise operations.. in a library that aims to be a bitreader.

With some more inspirations from other bitreader libraries, which does become clear when you read the code, and a big-time help from [@mlugg](https://github.com/mlugg), a better version was released that ACTUALLY used bitwise operations lol.

With some more tweaks, additions, adjustments, and a year later, BitReader v1.4.3 at the moment of writing is open to the public. So, how does it work exactly?

Simply, a reader is created when calling ```bitreader.NewReader(io.Reader, bool)``` or ```bitreader.NewReaderFromBytes([]byte, bool)```. First parameter becomes the stream for that reader, either in io.Reader form, or generated from a byte slice. The bool from the parameters specify the endianness state of the reader; true being little-endian and false being big-endian. This reader also has two more fields, where the ```index``` and ```currentByte``` data is kept. ```currentByte``` having the currently read byte, ```index``` containing the pointer position on that ```currentByte```.

When reading any amount of bits, the reader reads the stream one bit at a time, while reading a single byte to the buffer until the end of the byte is reached, so we can read the next byte as a buffer and keep reading it bit-by-bit. Seems a bit inefficient to do so, but the performance of Golang and modern hardware makes it indistinguishable.

On the case of endianness, this code block explains how the value is retrieved:

```go
var val bool
if reader.littleEndian {
	val = (reader.currentByte & (1 << reader.index)) != 0
} else {
	val = (reader.currentByte & (1 << (7 - reader.index))) != 0
}
```

Basing everything on this, more functions were added to skip bits/bytes, read null-terminated and length-specified string, read bits/bytes into a []byte, fork the reader to duplicate it; and wrapper functions for all of the signed/unsigned integer types, and the rest of the functions, if one would rather live on the edge and not handle errors.

### Usage Examples

Straigt from the README, you can see all use case examples from the library:

```go
import "github.com/pektezol/bitreader"

// ioStream:        io.Reader  Data to read from an io stream
// byteStream:      []byte     Data to read from a byte slice
// littleEndian:    bool       Little-endian(true) or big-endian(false) state
reader := bitreader.NewReader(ioStream, le)
reader := bitreader.NewReaderFromBytes(byteStream, le)

// Fork Reader, Copies Current Reader
newReader, err := reader.Fork()

// Read Total Number of Bits Left
bits, err := reader.ReadRemainingBits()

// Read First Bit
state, err := reader.ReadBool()

// Read Bits/Bytes
value, err := reader.ReadBits(64)       // up to 64 bits
value, err := reader.ReadBytes(8)       // up to 8 bytes

// Read String
text, err := reader.ReadString()            // null-terminated
text, err := reader.ReadStringLength(256)   // length-specified

// Read Bits/Bytes into Slice
arr, err := reader.ReadBitsToSlice(128)
arr, err := reader.ReadBytesToSlice(64)

// Skip Bits/Bytes
err := reader.SkipBits(8)
err := reader.SkipBytes(4)

// Wrapper functions
state := reader.TryReadBool()           // bool
value := reader.TryReadInt1()           // uint8
value := reader.TryReadUInt8()          // uint8
value := reader.TryReadSInt8()          // int8
value := reader.TryReadUInt16()         // uint16
value := reader.TryReadSInt16()         // int16
value := reader.TryReadUInt32()         // uint32
value := reader.TryReadSInt32()         // int32
value := reader.TryReadUInt64()         // uint64
value := reader.TryReadSInt64()         // int64
value := reader.TryReadFloat32()        // float32
value := reader.TryReadFloat64()        // float64
value := reader.TryReadBits(64)         // uint64
value := reader.TryReadBytes(8)         // uint64
text := reader.TryReadString()          // string
text := reader.TryReadStringLength(64)  // string
arr := reader.TryReadBitsToSlice(1024)  // []byte
arr := reader.TryReadBytesToSlice(128)  // []byte
bits := reader.TryReadRemainingBits()   // uint64
```

## What's Next?

I don't really have much else to add as a feature, but who knows? There might still be a couple bugs lying around where I may need your help. Using [GitHub Issues](https://github.com/pektezol/BitReader/issues/new), you can report a bug that you encountered and/or request a feature that you would like to be added.

Honestly, I am proud of this work and what it accomplishes. I have learned a lot, and I am continuing to learn every single day. If you liked what you see and/or read, consider giving a star to the [GitHub repository](https://github.com/pektezol/bitreader). Thank you for your time.