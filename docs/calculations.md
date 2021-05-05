# Constants
I2C_speed = 100 [kbit/s] = 12.5 [kbyte/s], standard 100 kHz
<br>I2C_request_length = 24 [bits] = 3 [bytes], mostly (address + command + data)

# Cycle time depending on count of inputs
Reading of one chip with 8 inputs take 1/I2C_speed * I2C_request_length = 0.24 [ms]
<br>Reading 100 chips (800 inputs) takes 24 ms

>Please consider that the current state of gobrail lacks in following points to reach this theoretical value:
* "Type 2" board is designed for flexibility not for speed, therefore only 4 inputs can be read on that chip
* output rail devices are directly chained to input rail devices, so each switch action with a delay (e.g. 100 ms for turnouts) will decrease the speed directly
This issues will be fixed in the future.

# How long a sensor needs to be active to be recognized by controller?
This directly depends on the cycle time. The active time must be greater than cycle time.

# Can a sensor made active long enough by passing a locomotive?
This question is important for passing sensors, e.g. optical gates. When the time will not long enough, the sensor needs an additional latch. So lets make a worst case calculation.

minimal_locomotive_length = 5 [cm], my shortest part is the BR81 with 7 cm
maximal_train_speed = 0.4 [m/s], measured by 10 times passing a loop by my newest train BR285

passing_time = minimal_locomotive_length / maximal_train_speed = 0.125 [s] = 125 [ms]

This means, when no output rail device is currently blocks or increase the cycle time, we don't need a latch.