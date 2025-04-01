# Description

This is my response to an code exam. This is the "library" part

## PROMPT
```
Write a library and CLI for generating random PIN codes. You probably know what a PIN code is; it's a short sequence of numbers, often used as a passcode for bank cards.

Here are the requirements:

- The library should export a function that returns a batch of 1,000 PIN codes in random order
- Each PIN code in the batch should be unique
- Each PIN should be:
    - 4 digits long
    - Two consecutive digits should not be the same (e.g. 1156 is invalid)
    - Three consecutive digits should not be incremental (e.g. 1236 is invalid)
- The library should have automated tests.
```

## REASONING

The strategy I chose for generating the codes was to
1. first generate all possible codes. For 4 digit or even 6 digit codes this isn't that much, programmatically
2. remove any codes that are not valid
3. shuffle the possible codes to induce randomness
4. return the requested number of codes

I chose this method because:
- it made it easier to determine if the number of requested codes was invalid (instead of having to do some math ahead of time to calculate possible codes)
- it is inherintly unique and we don't have to clean the codes
- The other method I had thought of: generating a random number, and filling a "bucket" of valid codes was inefficient. Especially as the number of requested codes gets close to the max. You could spend much time "guessing" for those last numbers. This later method also introduces a bug if your number of requested codes is greater than the number of possible codes.

I also chose setters for options, and didn't hardcode the digit length. This way, this method could easily be expanded if a user wanted different length codes or more codes than the 1000 that the prompt had 
