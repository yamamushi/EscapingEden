# Window

One of the most important concepts in the UI portion of this project is the Window.

The window ultimately is responsible for handling events from the console, which in most cases originated from the user at some level up the chain.

It is responsible for updating its own internal content and printing it to its point map, which is then retrieved by the console at a later time and parsed for printing to the screen.

Each directory in this folder contains a different type of window, as well as the functions they use to operate.

The base interface is the Window interface, and its methods can be overridden by individual window types for customizing their behavior.

A window can attempt to write to a pointmap position that is outside the bounds of the window space it occupies in the console, however the console will ensure that only the portion of the window that is visible is output.

