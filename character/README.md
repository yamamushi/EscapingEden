# Character Manager

The Character manager only manages character data. Any changes to a character pass through this manager, and any
requests about the state of a character pass through this manager. 

If a player presses a move key, the request comes through the player manager. All player actions come through this.

If a monster attacks a player, or a player is injured in some way, the final storing of the player state comes here.

The player may try to move, and the player manager will check to see if the move is valid with the game manager, and
if it is valid, it will perform the action and report the results as needed. 

This behavior and pattern is repeated for a variety of other actions, such as picking up an item, attacking a monster,
and so on.

Ultimately the character manager requires a high amount of throughput, because of this much of the work it performs
is done in stateless goroutines.

While the manager does read/write data to the database, most of the work is done with players loaded into memory. Because
of this, it is essential that the player manager be given time to write the player data to the database before the
server is shut down. Otherwise, player data will be lost.
