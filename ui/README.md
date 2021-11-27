
A console belongs to every client connection, and manages what is drawn on their session.

A console is a collection of windows, which are the actual drawable objects.

Windows have properties, but ultimately they must have a draw function to display themselves.

Windows can be moved around, and resized, and can be hidden.

Windows can request to be redrawn, and will be redrawn when the console is next drawn.
Windows can request properties from the game server, and depending on the permissions, they will receive a response.


There are two types of drawing mechanisms. 

    - You can write to the window content, and it will be drawn as is. This accepts string input.
    - You can draw to the PointMap, and it will be drawn to the grid of points. 
    - DO NOT mix the two, or your content will be overwritten.
    - Drawing to the PointMap has no bounds checking, so you must ensure that you are drawing to the correct area.