## Example "type2_2turnout_4keys"

#### General
This operates 2 turnouts controlled by 2 keys. Two boards are required, a basic "Type 2" board and an additional breadboard for wiring the keys. The plan connects the "Turnout1" to the switch "K1rt", which should be a real switch. When using a button you will determine that pressed and not pressed state will directly affect the position of the turnout. "Turnout2" is connected to a button "K2gn" in toggle mode. This means the turnout will change the position at each press of the button. 

#### Mod: remove LED's
The 4 LED's are used as an optical feedback and are not really needed for operation. When removing an LED replace LED and the 150 Ohm current limiting resistor by an 10 kOhm ... 15 kOhm resistor as pull up to improve stability.

#### Mod: Connect turnouts together
It is possible to switch one turnout by the state of another one. In the given example plan just change the line `"Connect": "K1rt"` to `"Connect": "Turnout2"`.
