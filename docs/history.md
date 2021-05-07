## The first years
1977 Become an owner of a "Starterset" with 1 locomotive and 2 wagons of "Berliner TT Bahnen"
1978 Small analog system with some track switches illuminated houses and a tunnel
1985 Full featured analog system with 36 buttons for track switches, illumination of houses and signals on a base area of 1.00 x 1.40 qm
1986 System was mothballed

There were ~8m track material, 13 track switches, 6 signals, 10 decoupling mechanics, 5 locomotives, some freight and some person wagons, all of the 80's and 90's.

## Revival and cleanup
2004 The system was startup again. Houses, trees, road and grass were cleaned. The track system was corroded, the motors of locomotives where partially blocked or stuttering. 
After some hours of cleanup the system worked again, but more randomly than stable. There where a few relays which won't work smoothly.
Cables brake from time to time, old solder connections brake. In summary, it was now fun to play.

## First core removal and start developing for digitalizing
2005 All copper material was removed including relays, switches ...
First idea to digitalize the system was born. The first hardware revision was build with NAND, OR and shift registers.
The software was running on a Basic stamp. The prototype was able to switch lamps and signals, but was much to slow for all the planned stuff.
The "Version 1" and "Version 2" was done.

## Using an arduino board and I2C
2006 A hardware with 1x PCA9501 which cascades 2x PCA9533 was developed, including a circuit board, exposed, etched and assembled by my own.
The software was written in C++ for arduino uno. The hardware and software was designed to drive lamps, signals and track switches. 
A first prototype drives a track switch stable for some hours, resulting in a overheated and destroyed magnet system.
The configuration part was developed to detect new boards with i2c chips and store configuration to EEPROM.
The "Version 3" was born.

## Architectural ideas and hardware simplification
2009 Due to financial limitations the customized board was no good idea, because prepare of exposing masks and etching for SMD components
is very difficult. Therefore some standardization was preferred, which also simplifies the prototyping procedure with bread boards.
The "Version 4" was born and finally 2013 the "Type 2" board with "PCA9501" was used for further software development.

At software level the system was divided in "board" (knows chips on a board), "board interface" (knows all boards), "extDev" (knows all devices), "Modellbahn" (API, CLI, UI).
The project lacks of good knowledge of C++ programming language and time to make progress, so it was mothballed again.

## Restart the project with go
2021 The track material of model railroad was completely removed, also lamps, signals and most of landscape. 
A new analog "Starterset" was purchased with shiny track. This set will be the prototype for the new system to test 
control of track switches, signal and lamps. Using of locomotive decoder is planned.

Software was rewritten in go using the gobot framework, which provides the advantage to use nearly each adaptor I (you) want. 
With my slow progress, this would be an big advantage for me.

Startup is done with digispark and mcp2221. The project was published initially on github.

## Future
The main progress will be limited to the winter month.
